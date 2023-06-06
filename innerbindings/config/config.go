package config

import (
	"github.com/PlatONnetwork/AppChain-Go/common"
	"math/big"
)

type RootChainContractConfig struct {
	PlatonRPCAddr      string         `json:"rpcAddress"`
	RootChainID        *big.Int       `json:"rootChainId"`
	StakingInfoAddress common.Address `json:"stakingInfoAddress"`
	RootChainAddress   common.Address `json:"rootChainAddress"`
}
