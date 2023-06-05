package types

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/PlatONnetwork/AppChain-Go/accounts/abi"
	"github.com/PlatONnetwork/AppChain-Go/common"
	"github.com/PlatONnetwork/AppChain-Go/common/hexutil"
)

var (
	// CheckpointPrimitives is the primitive ABI types for each Checkpoint field.
	CheckpointPrimitives = []abi.ArgumentMarshaling{
		{Name: "proposer", InternalType: "Proposer", Type: "address"},
		{Name: "start", InternalType: "Start", Type: "uint256"},
		{Name: "end", InternalType: "End", Type: "uint256"},
		{Name: "rootHash", InternalType: "RootHash", Type: "bytes32"},
		{Name: "accountHash", InternalType: "AccountHash", Type: "bytes32"},
		{Name: "chainId", InternalType: "ChainId", Type: "uint256"},
		{Name: "current", InternalType: "Current", Type: "uint32[]"},
		{Name: "rewards", InternalType: "Rewards", Type: "uint32[]"},
		{Name: "slashing", InternalType: "Slashing", Type: "uint32[]"},
	}

	// CheckpointType is the ABI type of a Checkpoint.
	CheckpointType, _ = abi.NewType("tuple", "checkpoint", CheckpointPrimitives)
)

func getAbiArgs() abi.Arguments {
	return abi.Arguments{
		{Name: "Checkpoint", Type: CheckpointType},
	}
}

// Checkpoint represents snapshots of the AppChain state and is supposed to be attested by 2/3+ of the
// validator set before it is validated and submitted on the contracts deployed on PlatON.
type Checkpoint struct {
	Proposer    common.Address
	Start       *big.Int
	End         *big.Int
	RootHash    common.Hash
	AccountHash common.Hash
	ChainId     *big.Int
	Current     []uint32
	Rewards     []uint32
	Slashing    []uint32
}

// Pack returns a standard message of the Checkpoint.
func (cp *Checkpoint) Pack() []byte {
	args := getAbiArgs()
	packed, _ := args.Pack(&struct {
		Proposer    common.Address
		Start       *big.Int
		End         *big.Int
		RootHash    common.Hash
		AccountHash common.Hash
		ChainId     *big.Int
		Current     []uint32
		Rewards     []uint32
		Slashing    []uint32
	}{
		cp.Proposer,
		cp.Start,
		cp.End,
		cp.RootHash,
		cp.AccountHash,
		cp.ChainId,
		cp.Current,
		cp.Rewards,
		cp.Slashing,
	})

	enc := hexutil.Encode(packed)
	enc = "0x" + enc[66:]
	return (hexutil.MustDecode(enc))

}

func (cp *Checkpoint) Equal(other *Checkpoint) bool {
	return reflect.DeepEqual(cp, other)
}
