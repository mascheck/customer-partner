package db

import (
	"customer-partner/internal/entities"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPartnerInMemoryRepository_GetPartnersByMaterial(t *testing.T) {
	type testCase struct {
		name   string
		data   []entities.Partner
		expIDs []string
		expLen int
	}
	tests := []testCase{
		{
			name:   "Returns empty list when no partners found (no data)",
			data:   []entities.Partner{},
			expIDs: []string{},
			expLen: 0,
		},
		{
			name:   "Returns empty list when no partners found",
			data:   []entities.Partner{{ID: "123", ExperiencedMaterial: []string{"tiles"}}},
			expIDs: []string{},
			expLen: 0,
		},
		{
			name: "Returns empty list when no partners found",
			data: []entities.Partner{
				{ID: "123", ExperiencedMaterial: []string{"tiles"}},
				{ID: "234", ExperiencedMaterial: []string{"wood"}},
				{ID: "345", ExperiencedMaterial: []string{"carpet", "wood"}},
				{ID: "456", ExperiencedMaterial: []string{"carpet", "tiles"}},
				{ID: "567", ExperiencedMaterial: []string{"carpet", "tiles", "wood"}},
			},
			expIDs: []string{"234", "345", "567"},
			expLen: 3,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repo := NewPartnerInMemoryRepository()
			repo.partners = tt.data

			actual := repo.GetPartnersByMaterial("wood")

			assert.Len(t, actual, tt.expLen)
			for _, partner := range actual {
				assert.Contains(t, tt.expIDs, partner.ID)
			}
		})
	}
}

func TestPartnerInMemoryRepository_GetPartnerByID(t *testing.T) {
	type testCase struct {
		name      string
		data      []entities.Partner
		expResult entities.Partner
		expErr    error
	}
	tests := []testCase{
		{
			name:      "Returns valid partner when present",
			data:      []entities.Partner{{ID: "123"}},
			expResult: entities.Partner{ID: "123"},
			expErr:    nil,
		},
		{
			name:      "Returns ErrRecordNotExist when partner not present",
			data:      []entities.Partner{},
			expResult: entities.Partner{},
			expErr:    entities.ErrRecordNotExist,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repo := NewPartnerInMemoryRepository()
			repo.partners = tt.data

			actual, err := repo.GetPartnerByID("123")

			assert.Equal(t, tt.expResult, actual)
			assert.Equal(t, tt.expErr, err)
		})
	}
}
