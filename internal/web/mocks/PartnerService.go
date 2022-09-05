// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	domain "customer-partner/internal/domain"
	entities "customer-partner/internal/entities"

	mock "github.com/stretchr/testify/mock"
)

// PartnerService is an autogenerated mock type for the PartnerService type
type PartnerService struct {
	mock.Mock
}

// GetPartner provides a mock function with given fields: id
func (_m *PartnerService) GetPartner(id string) (entities.Partner, error) {
	ret := _m.Called(id)

	var r0 entities.Partner
	if rf, ok := ret.Get(0).(func(string) entities.Partner); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(entities.Partner)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPartners provides a mock function with given fields: opts
func (_m *PartnerService) GetPartners(opts domain.GetPartnersOpts) []entities.Partner {
	ret := _m.Called(opts)

	var r0 []entities.Partner
	if rf, ok := ret.Get(0).(func(domain.GetPartnersOpts) []entities.Partner); ok {
		r0 = rf(opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entities.Partner)
		}
	}

	return r0
}

type mockConstructorTestingTNewPartnerService interface {
	mock.TestingT
	Cleanup(func())
}

// NewPartnerService creates a new instance of PartnerService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPartnerService(t mockConstructorTestingTNewPartnerService) *PartnerService {
	mock := &PartnerService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}