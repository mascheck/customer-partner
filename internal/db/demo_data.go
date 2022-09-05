package db

import "customer-partner/internal/entities"

var demoData = []entities.Partner{
	{
		ID:                  "1",
		Name:                "John and Johnson",
		ExperiencedMaterial: []string{"wood", "carpet", "tiles"},
		Address: entities.Address{
			Latitude:  48.1360,
			Longitude: 11.6875,
		},
		OperatingRadius: 100,
		Rating:          3,
	},
	{
		ID:                  "2",
		Name:                "Peter Skywalker",
		ExperiencedMaterial: []string{"carpet", "tiles"},
		Address: entities.Address{
			Latitude:  48.1360,
			Longitude: 11.6875,
		},
		OperatingRadius: 50,
		Rating:          4,
	}, {
		ID:                  "3",
		Name:                "Wood-hugger Gmbh",
		ExperiencedMaterial: []string{"wood"},
		Address: entities.Address{
			Latitude:  48.1360,
			Longitude: 11.6875,
		},
		OperatingRadius: 50,
		Rating:          5,
	},
}
