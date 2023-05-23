package models

import "gorm.io/gorm"

type TrashBinType int

const (
	Paper TrashBinType = iota
	Glass
	PlasticMetal
	Organic
	Other
)

type Product struct {
	gorm.Model           `json:"-"`
	EAN                  string       `gorm:"unique" json:"ean"`
	Name                 string       `json:"name"`
	CO2EmissionPrevented int          `json:"emission_prevented"`
	Mass                 int          `json:"mass"`
	TypeOfTrash          TrashBinType `json:"type_of_trash"`
}
