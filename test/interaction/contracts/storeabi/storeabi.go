// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package storeabi

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// StoreabiMetaData contains all meta data concerning the Storeabi contract.
var StoreabiMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_version\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"value\",\"type\":\"bytes32\"}],\"name\":\"ItemSet\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"items\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"value\",\"type\":\"bytes32\"}],\"name\":\"setItem\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// StoreabiABI is the input ABI used to generate the binding from.
// Deprecated: Use StoreabiMetaData.ABI instead.
var StoreabiABI = StoreabiMetaData.ABI

// Storeabi is an auto generated Go binding around an Ethereum contract.
type Storeabi struct {
	StoreabiCaller     // Read-only binding to the contract
	StoreabiTransactor // Write-only binding to the contract
	StoreabiFilterer   // Log filterer for contract events
}

// StoreabiCaller is an auto generated read-only Go binding around an Ethereum contract.
type StoreabiCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StoreabiTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StoreabiTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StoreabiFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StoreabiFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StoreabiSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StoreabiSession struct {
	Contract     *Storeabi         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StoreabiCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StoreabiCallerSession struct {
	Contract *StoreabiCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// StoreabiTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StoreabiTransactorSession struct {
	Contract     *StoreabiTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// StoreabiRaw is an auto generated low-level Go binding around an Ethereum contract.
type StoreabiRaw struct {
	Contract *Storeabi // Generic contract binding to access the raw methods on
}

// StoreabiCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StoreabiCallerRaw struct {
	Contract *StoreabiCaller // Generic read-only contract binding to access the raw methods on
}

// StoreabiTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StoreabiTransactorRaw struct {
	Contract *StoreabiTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStoreabi creates a new instance of Storeabi, bound to a specific deployed contract.
func NewStoreabi(address common.Address, backend bind.ContractBackend) (*Storeabi, error) {
	contract, err := bindStoreabi(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Storeabi{StoreabiCaller: StoreabiCaller{contract: contract}, StoreabiTransactor: StoreabiTransactor{contract: contract}, StoreabiFilterer: StoreabiFilterer{contract: contract}}, nil
}

// NewStoreabiCaller creates a new read-only instance of Storeabi, bound to a specific deployed contract.
func NewStoreabiCaller(address common.Address, caller bind.ContractCaller) (*StoreabiCaller, error) {
	contract, err := bindStoreabi(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StoreabiCaller{contract: contract}, nil
}

// NewStoreabiTransactor creates a new write-only instance of Storeabi, bound to a specific deployed contract.
func NewStoreabiTransactor(address common.Address, transactor bind.ContractTransactor) (*StoreabiTransactor, error) {
	contract, err := bindStoreabi(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StoreabiTransactor{contract: contract}, nil
}

// NewStoreabiFilterer creates a new log filterer instance of Storeabi, bound to a specific deployed contract.
func NewStoreabiFilterer(address common.Address, filterer bind.ContractFilterer) (*StoreabiFilterer, error) {
	contract, err := bindStoreabi(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StoreabiFilterer{contract: contract}, nil
}

// bindStoreabi binds a generic wrapper to an already deployed contract.
func bindStoreabi(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StoreabiMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Storeabi *StoreabiRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Storeabi.Contract.StoreabiCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Storeabi *StoreabiRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Storeabi.Contract.StoreabiTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Storeabi *StoreabiRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Storeabi.Contract.StoreabiTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Storeabi *StoreabiCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Storeabi.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Storeabi *StoreabiTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Storeabi.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Storeabi *StoreabiTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Storeabi.Contract.contract.Transact(opts, method, params...)
}

// Items is a free data retrieval call binding the contract method 0x48f343f3.
//
// Solidity: function items(bytes32 ) view returns(bytes32)
func (_Storeabi *StoreabiCaller) Items(opts *bind.CallOpts, arg0 [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Storeabi.contract.Call(opts, &out, "items", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Items is a free data retrieval call binding the contract method 0x48f343f3.
//
// Solidity: function items(bytes32 ) view returns(bytes32)
func (_Storeabi *StoreabiSession) Items(arg0 [32]byte) ([32]byte, error) {
	return _Storeabi.Contract.Items(&_Storeabi.CallOpts, arg0)
}

// Items is a free data retrieval call binding the contract method 0x48f343f3.
//
// Solidity: function items(bytes32 ) view returns(bytes32)
func (_Storeabi *StoreabiCallerSession) Items(arg0 [32]byte) ([32]byte, error) {
	return _Storeabi.Contract.Items(&_Storeabi.CallOpts, arg0)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_Storeabi *StoreabiCaller) Version(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Storeabi.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_Storeabi *StoreabiSession) Version() (string, error) {
	return _Storeabi.Contract.Version(&_Storeabi.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_Storeabi *StoreabiCallerSession) Version() (string, error) {
	return _Storeabi.Contract.Version(&_Storeabi.CallOpts)
}

// SetItem is a paid mutator transaction binding the contract method 0xf56256c7.
//
// Solidity: function setItem(bytes32 key, bytes32 value) returns()
func (_Storeabi *StoreabiTransactor) SetItem(opts *bind.TransactOpts, key [32]byte, value [32]byte) (*types.Transaction, error) {
	return _Storeabi.contract.Transact(opts, "setItem", key, value)
}

// SetItem is a paid mutator transaction binding the contract method 0xf56256c7.
//
// Solidity: function setItem(bytes32 key, bytes32 value) returns()
func (_Storeabi *StoreabiSession) SetItem(key [32]byte, value [32]byte) (*types.Transaction, error) {
	return _Storeabi.Contract.SetItem(&_Storeabi.TransactOpts, key, value)
}

// SetItem is a paid mutator transaction binding the contract method 0xf56256c7.
//
// Solidity: function setItem(bytes32 key, bytes32 value) returns()
func (_Storeabi *StoreabiTransactorSession) SetItem(key [32]byte, value [32]byte) (*types.Transaction, error) {
	return _Storeabi.Contract.SetItem(&_Storeabi.TransactOpts, key, value)
}

// StoreabiItemSetIterator is returned from FilterItemSet and is used to iterate over the raw logs and unpacked data for ItemSet events raised by the Storeabi contract.
type StoreabiItemSetIterator struct {
	Event *StoreabiItemSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StoreabiItemSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StoreabiItemSet)
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
		it.Event = new(StoreabiItemSet)
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
func (it *StoreabiItemSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StoreabiItemSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StoreabiItemSet represents a ItemSet event raised by the Storeabi contract.
type StoreabiItemSet struct {
	Key   [32]byte
	Value [32]byte
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterItemSet is a free log retrieval operation binding the contract event 0xe79e73da417710ae99aa2088575580a60415d359acfad9cdd3382d59c80281d4.
//
// Solidity: event ItemSet(bytes32 key, bytes32 value)
func (_Storeabi *StoreabiFilterer) FilterItemSet(opts *bind.FilterOpts) (*StoreabiItemSetIterator, error) {

	logs, sub, err := _Storeabi.contract.FilterLogs(opts, "ItemSet")
	if err != nil {
		return nil, err
	}
	return &StoreabiItemSetIterator{contract: _Storeabi.contract, event: "ItemSet", logs: logs, sub: sub}, nil
}

// WatchItemSet is a free log subscription operation binding the contract event 0xe79e73da417710ae99aa2088575580a60415d359acfad9cdd3382d59c80281d4.
//
// Solidity: event ItemSet(bytes32 key, bytes32 value)
func (_Storeabi *StoreabiFilterer) WatchItemSet(opts *bind.WatchOpts, sink chan<- *StoreabiItemSet) (event.Subscription, error) {

	logs, sub, err := _Storeabi.contract.WatchLogs(opts, "ItemSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StoreabiItemSet)
				if err := _Storeabi.contract.UnpackLog(event, "ItemSet", log); err != nil {
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

// ParseItemSet is a log parse operation binding the contract event 0xe79e73da417710ae99aa2088575580a60415d359acfad9cdd3382d59c80281d4.
//
// Solidity: event ItemSet(bytes32 key, bytes32 value)
func (_Storeabi *StoreabiFilterer) ParseItemSet(log types.Log) (*StoreabiItemSet, error) {
	event := new(StoreabiItemSet)
	if err := _Storeabi.contract.UnpackLog(event, "ItemSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
