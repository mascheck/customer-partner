package domain

import (
	"customer-partner/internal/entities"
	"math"
	"sort"
)

type match struct {
	partner  entities.Partner
	distance float64
}

// byRatingAndDistance implements the sort.Interface for a custom sorting.
type byRatingAndDistance []match

func (rd byRatingAndDistance) Len() int {
	return len(rd)
}
func (rd byRatingAndDistance) Less(i, j int) bool {
	if rd[i].partner.Rating > rd[j].partner.Rating {
		return true
	}
	return rd[i].distance < rd[j].distance
}
func (rd byRatingAndDistance) Swap(i, j int) {
	rd[i], rd[j] = rd[j], rd[i]
}

// GetPartnersOpts combines attributes necessary for finding the best match.
type GetPartnersOpts struct {
	Material            string
	CustomerAddressLong float64
	CustomerAddressLat  float64
}

// PartnerRepository defines an interface which a persistence storage must provide.
type PartnerRepository interface {
	GetPartnersByMaterial(material string) []entities.Partner
	GetPartnerByID(id string) (entities.Partner, error)
}

func NewPartnerService(repository PartnerRepository) *PartnerService {
	return &PartnerService{repository: repository}
}

// PartnerService implements the domain logic of the partner domain.
type PartnerService struct {
	repository PartnerRepository
}

// GetPartners retrieves the partners from the persistence storage and sorts them after best match. Partners not in
// operating radius are sorted out.
func (s *PartnerService) GetPartners(opts GetPartnersOpts) []entities.Partner {
	partners := s.repository.GetPartnersByMaterial(opts.Material)
	matches := convertPartnersToMatchesAndFilterByOperatingRadius(
		partners,
		opts.CustomerAddressLat,
		opts.CustomerAddressLong,
	)
	sort.Sort(byRatingAndDistance(matches))
	return convertMatchesToPartners(matches)
}

// GetPartner finds a partner by its id.
// Can return entities.ErrRecordNotExist when partner with given id does not exist.
func (s *PartnerService) GetPartner(id string) (entities.Partner, error) {
	return s.repository.GetPartnerByID(id)
}

func convertPartnersToMatchesAndFilterByOperatingRadius(
	partners []entities.Partner,
	customerLat float64,
	customerLong float64,
) []match {
	var matches []match
	for _, partner := range partners {
		d := distance(
			customerLat,
			customerLong,
			partner.Address.Latitude,
			partner.Address.Longitude,
			"K",
		)
		if d < float64(partner.OperatingRadius) {
			match := match{
				partner:  partner,
				distance: d,
			}
			matches = append(matches, match)
		}
	}
	return matches
}

func convertMatchesToPartners(matches []match) []entities.Partner {
	var partners []entities.Partner
	for _, m := range matches {
		partners = append(partners, m.partner)
	}
	return partners
}

// :::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
// :::                                                                         :::
// :::  This routine calculates the distance between two points (given the     :::
// :::  latitude/longitude of those points). It is being used to calculate     :::
// :::  the distance between two locations using GeoDataSource (TM) products  :::
// :::                                                                         :::
// :::  Definitions:                                                           :::
// :::    South latitudes are negative, east longitudes are positive           :::
// :::                                                                         :::
// :::  Passed to function:                                                    :::
// :::    lat1, lon1 = Latitude and Longitude of point 1 (in decimal degrees)  :::
// :::    lat2, lon2 = Latitude and Longitude of point 2 (in decimal degrees)  :::
// :::    unit = the unit you desire for results                               :::
// :::           where: 'M' is statute miles (default)                         :::
// :::                  'K' is kilometers                                      :::
// :::                  'N' is nautical miles                                  :::
// :::                                                                         :::
// :::  Worldwide cities and other features databases with latitude longitude  :::
// :::  are available at https://www.geodatasource.com                         :::
// :::                                                                         :::
// :::  For enquiries, please contact sales@geodatasource.com                  :::
// :::                                                                         :::
// :::  Official Web site: https://www.geodatasource.com                       :::
// :::                                                                         :::
// :::               GeoDataSource.com (C) All Rights Reserved 2022            :::
// :::                                                                         :::
// :::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
func distance(lat1 float64, lng1 float64, lat2 float64, lng2 float64, unit ...string) float64 {
	const PI float64 = 3.141592653589793

	radlat1 := PI * lat1 / 180
	radlat2 := PI * lat2 / 180

	theta := lng1 - lng2
	radtheta := PI * theta / 180

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515

	if len(unit) > 0 {
		if unit[0] == "K" {
			dist = dist * 1.609344
		} else if unit[0] == "N" {
			dist = dist * 0.8684
		}
	}
	return dist
}
