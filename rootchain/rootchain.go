package rootchain

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PlatONnetwork/AppChain-Go/common"
	vm2 "github.com/PlatONnetwork/AppChain-Go/common/vm"
	"github.com/PlatONnetwork/AppChain-Go/core/state"
	"github.com/PlatONnetwork/AppChain-Go/core/types"
	"github.com/PlatONnetwork/AppChain-Go/core/vm"
	"github.com/PlatONnetwork/AppChain-Go/innerbindings/helper"
	"math/big"
)

type RootChainReader interface {
	//get root chain stake logs, return logs, end block number
	GetStakeLogs(start *big.Int, limit uint64) ([]*types.Log, *big.Int, error)
	//get logs base on block height range
	GetStakeLogsRange(start, end *big.Int) ([]*types.Log, error)
}

type RootChainCheck interface {
	CheckStakeStateSyncExtra(header *types.Header, tx *types.Transaction) error
}
type StateReader interface {
	StateAt(root common.Hash) (*state.StateDB, error)
}
type RootChain struct {
	stateReader StateReader
}

func NewRootChain(blockChainCache StateReader) (*RootChain, error) {
	return &RootChain{
		stateReader: blockChainCache,
	}, nil
}

func (r RootChain) CheckStakeStateSyncExtra(header *types.Header, tx *types.Transaction) error {
	end := types.DecodeStakeExtra(header.Extra)
	//if there is no stake tx in block, both tx and end is must empty
	if tx == nil {
		if end.Uint64() != 0 {
			return fmt.Errorf("stake tx is nil, but extra is not empty")
		}
		return nil
	}

	stateDb, err := r.stateReader.StateAt(header.ParentHash)
	if err != nil {
		return err
	}
	start := new(big.Int).SetBytes(stateDb.GetState(vm2.StakingContractAddr, vm.BlockNumberKey))
	logs, err := r.GetStakeLogsRange(big.NewInt(start.Int64()+1), end)
	if err != nil {
		return err
	}
	data, err := helper.EncodeStakeStateSync(end, logs)
	if err != nil {
		return err
	}
	if bytes.Equal(tx.Data(), data) {
		return errors.New(fmt.Sprintf("header extra check failed: expect tx data doesn't match actual tx data"))
	}
	return nil
}

func (r RootChain) GetStakeLogs(start *big.Int, limit uint64) ([]*types.Log, *big.Int, error) {
	//TODO implement me
	panic("implement me")
}

func (r RootChain) GetStakeLogsRange(start, end *big.Int) ([]*types.Log, error) {
	//TODO implement me
	panic("implement me")
}
