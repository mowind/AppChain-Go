package helper

import (
	"bytes"
	"github.com/PlatONnetwork/AppChain-Go/core/types"
	"math/big"
)

func EncodeStakeStateSync(blockNumber *big.Int, logs []*types.Log) ([]byte, error) {
	var events [][]byte
	for _, log := range logs {
		buffer := bytes.NewBuffer([]byte{})
		if err := log.EncodeRLP(buffer); err != nil {
			return nil, err
		}
		events = append(events, buffer.Bytes())
	}
	bytes, err := InnerStakeAbi.Pack(StakeStateSync, blockNumber, events)
	return bytes, err
}
