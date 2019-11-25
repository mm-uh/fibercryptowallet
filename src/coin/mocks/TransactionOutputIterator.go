// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import core "github.com/fibercrypto/FiberCryptoWallet/src/core"
import mock "github.com/stretchr/testify/mock"

// TransactionOutputIterator is an autogenerated mock type for the TransactionOutputIterator type
type TransactionOutputIterator struct {
	mock.Mock
}

// HasNext provides a mock function with given fields:
func (_m *TransactionOutputIterator) HasNext() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Next provides a mock function with given fields:
func (_m *TransactionOutputIterator) Next() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Value provides a mock function with given fields:
func (_m *TransactionOutputIterator) Value() core.TransactionOutput {
	ret := _m.Called()

	var r0 core.TransactionOutput
	if rf, ok := ret.Get(0).(func() core.TransactionOutput); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(core.TransactionOutput)
		}
	}

	return r0
}