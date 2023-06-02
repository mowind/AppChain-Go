package vm

import (
	"math/big"
	"testing"

	"github.com/PlatONnetwork/AppChain-Go/accounts/abi"
	"github.com/PlatONnetwork/AppChain-Go/common"
	"github.com/PlatONnetwork/AppChain-Go/common/mock"
	"github.com/PlatONnetwork/AppChain-Go/common/vm"
	"github.com/PlatONnetwork/AppChain-Go/core/vm/solidity"
	"github.com/PlatONnetwork/AppChain-Go/core/vm/solidity/checkpoint"
	"github.com/PlatONnetwork/AppChain-Go/core/vm/solidity/types"
	"github.com/PlatONnetwork/AppChain-Go/crypto"
	"github.com/PlatONnetwork/AppChain-Go/x/xutil"
	"github.com/stretchr/testify/assert"
)

func makePropose() *checkpoint.ICheckpointSigAggregatorCheckpoint {
	return &checkpoint.ICheckpointSigAggregatorCheckpoint{
		Proposer:    common.HexToAddress("0x134234"),
		Start:       big.NewInt(1),
		End:         big.NewInt(100),
		ChainId:     big.NewInt(101),
		RootHash:    common.HexToHash("0x13421423134"),
		AccountHash: common.HexToHash("0x14321421432"),
		Current:     []uint32{1, 2, 3},
		Rewards:     []uint32{1, 2, 3, 5, 6, 7},
		Slashing:    []uint32{4},
	}
}

func TestPropose_fail(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()
	newPlugins()

	contract := newContract(big.NewInt(10000), common.HexToAddress("0x1234"))
	c := &CheckpointSigAggregatorContract{
		Contract: contract,
		Evm:      newEvm(big1, blockHash, chain),
	}
	cp := makePropose()

	// Validator not exists
	var validatorId uint32 = 99
	var signature []byte = []byte{0x1, 0x2, 0x3}
	input, err := CheckpointABI().Pack("propose", cp, validatorId, signature)
	assert.Nil(t, err)

	out, err := c.Run(input)
	assert.Error(t, err, "The validator does not exist")
	assert.True(t, len(out) == 0)

	build_staking_data(chain.SnapDB, chain.Genesis.Hash())
	chain.SnapDB.Commit(blockHash)

	// Validator not found
	out, err = c.Run(input)
	assert.Equal(t, err, ErrValidatorNotFound)

	// Invalid caller
	cp.Proposer = addrArr[0]
	validatorId = 0
	input, err = CheckpointABI().Pack("propose", cp, validatorId, signature)
	assert.Nil(t, err)
	_, err = c.Run(input)
	assert.Equal(t, err, ErrInvalidCaller)

	addr, _ := xutil.NodeId2Addr(nodeIdArr[0])
	correctCaller := common.Address(addr)

	// Invalid proposer
	cp.Proposer = common.HexToAddress("0x12")
	c.Contract.CallerAddress = correctCaller
	input, err = CheckpointABI().Pack("propose", cp, validatorId, signature)
	assert.Nil(t, err)
	_, err = c.Run(input)
	assert.Equal(t, err, ErrInvalidProposer)

	// Verify signature fail
	cp.Proposer = correctCaller
	input, err = CheckpointABI().Pack("propose", cp, validatorId, signature)
	assert.Nil(t, err)
	_, err = c.Run(input)
	assert.Equal(t, err, ErrVerifySigFail)

	// Read latest checkpoint fail
	chain.StateDB.SetState(vm.CheckpointSigAggAddr, latestCheckpointKey, []byte{1, 2, 3})

	tcp := solidity.ICheckpointToCheckpoint(cp)
	packed := tcp.Pack()
	sig := blsKey1.Sign(string(crypto.Keccak256(packed)))
	signature = sig.Serialize()
	input, err = CheckpointABI().Pack("propose", cp, validatorId, signature)
	assert.Nil(t, err)
	_, err = c.Run(input)
	assert.Error(t, err, "expected input list")

	// Invalid propose
	WriteLatestCheckpoint(chain.StateDB, &StorageCheckpoint{
		Checkpoint: types.Checkpoint{
			End: big.NewInt(10),
		},
	})
	_, err = c.Run(input)
	assert.Equal(t, err, ErrInvalidProposal)

	// Read pending checkpoint fail
	chain.StateDB.SetState(vm.CheckpointSigAggAddr, pendingCheckpointKey, []byte{1, 2, 3})
	cp.Start = big.NewInt(11)
	cp.End = big.NewInt(20)
	tcp = solidity.ICheckpointToCheckpoint(cp)
	packed = tcp.Pack()
	sig = blsKey1.Sign(string(crypto.Keccak256(packed)))
	signature = sig.Serialize()
	input, err = CheckpointABI().Pack("propose", cp, validatorId, signature)
	assert.Nil(t, err)
	_, err = c.Run(input)
	assert.Error(t, err, "expected input list")

	// Invalid popose(other proposer)
	WritePendingCheckpoint(chain.StateDB, &StorageCheckpoint{
		Checkpoint: *tcp,
		BlockNum:   0,
	})
	addr, _ = xutil.NodeId2Addr(nodeIdArr[1])
	cp.Proposer = common.Address(addr)
	validatorId = 1
	tcp = solidity.ICheckpointToCheckpoint(cp)
	packed = tcp.Pack()
	sig = blsKey2.Sign(string(crypto.Keccak256(packed)))
	signature = sig.Serialize()
	input, err = CheckpointABI().Pack("propose", cp, validatorId, signature)
	assert.Nil(t, err)
	c.Contract.CallerAddress = common.Address(addr)
	_, err = c.Run(input)
	assert.Equal(t, err, ErrInvalidProposal)

	// Propose timeout
	cp.Proposer = correctCaller
	validatorId = 0
	c.Contract.CallerAddress = correctCaller
	c.Evm.Context.BlockNumber = big.NewInt(11)
	tcp = solidity.ICheckpointToCheckpoint(cp)
	packed = tcp.Pack()
	sig = blsKey1.Sign(string(crypto.Keccak256(packed)))
	signature = sig.Serialize()
	input, err = CheckpointABI().Pack("propose", cp, validatorId, signature)
	assert.Nil(t, err)
	_, err = c.Run(input)
	assert.Equal(t, err, ErrProposalTimeout)

	// Emitted
	WritePendingCheckpoint(chain.StateDB, &StorageCheckpoint{
		Checkpoint: *tcp,
		Emitted:    true,
	})
	c.Evm.Context.BlockNumber = big.NewInt(1)
	_, err = c.Run(input)
	assert.Equal(t, err, ErrEmitted)

	// Checkpoint not equal
	cp.End = big.NewInt(13)
	tcp = solidity.ICheckpointToCheckpoint(cp)
	packed = tcp.Pack()
	sig = blsKey1.Sign(string(crypto.Keccak256(packed)))
	signature = sig.Serialize()
	input, err = CheckpointABI().Pack("propose", cp, validatorId, signature)
	assert.Nil(t, err)
	_, err = c.Run(input)
	assert.Equal(t, err, ErrInvalidProposal)
}

func TestPropose(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	build_staking_data(chain.SnapDB, chain.Genesis.Hash())
	chain.SnapDB.Commit(blockHash)

	addr, _ := xutil.NodeId2Addr(nodeIdArr[0])
	sender := common.Address(addr)

	evm := newEvm(big1, blockHash, chain)
	contract := newContract(big.NewInt(10000), sender)

	cp := makePropose()
	cp.Proposer = sender

	var validatorId uint32
	tcp := solidity.ICheckpointToCheckpoint(cp)
	packed := tcp.Pack()
	sig := blsKey1.Sign(string(crypto.Keccak256(packed)))
	signature := sig.Serialize()

	data, err := CheckpointABI().Pack("propose", cp, validatorId, signature)
	assert.Nil(t, err)

	c := &CheckpointSigAggregatorContract{
		Contract: contract,
		Evm: evm,
	}

	_, err = c.Run(data)
	assert.Nil(t, err)
}

func TestConfirm(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	build_staking_data(chain.SnapDB, chain.Genesis.Hash())
	chain.SnapDB.Commit(blockHash)

	addr, _ := xutil.NodeId2Addr(nodeIdArr[0])
	sender := common.Address(addr)

	evm := newEvm(big1, blockHash, chain)
	contract := newContract(big.NewInt(10000), sender)


	c := &CheckpointSigAggregatorContract{
		Contract: contract,
		Evm: evm,
	}

	scp := &StorageCheckpoint{
		Checkpoint: types.Checkpoint{
			Proposer:    sender,
			Start:       big.NewInt(1),
			End:         big.NewInt(100),
			ChainId:     big.NewInt(101),
			RootHash:    common.HexToHash("0x13421423134"),
			AccountHash: common.HexToHash("0x14321421432"),
			Current:     []uint32{1, 2, 3},
			Rewards:     []uint32{1, 2, 3, 5, 6, 7},
			Slashing:    []uint32{4},
		},
	}

	WritePendingCheckpoint(chain.StateDB, scp)

	input, err := CheckpointABI().Pack("confirm", scp.Proposer, scp.RootHash)
	assert.Nil(t, err)

	out, err := c.Run(input)
	assert.Nil(t, err)
	assert.True(t, len(out) == 0)

	scp1, err := ReadPendingCheckpoint(chain.StateDB)
	assert.Equal(t, err, ErrCheckpointNotFound)
	assert.Nil(t, scp1)

	scp2, err := ReadLatestCheckpoint(chain.StateDB)
	assert.Nil(t, err)
	assert.Equal(t, scp.Checkpoint, scp2.Checkpoint)
}

func TestLatestCheckpoint(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	c := &CheckpointSigAggregatorContract{
		Evm: newEvm(big.NewInt(1), common.HexToHash("0x13412412343"), chain),
	}

	input, err := CheckpointABI().Pack("latestCheckpoint")
	assert.Nil(t, err)

	out, err := c.Run(input)
	assert.Nil(t, err)
	assert.True(t, len(out) == 0)

	scp := &StorageCheckpoint{
		Checkpoint: types.Checkpoint{
			Proposer:    common.HexToAddress("0x1234"),
			Start:       big.NewInt(1),
			End:         big.NewInt(2),
			ChainId:     big.NewInt(5),
			RootHash:    common.HexToHash("0x12342314"),
			AccountHash: common.HexToHash("0x1342134"),
			Current:     []uint32{1, 2, 3},
			Rewards:     []uint32{1, 2, 3, 8, 9},
			Slashing:    []uint32{4, 5, 6},
		},
	}

	WriteLatestCheckpoint(chain.StateDB, scp)

	out, err = c.Run(input)
	assert.Nil(t, err)

	var cp checkpoint.ICheckpointSigAggregatorCheckpoint
	out0, err := CheckpointABI().Unpack("latestCheckpoint", out)
	assert.Nil(t, err)
	abi.ConvertType(out0[0], &cp)
	assert.Equal(t, cp.Proposer[:], scp.Proposer[:])
	assert.Equal(t, scp.Start, cp.Start)
	assert.Equal(t, scp.End, cp.End)
}
