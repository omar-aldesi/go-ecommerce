package models

import "gorm.io/gorm"

type ShippingAddress struct {
	gorm.Model
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2"`
	City         string `json:"city"`
	Country      string `json:"country"`
	Postcode     string `json:"postcode"`
	State        string `json:"state"`

	OrderID uint `json:"order_id"`
}
