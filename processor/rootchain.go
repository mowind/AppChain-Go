package processor

import (
	"math/big"

	"github.com/PlatONnetwork/AppChain-Go/common"
	"github.com/PlatONnetwork/AppChain-Go/core/types"
	"github.com/PlatONnetwork/AppChain-Go/ethclient"
	"github.com/PlatONnetwork/AppChain-Go/innerbindings/rootchain"
	"github.com/PlatONnetwork/AppChain-Go/log"
)

type RootchainConnector struct {
	platonClient *ethclient.Client
	caller       *rootchain.RootchainCaller
	transactor   *rootchain.RootchainTransactor
}

func NewRootchainConnector(platonAddr string, contractAddr common.Address) (*RootchainConnector, error) {
	client, err := ethclient.Dial(platonAddr)
	if err != nil {
		return nil, err
	}
	caller, err := rootchain.NewRootchainCaller(contractAddr, client)
	if err != nil {
		return nil, err
	}
	transactor, err := rootchain.NewRootchainTransactor(contractAddr, client)
	if err != nil {
		return nil, err
	}

	return &RootchainConnector{
		platonClient: client,
		caller:       caller,
		transactor:   transactor,
	}, nil
}

func (c *RootchainConnector) CurrentHeaderBlock(blockInterval uint64) (uint64, error) {
	currentHeaderBlock, err := c.caller.CurrentHeaderBlock(nil)
	if err != nil {
		log.Error("Cannot fetch current header block from rootchain contract", "err", err)
		return 0, err
	}

	return currentHeaderBlock.Uint64() / blockInterval, nil
}

func (c *RootchainConnector) GetHeaderInfo(number, blockInterval uint64) (
	root common.Hash,
	start uint64,
	end uint64,
	createdAt uint64,
	proposer common.Address,
	err error,
) {
	checkpointBigInt := big.NewInt(0).Mul(big.NewInt(0).SetUint64(number), big.NewInt(0).SetUint64(blockInterval))

	headerBlock, err := c.caller.HeaderBlocks(nil, checkpointBigInt)
	if err != nil {
		return root, start, end, createdAt, proposer, err
	}

	return headerBlock.Root,
		headerBlock.Start.Uint64(),
		headerBlock.End.Uint64(),
		headerBlock.CreatedAt.Uint64(),
		headerBlock.Proposer,
		nil
}

func (c *RootchainConnector) GetLatestChildBlock() (uint64, error) {
	latestChildBlock, err := c.caller.GetLastChildBlock(nil)
	if err != nil{
		log.Error("Could not fetch current child block from rootchain contract", "err", err)
		return 0, err
	}

	return latestChildBlock.Uint64(), nil
}

type signFn func(*types.Transaction, *big.Int) (*types.Transaction, error)
func (c *RootchainConnector) SendCheckpoint(signedData []byte, signedValidators []*big.Int, signature []byte, signFn signFn) error {
	// TODO: submit checkpoint to rootchain
	return nil
}
