package eventmanager

import (
	"context"
	appchain "github.com/PlatONnetwork/AppChain-Go"
	"github.com/PlatONnetwork/AppChain-Go/common"
	"github.com/PlatONnetwork/AppChain-Go/core/types"
	"github.com/PlatONnetwork/AppChain-Go/ethclient"
	"github.com/PlatONnetwork/AppChain-Go/ethdb"
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

	mu sync.RWMutex
}

type RootChainContractConfig struct {
	RootChainID        string
	StakingInfoAddress common.Address
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
				},
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

// PackEventList Out block node call to select the event to be packed.
func (em *EventManager) PackEventList(startBlockNumber uint64) (*big.Int, []*types.Log, error) {
	endBlockNumber := em.fromBlockNumber - em.backNumbers - 1
	return em.BuildEventList(startBlockNumber, endBlockNumber)
}

func (em *EventManager) BuildEventList(startBlockNumber uint64, endBlockNumber uint64) (*big.Int, []*types.Log, error) {
	em.mu.RLock()
	defer em.mu.RUnlock()
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

func (em *EventManager) Stop() {
	close(em.exit)
}
