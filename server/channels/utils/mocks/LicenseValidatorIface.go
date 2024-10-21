// Code generated by mockery v2.42.2. DO NOT EDIT.

// Regenerate this file using `make misc-mocks`.

package mocks

import (
	model "git.biggo.com/Funmula/BigGoChat/server/public/model"
	mock "github.com/stretchr/testify/mock"
)

// LicenseValidatorIface is an autogenerated mock type for the LicenseValidatorIface type
type LicenseValidatorIface struct {
	mock.Mock
}

// LicenseFromBytes provides a mock function with given fields: licenseBytes
func (_m *LicenseValidatorIface) LicenseFromBytes(licenseBytes []byte) (*model.License, *model.AppError) {
	ret := _m.Called(licenseBytes)

	if len(ret) == 0 {
		panic("no return value specified for LicenseFromBytes")
	}

	var r0 *model.License
	var r1 *model.AppError
	if rf, ok := ret.Get(0).(func([]byte) (*model.License, *model.AppError)); ok {
		return rf(licenseBytes)
	}
	if rf, ok := ret.Get(0).(func([]byte) *model.License); ok {
		r0 = rf(licenseBytes)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.License)
		}
	}

	if rf, ok := ret.Get(1).(func([]byte) *model.AppError); ok {
		r1 = rf(licenseBytes)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// ValidateLicense provides a mock function with given fields: signed
func (_m *LicenseValidatorIface) ValidateLicense(signed []byte) (string, error) {
	ret := _m.Called(signed)

	if len(ret) == 0 {
		panic("no return value specified for ValidateLicense")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte) (string, error)); ok {
		return rf(signed)
	}
	if rf, ok := ret.Get(0).(func([]byte) string); ok {
		r0 = rf(signed)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(signed)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewLicenseValidatorIface creates a new instance of LicenseValidatorIface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewLicenseValidatorIface(t interface {
	mock.TestingT
	Cleanup(func())
}) *LicenseValidatorIface {
	mock := &LicenseValidatorIface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
