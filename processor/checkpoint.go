package processor

import (
	"math/big"

	"github.com/PlatONnetwork/AppChain-Go/core/cbfttypes"
	"github.com/PlatONnetwork/AppChain-Go/core/types"
	"github.com/PlatONnetwork/AppChain-Go/event"
	"github.com/PlatONnetwork/AppChain-Go/log"
)

type TxSigner interface {
	Nonce() uint64
	Sign(tx *types.Transaction, chainId *big.Int) (*types.Transaction, error)
}

type CheckpointProcessor struct {
	managerAccount TxSigner
	newHeadBlockCh chan *types.Log
	bftResultSub   *event.TypeMuxSubscription

	exitCh chan struct{}
}

func NewCheckpointProcessor(
	magnerAccount TxSigner,
	newHeadBlockCh chan *types.Log,
	bftResultSub *event.TypeMuxSubscription,
) *CheckpointProcessor {
	p := &CheckpointProcessor{
		managerAccount: magnerAccount,
		newHeadBlockCh: newHeadBlockCh,
		bftResultSub:   bftResultSub,
		exitCh:          make(chan struct{}),
	}

	go p.loop()

	return p
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
		case <- p.exitCh:
			log.Info("Checkpoint processor stopping...")
			return
		}
	}
}

func (p *CheckpointProcessor) handleNewHeadBlock(log *types.Log){
	// TODO: check if need to confirm checkpoint
	// TODO: submit confirm checkpoint transaction
}

func (p *CheckpointProcessor) handleBlock(block *types.Block) {
	// TODO: check if need to propose checkpoint
	// TODO: propose checkpoint
}
