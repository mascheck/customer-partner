package entities

import "errors"

var ErrRecordNotExist = errors.New("record not exist")

type Address struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Partner struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	ExperiencedMaterial []string `json:"experienced_material"`
	Address             Address  `json:"address"`
	OperatingRadius     int      `json:"operating_radius"`
	Rating              int      `json:"rating"`
}
