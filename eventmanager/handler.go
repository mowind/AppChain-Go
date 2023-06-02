package eventmanager

import (
	"context"
	appchain "github.com/PlatONnetwork/AppChain-Go"
	"github.com/PlatONnetwork/AppChain-Go/common"
	"github.com/PlatONnetwork/AppChain-Go/core/types"
	"github.com/PlatONnetwork/AppChain-Go/crypto"
	"github.com/PlatONnetwork/AppChain-Go/ethclient"
	"github.com/PlatONnetwork/AppChain-Go/ethdb"
	"github.com/PlatONnetwork/AppChain-Go/event"
	"github.com/PlatONnetwork/AppChain-Go/innerbindings/helper"
	"github.com/PlatONnetwork/AppChain-Go/log"
	"math/big"
	"sort"
	"sync"
)

// EventManager Managing events issued on the master chain.
// Includes listening for events, storing events and assembling a batch of packaged events.
type EventManager struct {
	platonAddr string
	exit       chan struct{}
	db         ethdb.Database
	RCConfig   *RootChainContractConfig
	// When packing events, the height of the latest listened event - x = the block height of the packing cut-off.
	backNumbers uint64

	// Get events from the specified block height
	fromBlockNumber uint64
	// Event Storage
	blockLogs map[uint64][]*types.Log

	checkpointEventFeed event.Feed
	mu                  sync.RWMutex
}

type RootChainContractConfig struct {
	RootChainID        string
	StakingInfoAddress common.Address
	RootChainAddress   common.Address
}

func NewEventManager(platonAddr string, db ethdb.Database) *EventManager {
	eventManager := &EventManager{
		platonAddr:      platonAddr,
		exit:            make(chan struct{}),
		db:              db,
		backNumbers:     10,
		fromBlockNumber: 1,
		blockLogs:       make(map[uint64][]*types.Log, 0),
	}
	return eventManager
}

// SubscribeEvents Subscribe to Checkpoint events that occur on RootChain.
func (em *EventManager) SubscribeEvents(ch chan *types.Log) event.Subscription {
	return em.checkpointEventFeed.Subscribe(ch)
}

func (em *EventManager) Listen() error {
	// If it is an authenticator node, this rpc address needs to be configured.
	// Not required if it is a normal node.
	if em.platonAddr == "" {
		log.Warn("the rpc address for platon is empty, please check if it is required")
		return nil
	}
	client, err := ethclient.Dial(em.platonAddr)
	if err != nil {
		log.Error("Failed to connect to Platon's RPC address", "addr", em.platonAddr, "error", err)
		return err
	}
	newHeadChan := make(chan *types.Header)
	newHeadSubscribe, err := client.SubscribeNewHead(context.Background(), newHeadChan)
	if err != nil {
		close(newHeadChan)
		log.Error("listening to the block header fails", "error", err)
		return err
	}
	defer func() {
		newHeadSubscribe.Unsubscribe()
		client.Close()
	}()

	for {
		select {
		case <-em.exit:
			log.Warn("event listener exit")
			return nil
		case err := <-newHeadSubscribe.Err():
			log.Error("subscription failure", "error", err)
			// TODO 处理订阅区块头失败的情况
			return err
		case newHead := <-newHeadChan:
			log.Trace("listening for a new block", "blockNumber", newHead.Number, "blockHash", newHead.Hash().TerminalString())
			em.mu.RLock()
			fromBlockNumber := em.fromBlockNumber
			em.mu.RUnlock()
			filterParams := appchain.FilterQuery{
				FromBlock: new(big.Int).SetUint64(fromBlockNumber),
				ToBlock:   newHead.Number,
				Addresses: []common.Address{
					em.RCConfig.StakingInfoAddress,
					em.RCConfig.RootChainAddress,
				},
				Topics: [][]common.Hash{{helper.StakedID, helper.UnstakeInitID,
					helper.SignerChangeID, helper.StakeUpdateID, helper.NewHeaderBlockID}},
			}

			logs, err := client.FilterLogs(context.Background(), filterParams)
			if err != nil {
				log.Error("failed to get filtered logs", "fromBlock", filterParams.FromBlock, "toBlock", filterParams.ToBlock, "error", err)
				// TODO 处理获取事件失败的情况
				break
			}
			log.Debug("get event success", "fromBlock", filterParams.FromBlock, "toBlock", filterParams.ToBlock, "logLength", len(logs))
			blockLogsTemp := make(map[uint64][]*types.Log)
			for _, log := range logs {
				// checkpoint events are not stored and are notified directly to the special handling logic.
				// feed.Send()
				if log.Topics[0] == helper.NewHeaderBlockID {
					em.checkpointEventFeed.Send(log)
					continue
				}

				logs, ok := blockLogsTemp[log.BlockNumber]
				if !ok {
					logs = make([]*types.Log, 0)
				}
				logs = append(logs, &log)
				blockLogsTemp[log.BlockNumber] = logs
			}
			em.mu.Lock()
			// If a block that has already been listened to appears, it is skipped.
			// This does not occur normally.
			for k, v := range blockLogsTemp {
				if _, ok := em.blockLogs[k]; ok {
					continue
				}
				em.blockLogs[k] = v
			}
			// Make the latest block +1, as the starting block high for the next fetch event.
			em.fromBlockNumber = newHead.Number.Uint64() + 1
			em.mu.Unlock()
		}
	}
}

type BlockNumberListSort []uint64

func (bnl BlockNumberListSort) Len() int {
	return len(bnl)
}

func (bnl BlockNumberListSort) Less(i, j int) bool {
	return bnl[i] < bnl[j]
}

func (bnl BlockNumberListSort) Swap(i, j int) {
	bnl[i], bnl[j] = bnl[j], bnl[i]
}

// BuildEventList Get all the specified events in the range based on the start and end block heights.
func (em *EventManager) BuildEventList(startBlockNumber uint64, endBlockNumber uint64) (*big.Int, []*types.Log, error) {
	em.mu.RLock()
	defer em.mu.RUnlock()
	if endBlockNumber == 0 {
		// If it is the node that is out of the block, that logic is taken.
		// Calculate the cut-off block height for packing events based on the estimated inter-node synchronization block delay.
		endBlockNumber = em.fromBlockNumber - em.backNumbers - 1
	}
	if startBlockNumber > em.fromBlockNumber {
		log.Warn("starting block height is greater than the latest height listened to", "startBlockNumber", startBlockNumber, "latestHeight", em.fromBlockNumber-1)
		return nil, nil, nil
	}
	if endBlockNumber >= em.fromBlockNumber || endBlockNumber < startBlockNumber {
		log.Debug("Not enough events", "startBlockNumber", startBlockNumber, "latestHeight", em.fromBlockNumber-1,
			"backNumbers", em.backNumbers, "endBlockNumber", endBlockNumber)
		return nil, nil, nil
	}
	blockNumberList := make(BlockNumberListSort, 0)
	for blockNumber := range em.blockLogs {
		if blockNumber >= startBlockNumber && blockNumber <= endBlockNumber {
			blockNumberList = append(blockNumberList, blockNumber)
		}
	}
	sort.Sort(&blockNumberList)
	logList := make([]*types.Log, 0)
	for _, blockNumber := range blockNumberList {
		logs := em.blockLogs[blockNumber]
		logList = append(logList, logs...)
	}
	return new(big.Int).SetUint64(endBlockNumber), logList, nil
}

// PackStakeEvents Encodes a batch of events and constructs a transaction input, and computes a hash on the input.
// Range of events: [startBlockNumber,endBlockNumber]
// Return:
// 1. Get the event's cut-off block height
// 2. Input of transaction.
// 3. Hash of Input.
// 4. Error
func (em *EventManager) PackStakeEvents(startBlockNumber uint64, endBlockNumber uint64) (*big.Int, []byte, common.Hash, error) {
	stopBlockNumber, logs, _ := em.BuildEventList(startBlockNumber, endBlockNumber)
	if stopBlockNumber != nil {
		encodeInput, err := helper.EncodeStakeStateSync(stopBlockNumber, logs)
		if err != nil {
			log.Error("packed events fail", "startBlockNumber", startBlockNumber, "endBlockNumber", endBlockNumber,
				"stopBlockNumber", stopBlockNumber, "error", err)
			return nil, nil, common.ZeroHash, nil
		}
		return stopBlockNumber, encodeInput, crypto.Keccak256Hash(encodeInput), nil
	}
	return nil, nil, common.ZeroHash, nil
}

func (em *EventManager) Stop() {
	close(em.exit)
}
