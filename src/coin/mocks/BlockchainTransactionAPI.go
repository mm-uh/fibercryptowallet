// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import core "github.com/fibercrypto/fibercryptowallet/src/core"
import mock "github.com/stretchr/testify/mock"

// BlockchainTransactionAPI is an autogenerated mock type for the BlockchainTransactionAPI type
type BlockchainTransactionAPI struct {
	mock.Mock
}

// SendFromAddress provides a mock function with given fields: from, to, change, options
func (_m *BlockchainTransactionAPI) SendFromAddress(from []core.WalletAddress, to []core.TransactionOutput, change core.Address, options core.KeyValueStore) (core.Transaction, error) {
	ret := _m.Called(from, to, change, options)

	var r0 core.Transaction
	if rf, ok := ret.Get(0).(func([]core.WalletAddress, []core.TransactionOutput, core.Address, core.KeyValueStore) core.Transaction); ok {
		r0 = rf(from, to, change, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(core.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]core.WalletAddress, []core.TransactionOutput, core.Address, core.KeyValueStore) error); ok {
		r1 = rf(from, to, change, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Spend provides a mock function with given fields: unspent, new, change, options
func (_m *BlockchainTransactionAPI) Spend(unspent []core.WalletOutput, new []core.TransactionOutput, change core.Address, options core.KeyValueStore) (core.Transaction, error) {
	ret := _m.Called(unspent, new, change, options)

	var r0 core.Transaction
	if rf, ok := ret.Get(0).(func([]core.WalletOutput, []core.TransactionOutput, core.Address, core.KeyValueStore) core.Transaction); ok {
		r0 = rf(unspent, new, change, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(core.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]core.WalletOutput, []core.TransactionOutput, core.Address, core.KeyValueStore) error); ok {
		r1 = rf(unspent, new, change, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}