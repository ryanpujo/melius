package models

type Country struct {
	ID     uint     `json:"id"`
	Name   string   `json:"name"`
	States []*State `json:"states"`
}

type CountryPayload struct {
	Name  string        `json:"name" binding:"required"`
}

type State struct {
	ID      uint     `json:"id"`
	Name    string   `json:"name"`
	Country *Country `json:"country"`
	Cities  []*City  `json:"cities"`
}

type StatePayload struct {
	Name string       `json:"name" binding:"required"`
}

type City struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name"`
	State     *State     `json:"state"`
	Addresses []*Address `json:"addresses"`
}

type CityPayload struct {
	Name    string          `json:"name" binding:"required"`
}

type Address struct {
	ID          uint   `json:"id"`
	AddressLine string `json:"address_line"`
	PostalCode  string `json:"postal_code"`
	IsMain      bool   `json:"is_main"`
	City        *City  `json:"city"`
}

type AddressPayload struct {
	AddressLine string `json:"address_line" binding:"required"`
	PostalCode  string `json:"postal_code" binding:"required"`
	IsMain      bool   `json:"is_main"`
	CityID      uint   `json:"city_id" binding:"required"`
}
