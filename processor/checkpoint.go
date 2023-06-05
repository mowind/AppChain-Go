package processor

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"sync"

	"github.com/PlatONnetwork/AppChain-Go/accounts/abi"
	"github.com/PlatONnetwork/AppChain-Go/common"
	"github.com/PlatONnetwork/AppChain-Go/common/hexutil"
	cvm "github.com/PlatONnetwork/AppChain-Go/common/vm"
	"github.com/PlatONnetwork/AppChain-Go/consensus"
	"github.com/PlatONnetwork/AppChain-Go/core/cbfttypes"
	"github.com/PlatONnetwork/AppChain-Go/core/types"
	"github.com/PlatONnetwork/AppChain-Go/core/vm"
	"github.com/PlatONnetwork/AppChain-Go/core/vm/solidity/checkpoint"
	stypes "github.com/PlatONnetwork/AppChain-Go/core/vm/solidity/types"
	"github.com/PlatONnetwork/AppChain-Go/crypto"
	"github.com/PlatONnetwork/AppChain-Go/ethclient"
	"github.com/PlatONnetwork/AppChain-Go/event"
	"github.com/PlatONnetwork/AppChain-Go/internal/ethapi"
	"github.com/PlatONnetwork/AppChain-Go/log"
	"github.com/PlatONnetwork/AppChain-Go/processor/api"
	"github.com/PlatONnetwork/AppChain-Go/rpc"
	"github.com/PlatONnetwork/AppChain-Go/x/plugin"
	"github.com/PlatONnetwork/AppChain-Go/x/xutil"
	lru "github.com/hashicorp/golang-lru"
	"github.com/xsleonard/go-merkle"
)

type PendingCheckpoint checkpoint.ICheckpointSigAggregatorPendingCheckpoint
type Checkpoint checkpoint.ICheckpointSigAggregatorCheckpoint

type ChainReader interface {
	CurrentHeader() *types.Header
	GetHeaderByNumber(number uint64) *types.Header
}

type TxSigner interface {
	Nonce() uint64
	Sign(tx *types.Transaction, chainId *big.Int) (*types.Transaction, error)
}

type CheckpointProposal struct {
	Checkpoint *stypes.Checkpoint
	BlockNum   uint64
}

type CheckpointProcessor struct {
	chain ChainReader

	chainId        *big.Int
	platonClient   *ethclient.Client
	bft            consensus.Bft
	txPoolAPI      api.TxPoolAPI
	caller         api.Caller
	managerAccount TxSigner
	checkpointABI  *abi.ABI

	newHeadBlockCh chan *types.Log
	bftResultSub   *event.TypeMuxSubscription

	exitCh chan struct{}

	latestProposal *CheckpointProposal
	rootHashCache *lru.ARCCache
}

func NewCheckpointProcessor(
	chain ChainReader,
	chainId *big.Int,
	platonAddr string,
	bft consensus.Bft,
	txPoolAPI api.TxPoolAPI,
	caller api.Caller,
	magnerAccount TxSigner,
	newHeadBlockCh chan *types.Log,
	bftResultSub *event.TypeMuxSubscription,
) (*CheckpointProcessor, error) {
	p := &CheckpointProcessor{
		chain:          chain,
		chainId:        chainId,
		bft:            bft,
		txPoolAPI:      txPoolAPI,
		caller:         caller,
		managerAccount: magnerAccount,
		newHeadBlockCh: newHeadBlockCh,
		bftResultSub:   bftResultSub,
		checkpointABI:  vm.CheckpointABI(),
		exitCh:         make(chan struct{}),
	}

	client, err := ethclient.Dial(platonAddr)
	if err != nil {
		log.Error("Failed to connect to PlatON node", "addr", platonAddr, "err", err)
		return nil, err
	}
	p.platonClient = client

	p.rootHashCache, err = lru.NewARC(10)
	if err != nil {
		return nil, err
	}

	go p.loop()

	return p, nil
}

func (p *CheckpointProcessor) Stop() {
	close(p.exitCh)
}

func (p *CheckpointProcessor) loop() {
	defer p.bftResultSub.Unsubscribe()

	for {
		select {
		case log := <-p.newHeadBlockCh:
			if log == nil {
				continue
			}
			p.handleNewHeadBlock(log)
		case result := <-p.bftResultSub.Chan():
			if result == nil {
				continue
			}
			cbftRsult, ok := result.Data.(cbfttypes.CbftResult)
			if !ok {
				log.Error("Receive bft result type error")
				continue
			}
			block := cbftRsult.Block
			if block == nil {
				log.Error("Cbft result error: block is nil")
				continue
			}
			p.handleBlock(block)
		case <-p.exitCh:
			log.Info("Checkpoint processor stopping...")
			return
		}
	}
}

func (p *CheckpointProcessor) handleNewHeadBlock(log *types.Log) {
	// TODO: check if need to confirm checkpoint
	// TODO: submit confirm checkpoint transaction

}

func (p *CheckpointProcessor) handleBlock(block *types.Block) {
	// TODO: check if need to propose checkpoint
	// TODO: propose checkpoint
	p.shouldPropose(block)

	// TODO: should submit checkpoint
}

func (p *CheckpointProcessor) shouldPropose(block *types.Block) {
	validators, err := plugin.StakingInstance().GetValidator(block.NumberU64())
	if err != nil {
		log.Error("Failed to validator list", "number", block.Number, "hash", block.Hash, "err", err)
		return
	}

	validator, err := validators.FindNodeByID(p.bft.NodeID())
	if err != nil {
		log.Debug("Current node is not a validator", "number", block.Number(), "hash", block.Hash(), "err", err)
		return
	}

	currentValidator, err := validators.FindNodeByID(p.bft.CurrentProposer())
	if err != nil {
		log.Error("Failed to get current validator", "number", block.Number, "hash", block.Hash, "err", err)
		return
	}

	pending, err := p.pendingCheckpoint()
	if err != nil {
		log.Error("Failed to get pending checkpoint", "number", block.Number(), "hash", block.Hash(), "err", err)
		return
	}
	if pending != nil &&
		p.latestProposal != nil &&
		block.NumberU64()-pending.BlockNum.Uint64() < vm.NextProposeDelay &&
		bytes.Equal(pending.Checkpoint.Proposer[:], p.latestProposal.Checkpoint.Proposer[:]) &&
		pending.Checkpoint.Start.Cmp(p.latestProposal.Checkpoint.Start) == 0 &&
		pending.Checkpoint.End.Cmp(p.latestProposal.Checkpoint.End) == 0 &&
		bytes.Equal(pending.Checkpoint.RootHash[:], p.latestProposal.Checkpoint.RootHash[:]) &&
		bytes.Equal(pending.Checkpoint.AccountHash[:], p.latestProposal.Checkpoint.AccountHash[:]) {
		// Already send proposal
		log.Debug("Already sign checkpoint proposal", "number", block.Number, "hash", block.Hash(),
			"proposer", pending.Checkpoint.Proposer, "start", pending.Checkpoint.Start, "end", pending.Checkpoint.End)
		return
	}

	latest, err := p.getLatestCheckpoint()
	if err != nil {
		log.Error("Failed to get latest checkpoint", "number", block.Number(), "hash", block.Hash(), "err", err)
		return
	}

	var end uint64 = 1
	if latest != nil {
		end = latest.End.Uint64() + 1
	}
	blockEpoch := xutil.CalcBlocksEachEpoch()
	if block.NumberU64()-end < blockEpoch {
		log.Debug("Current block hight not arrive checkpoint epoch", "number", block.NumberU64(), "blockEpoch", blockEpoch)
		return
	}

	proposal := &Checkpoint{
		Proposer: common.Address(currentValidator.Address),
		Start:    big.NewInt(0).SetUint64(end),
		End:      big.NewInt(0).SetUint64(latest.End.Uint64() + blockEpoch),
		ChainId:  p.chainId,
	}

	rootHash, err := p.rootHash(proposal.Start.Uint64(), proposal.End.Uint64())
	if err != nil {
		log.Error("Failed to get root hash", "start", proposal.Start, "end", proposal.End, "err", err)
		return
	}

	// TODO: Should submit checkpoint?
}

func (p *CheckpointProcessor) confirm(log *types.Log) {}

func (p *CheckpointProcessor) pendingCheckpoint() (*PendingCheckpoint, error) {
	blockNr := rpc.BlockNumber(rpc.PendingBlockNumber)

	const method = "pendingCheckpoint"

	data, err := p.checkpointABI.Pack(method)
	if err != nil {
		return nil, err
	}

	msgData := (hexutil.Bytes)(data)
	toAddress := cvm.CheckpointSigAggAddr
	gas := (hexutil.Uint64)(uint64(math.MaxUint64 / 2))

	result, err := p.caller.Call(context.Background(), ethapi.CallArgs{
		To:   &toAddress,
		Data: &msgData,
		Gas:  &gas,
	}, rpc.BlockNumberOrHash{BlockNumber: &blockNr}, nil)
	if err != nil {
		return nil, err
	}

	pending := new(checkpoint.ICheckpointSigAggregatorPendingCheckpoint)
	if err := p.checkpointABI.UnpackIntoInterface(pending, method, result); err != nil {
		return nil, err
	}
	return (*PendingCheckpoint)(pending), nil
}

func (p *CheckpointProcessor) getLatestCheckpoint() (*Checkpoint, error) {
	blockNr := rpc.BlockNumber(rpc.PendingBlockNumber)

	const method = "latestCheckpoint"

	data, err := p.checkpointABI.Pack(method)
	if err != nil {
		return nil, err
	}

	msgData := (hexutil.Bytes)(data)
	toAddress := cvm.CheckpointSigAggAddr
	gas := (hexutil.Uint64)(uint64(math.MaxUint64 / 2))

	result, err := p.caller.Call(context.Background(), ethapi.CallArgs{
		To:   &toAddress,
		Data: &msgData,
		Gas:  &gas,
	}, rpc.BlockNumberOrHash{BlockNumber: &blockNr}, nil)
	if err != nil {
		return nil, err
	}

	cp := new(checkpoint.ICheckpointSigAggregatorCheckpoint)
	if err := p.checkpointABI.UnpackIntoInterface(cp, method, result); err != nil {
		return nil, err
	}
	return (*Checkpoint)(cp), nil
}

func (p *CheckpointProcessor) rootHash(start, end uint64) (common.Hash, error) {
	key := getRootHashKey(start, end)

	if root, known := p.rootHashCache.Get(key); known {
		return common.BytesToHash(root), nil
	}

	length := end - start + 1
	currentHeaderNumber := p.chain.CurrentHeader().Number.Uint64()

	if start > end || end > currentHeaderNumber {
		return common.ZeroHash, fmt.Errorf("invalid start end block(start: %d, end: %d, current: %d)", start, end, currentHeaderNumber)
	}

	blockHeaders := make([]*types.Header, end-start+1)
	wg := new(sync.WaitGroup)
	concurrent := make(chan bool, 20)

	for i := start; i <= end; i++ {
		wg.Add(1)
		concurrent <- true

		go func(number uint64) {
			blockHeaders[number-start] = p.chain.GetHeaderByNumber(number)

			<-concurrent
			wg.Done()
		}(i)
	}
	wg.Wait()
	close(concurrent)

	headers := make([][32]byte, nextPowerOfTwo(length))

	for i := 0; i < len(blockHeaders); i++ {
		blockHeader := blockHeaders[i]
		header := crypto.Keccak256(appendBytes32(
			blockHeader.Number.Bytes(),
			new(big.Int).SetUint64(blockHeader.Time).Bytes(),
			blockHeader.TxHash.Bytes(),
			blockHeader.ReceiptHash.Bytes(),
		))

		var arr [32]byte
		copy(arr[:], header)
		headers[i] = arr
	}

	tree := merkle.NewTreeWithOpts(merkle.TreeOptions{EnableHashSorting: false, DisableHashLeves: true})
	if err := tree.Generate(convert(headers), sha3.blockHeaders256()); err != nil {
		return common.ZeroHash, err
	}

	p.rootHashCache.Add(key, tree.Root().Hash)

	return common.BytesToHash(tree.Root().Hash), nil
}

func getRootHashKey(start, end uint64) string {
	return strconv.FormatUint(start, 10) + "-" + strconv.FormatUint(end, 10)
}
