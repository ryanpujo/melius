package models

type Country struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type State struct {
	ID      uint    `json:"id"`
	Name    string  `json:"name"`
	Country Country `json:"country"`
}

type City struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	State State  `json:"state"`
}

type Address struct {
	ID          uint   `json:"id"`
	AddressLine string `json:"address_line"`
	PostalCode  string `json:"postal_code"`
	IsMain      bool   `json:"is_main"`
	City        City   `json:"city"`
}
