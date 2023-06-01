// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package checkpoint

import (
	"errors"
	"math/big"
	"strings"

	hskchain "github.com/PlatONnetwork/AppChain-Go"
	"github.com/PlatONnetwork/AppChain-Go/accounts/abi"
	"github.com/PlatONnetwork/AppChain-Go/accounts/abi/bind"
	"github.com/PlatONnetwork/AppChain-Go/common"
	"github.com/PlatONnetwork/AppChain-Go/core/types"
	"github.com/PlatONnetwork/AppChain-Go/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = hskchain.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ICheckpointSigAggregatorCheckpoint is an auto generated low-level Go binding around an user-defined struct.
type ICheckpointSigAggregatorCheckpoint struct {
	Proposer    common.Address
	Start       *big.Int
	End         *big.Int
	RootHash    [32]byte
	AccountHash [32]byte
	ChainId     *big.Int
	Current     []uint32
	Rewards     []uint32
	Slashing    []uint32
}

// CheckpointMetaData contains all meta data concerning the Checkpoint contract.
var CheckpointMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"end\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32[]\",\"name\":\"signedValidators\",\"type\":\"uint32[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"CheckpointSigAggregated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proposer\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"confirm\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestCheckpoint\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"proposer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"end\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"rootHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"accountHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"uint32[]\",\"name\":\"current\",\"type\":\"uint32[]\"},{\"internalType\":\"uint32[]\",\"name\":\"rewards\",\"type\":\"uint32[]\"},{\"internalType\":\"uint32[]\",\"name\":\"slashing\",\"type\":\"uint32[]\"}],\"internalType\":\"structICheckpointSigAggregator.Checkpoint\",\"name\":\"cp\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"proposer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"end\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"rootHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"accountHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"uint32[]\",\"name\":\"current\",\"type\":\"uint32[]\"},{\"internalType\":\"uint32[]\",\"name\":\"rewards\",\"type\":\"uint32[]\"},{\"internalType\":\"uint32[]\",\"name\":\"slashing\",\"type\":\"uint32[]\"}],\"internalType\":\"structICheckpointSigAggregator.Checkpoint\",\"name\":\"cp\",\"type\":\"tuple\"},{\"internalType\":\"uint32\",\"name\":\"validatorId\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"propose\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// CheckpointABI is the input ABI used to generate the binding from.
// Deprecated: Use CheckpointMetaData.ABI instead.
var CheckpointABI = CheckpointMetaData.ABI

// Checkpoint is an auto generated Go binding around an hskchain contract.
type Checkpoint struct {
	CheckpointCaller     // Read-only binding to the contract
	CheckpointTransactor // Write-only binding to the contract
	CheckpointFilterer   // Log filterer for contract events
}

// CheckpointCaller is an auto generated read-only Go binding around an hskchain contract.
type CheckpointCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CheckpointTransactor is an auto generated write-only Go binding around an hskchain contract.
type CheckpointTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CheckpointFilterer is an auto generated log filtering Go binding around an hskchain contract events.
type CheckpointFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CheckpointSession is an auto generated Go binding around an hskchain contract,
// with pre-set call and transact options.
type CheckpointSession struct {
	Contract     *Checkpoint       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CheckpointCallerSession is an auto generated read-only Go binding around an hskchain contract,
// with pre-set call options.
type CheckpointCallerSession struct {
	Contract *CheckpointCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// CheckpointTransactorSession is an auto generated write-only Go binding around an hskchain contract,
// with pre-set transact options.
type CheckpointTransactorSession struct {
	Contract     *CheckpointTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// CheckpointRaw is an auto generated low-level Go binding around an hskchain contract.
type CheckpointRaw struct {
	Contract *Checkpoint // Generic contract binding to access the raw methods on
}

// CheckpointCallerRaw is an auto generated low-level read-only Go binding around an hskchain contract.
type CheckpointCallerRaw struct {
	Contract *CheckpointCaller // Generic read-only contract binding to access the raw methods on
}

// CheckpointTransactorRaw is an auto generated low-level write-only Go binding around an hskchain contract.
type CheckpointTransactorRaw struct {
	Contract *CheckpointTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCheckpoint creates a new instance of Checkpoint, bound to a specific deployed contract.
func NewCheckpoint(address common.Address, backend bind.ContractBackend) (*Checkpoint, error) {
	contract, err := bindCheckpoint(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Checkpoint{CheckpointCaller: CheckpointCaller{contract: contract}, CheckpointTransactor: CheckpointTransactor{contract: contract}, CheckpointFilterer: CheckpointFilterer{contract: contract}}, nil
}

// NewCheckpointCaller creates a new read-only instance of Checkpoint, bound to a specific deployed contract.
func NewCheckpointCaller(address common.Address, caller bind.ContractCaller) (*CheckpointCaller, error) {
	contract, err := bindCheckpoint(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CheckpointCaller{contract: contract}, nil
}

// NewCheckpointTransactor creates a new write-only instance of Checkpoint, bound to a specific deployed contract.
func NewCheckpointTransactor(address common.Address, transactor bind.ContractTransactor) (*CheckpointTransactor, error) {
	contract, err := bindCheckpoint(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CheckpointTransactor{contract: contract}, nil
}

// NewCheckpointFilterer creates a new log filterer instance of Checkpoint, bound to a specific deployed contract.
func NewCheckpointFilterer(address common.Address, filterer bind.ContractFilterer) (*CheckpointFilterer, error) {
	contract, err := bindCheckpoint(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CheckpointFilterer{contract: contract}, nil
}

// bindCheckpoint binds a generic wrapper to an already deployed contract.
func bindCheckpoint(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CheckpointABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Checkpoint *CheckpointRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Checkpoint.Contract.CheckpointCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Checkpoint *CheckpointRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Checkpoint.Contract.CheckpointTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Checkpoint *CheckpointRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Checkpoint.Contract.CheckpointTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Checkpoint *CheckpointCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Checkpoint.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Checkpoint *CheckpointTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Checkpoint.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Checkpoint *CheckpointTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Checkpoint.Contract.contract.Transact(opts, method, params...)
}

// LatestCheckpoint is a free data retrieval call binding the contract method 0x907c0f92.
//
// Solidity: function latestCheckpoint() view returns((address,uint256,uint256,bytes32,bytes32,uint256,uint32[],uint32[],uint32[]) cp)
func (_Checkpoint *CheckpointCaller) LatestCheckpoint(opts *bind.CallOpts) (ICheckpointSigAggregatorCheckpoint, error) {
	var out []interface{}
	err := _Checkpoint.contract.Call(opts, &out, "latestCheckpoint")

	if err != nil {
		return *new(ICheckpointSigAggregatorCheckpoint), err
	}

	out0 := *abi.ConvertType(out[0], new(ICheckpointSigAggregatorCheckpoint)).(*ICheckpointSigAggregatorCheckpoint)

	return out0, err

}

// LatestCheckpoint is a free data retrieval call binding the contract method 0x907c0f92.
//
// Solidity: function latestCheckpoint() view returns((address,uint256,uint256,bytes32,bytes32,uint256,uint32[],uint32[],uint32[]) cp)
func (_Checkpoint *CheckpointSession) LatestCheckpoint() (ICheckpointSigAggregatorCheckpoint, error) {
	return _Checkpoint.Contract.LatestCheckpoint(&_Checkpoint.CallOpts)
}

// LatestCheckpoint is a free data retrieval call binding the contract method 0x907c0f92.
//
// Solidity: function latestCheckpoint() view returns((address,uint256,uint256,bytes32,bytes32,uint256,uint32[],uint32[],uint32[]) cp)
func (_Checkpoint *CheckpointCallerSession) LatestCheckpoint() (ICheckpointSigAggregatorCheckpoint, error) {
	return _Checkpoint.Contract.LatestCheckpoint(&_Checkpoint.CallOpts)
}

// Confirm is a paid mutator transaction binding the contract method 0xce157608.
//
// Solidity: function confirm(address proposer, bytes32 root) returns()
func (_Checkpoint *CheckpointTransactor) Confirm(opts *bind.TransactOpts, proposer common.Address, root [32]byte) (*types.Transaction, error) {
	return _Checkpoint.contract.Transact(opts, "confirm", proposer, root)
}

// Confirm is a paid mutator transaction binding the contract method 0xce157608.
//
// Solidity: function confirm(address proposer, bytes32 root) returns()
func (_Checkpoint *CheckpointSession) Confirm(proposer common.Address, root [32]byte) (*types.Transaction, error) {
	return _Checkpoint.Contract.Confirm(&_Checkpoint.TransactOpts, proposer, root)
}

// Confirm is a paid mutator transaction binding the contract method 0xce157608.
//
// Solidity: function confirm(address proposer, bytes32 root) returns()
func (_Checkpoint *CheckpointTransactorSession) Confirm(proposer common.Address, root [32]byte) (*types.Transaction, error) {
	return _Checkpoint.Contract.Confirm(&_Checkpoint.TransactOpts, proposer, root)
}

// Propose is a paid mutator transaction binding the contract method 0x044d1cca.
//
// Solidity: function propose((address,uint256,uint256,bytes32,bytes32,uint256,uint32[],uint32[],uint32[]) cp, uint32 validatorId, bytes signature) returns()
func (_Checkpoint *CheckpointTransactor) Propose(opts *bind.TransactOpts, cp ICheckpointSigAggregatorCheckpoint, validatorId uint32, signature []byte) (*types.Transaction, error) {
	return _Checkpoint.contract.Transact(opts, "propose", cp, validatorId, signature)
}

// Propose is a paid mutator transaction binding the contract method 0x044d1cca.
//
// Solidity: function propose((address,uint256,uint256,bytes32,bytes32,uint256,uint32[],uint32[],uint32[]) cp, uint32 validatorId, bytes signature) returns()
func (_Checkpoint *CheckpointSession) Propose(cp ICheckpointSigAggregatorCheckpoint, validatorId uint32, signature []byte) (*types.Transaction, error) {
	return _Checkpoint.Contract.Propose(&_Checkpoint.TransactOpts, cp, validatorId, signature)
}

// Propose is a paid mutator transaction binding the contract method 0x044d1cca.
//
// Solidity: function propose((address,uint256,uint256,bytes32,bytes32,uint256,uint32[],uint32[],uint32[]) cp, uint32 validatorId, bytes signature) returns()
func (_Checkpoint *CheckpointTransactorSession) Propose(cp ICheckpointSigAggregatorCheckpoint, validatorId uint32, signature []byte) (*types.Transaction, error) {
	return _Checkpoint.Contract.Propose(&_Checkpoint.TransactOpts, cp, validatorId, signature)
}

// CheckpointCheckpointSigAggregatedIterator is returned from FilterCheckpointSigAggregated and is used to iterate over the raw logs and unpacked data for CheckpointSigAggregated events raised by the Checkpoint contract.
type CheckpointCheckpointSigAggregatedIterator struct {
	Event *CheckpointCheckpointSigAggregated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  hskchain.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CheckpointCheckpointSigAggregatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CheckpointCheckpointSigAggregated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CheckpointCheckpointSigAggregated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CheckpointCheckpointSigAggregatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CheckpointCheckpointSigAggregatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CheckpointCheckpointSigAggregated represents a CheckpointSigAggregated event raised by the Checkpoint contract.
type CheckpointCheckpointSigAggregated struct {
	Proposer         common.Address
	Start            *big.Int
	End              *big.Int
	Root             [32]byte
	SignedValidators []uint32
	Signature        []byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterCheckpointSigAggregated is a free log retrieval operation binding the contract event 0xf670678827823b741813a19384f6df652e8bf8e39c172da4666fe9616fb19a6c.
//
// Solidity: event CheckpointSigAggregated(address indexed proposer, uint256 start, uint256 end, bytes32 root, uint32[] signedValidators, bytes signature)
func (_Checkpoint *CheckpointFilterer) FilterCheckpointSigAggregated(opts *bind.FilterOpts, proposer []common.Address) (*CheckpointCheckpointSigAggregatedIterator, error) {

	var proposerRule []interface{}
	for _, proposerItem := range proposer {
		proposerRule = append(proposerRule, proposerItem)
	}

	logs, sub, err := _Checkpoint.contract.FilterLogs(opts, "CheckpointSigAggregated", proposerRule)
	if err != nil {
		return nil, err
	}
	return &CheckpointCheckpointSigAggregatedIterator{contract: _Checkpoint.contract, event: "CheckpointSigAggregated", logs: logs, sub: sub}, nil
}

// WatchCheckpointSigAggregated is a free log subscription operation binding the contract event 0xf670678827823b741813a19384f6df652e8bf8e39c172da4666fe9616fb19a6c.
//
// Solidity: event CheckpointSigAggregated(address indexed proposer, uint256 start, uint256 end, bytes32 root, uint32[] signedValidators, bytes signature)
func (_Checkpoint *CheckpointFilterer) WatchCheckpointSigAggregated(opts *bind.WatchOpts, sink chan<- *CheckpointCheckpointSigAggregated, proposer []common.Address) (event.Subscription, error) {

	var proposerRule []interface{}
	for _, proposerItem := range proposer {
		proposerRule = append(proposerRule, proposerItem)
	}

	logs, sub, err := _Checkpoint.contract.WatchLogs(opts, "CheckpointSigAggregated", proposerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CheckpointCheckpointSigAggregated)
				if err := _Checkpoint.contract.UnpackLog(event, "CheckpointSigAggregated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCheckpointSigAggregated is a log parse operation binding the contract event 0xf670678827823b741813a19384f6df652e8bf8e39c172da4666fe9616fb19a6c.
//
// Solidity: event CheckpointSigAggregated(address indexed proposer, uint256 start, uint256 end, bytes32 root, uint32[] signedValidators, bytes signature)
func (_Checkpoint *CheckpointFilterer) ParseCheckpointSigAggregated(log types.Log) (*CheckpointCheckpointSigAggregated, error) {
	event := new(CheckpointCheckpointSigAggregated)
	if err := _Checkpoint.contract.UnpackLog(event, "CheckpointSigAggregated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
