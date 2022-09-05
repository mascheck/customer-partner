package domain_test

import (
	"customer-partner/internal/domain"
	"customer-partner/internal/domain/mocks"
	"customer-partner/internal/entities"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:generate mockery --name PartnerRepository

func TestPartnerService_GetPartners(t *testing.T) {
	type testCase struct {
		name       string
		opts       domain.GetPartnersOpts
		repoReturn []entities.Partner
		expLen     int
		expIDs     []string
	}
	tests := []testCase{
		{
			name: "Returns empty list on empty return from repo",
			opts: domain.GetPartnersOpts{
				Material:            "wood",
				CustomerAddressLat:  48.3535,
				CustomerAddressLong: 11.7812,
			},
			repoReturn: []entities.Partner{},
			expLen:     0,
		},
		{
			name: "Returns empty list when no partner in range",
			opts: domain.GetPartnersOpts{
				Material:            "wood",
				CustomerAddressLat:  48.3535,
				CustomerAddressLong: 11.7812,
			},
			repoReturn: []entities.Partner{
				{
					ID: "123",
					Address: entities.Address{
						Latitude:  48.4021,
						Longitude: 11.7511,
					},
					OperatingRadius: 5,
					Rating:          5,
				},
			},
			expLen: 0,
		},
		{
			name: "Returns partner when partner is in range",
			opts: domain.GetPartnersOpts{
				Material:            "wood",
				CustomerAddressLat:  48.3535,
				CustomerAddressLong: 11.7812,
			},
			repoReturn: []entities.Partner{
				{
					ID: "123",
					Address: entities.Address{
						Latitude:  48.4021,
						Longitude: 11.7511,
					},
					OperatingRadius: 10,
					Rating:          5,
				},
			},
			expLen: 1,
		},
		{
			name: "Returns partner ordered by rating",
			opts: domain.GetPartnersOpts{
				Material:            "wood",
				CustomerAddressLat:  48.3535,
				CustomerAddressLong: 11.7812,
			},
			repoReturn: []entities.Partner{
				{
					ID: "123",
					Address: entities.Address{
						Latitude:  48.4021,
						Longitude: 11.7511,
					},
					OperatingRadius: 10,
					Rating:          4,
				},
				{
					ID: "234",
					Address: entities.Address{
						Latitude:  48.4021,
						Longitude: 11.7511,
					},
					OperatingRadius: 10,
					Rating:          5,
				},
			},
			expLen: 2,
			expIDs: []string{"234", "123"},
		},
		{
			name: "Returns partner ordered by rating and distance",
			opts: domain.GetPartnersOpts{
				Material:            "wood",
				CustomerAddressLat:  48.3535,
				CustomerAddressLong: 11.7812,
			},
			repoReturn: []entities.Partner{
				{
					ID: "123",
					Address: entities.Address{
						Latitude:  48.2186,
						Longitude: 11.6236,
					},
					OperatingRadius: 20,
					Rating:          4,
				},
				{
					ID: "234",
					Address: entities.Address{
						Latitude:  48.4021,
						Longitude: 11.7511,
					},
					OperatingRadius: 10,
					Rating:          4,
				},
				{
					ID: "345",
					Address: entities.Address{
						Latitude:  48.4021,
						Longitude: 11.7511,
					},
					OperatingRadius: 10,
					Rating:          5,
				},
			},
			expLen: 3,
			expIDs: []string{"345", "234", "123"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.PartnerRepository{}
			repo.On("GetPartnersByMaterial", tt.opts.Material).Return(tt.repoReturn)
			service := domain.NewPartnerService(repo)

			actual := service.GetPartners(tt.opts)

			repo.AssertExpectations(t)
			assert.Len(t, actual, tt.expLen)
			for i, expID := range tt.expIDs {
				assert.Equal(t, expID, actual[i].ID)
			}
		})
	}
}
