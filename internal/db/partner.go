package db

import (
	"customer-partner/internal/entities"
)

func NewPartnerInMemoryRepository() *PartnerInMemoryRepository {
	return &PartnerInMemoryRepository{partners: demoData}
}

// PartnerInMemoryRepository saves partners in memory and initialises them with some demo data.
type PartnerInMemoryRepository struct {
	partners []entities.Partner
}

// GetPartnersByMaterial returns partners filtered by material.
func (r *PartnerInMemoryRepository) GetPartnersByMaterial(material string) []entities.Partner {
	var filtered []entities.Partner
	for _, partner := range r.partners {
		// TODO consider making partner.ExperiencedMaterial a set which would remove the for loop for this check.
		//    However, this will introduce a data mapping layer to still render it as a list in the api.
		for _, expMaterial := range partner.ExperiencedMaterial {
			if expMaterial == material {
				filtered = append(filtered, partner)
				continue
			}
		}
	}
	return filtered
}

// GetPartnerByID returns a partner by an id.
// Can return entities.ErrRecordNotExist when partner with given id does not exist.
func (r *PartnerInMemoryRepository) GetPartnerByID(id string) (entities.Partner, error) {
	for _, partner := range r.partners {
		if partner.ID == id {
			return partner, nil
		}
	}
	return entities.Partner{}, entities.ErrRecordNotExist
}
