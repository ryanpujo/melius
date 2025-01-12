package models

type Country struct {
	ID   uint   `json:"id"`
	Name string `json:"name" binding:"required"`
}

type State struct {
	ID      uint    `json:"id"`
	Name    string  `json:"name" binding:"required"`
	Country Country `json:"country" binding:"required"`
}

type City struct {
	ID    uint   `json:"id"`
	Name  string `json:"name" binding:"required"`
	State State  `json:"state" binding:"required"`
}

type Address struct {
	ID          uint   `json:"id"`
	AddressLine string `json:"address_line" binding:"required"`
	PostalCode  string `json:"postal_code" binding:"required"`
	IsMain      bool   `json:"is_main"`
	City        City   `json:"city" binding:"required"`
}
