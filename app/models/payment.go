package models

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	Amount              float64 `json:"amount"`
	Currency            string  `json:"currency"`
	Status              string  `json:"status"`
	Gateway             string  `json:"gateway"`
	PaymentIntentID     string  `json:"payment_intent_id"`
	PaymentClientSecret string  `json:"payment_client_secret"`
	ReceiptEmail        string  `json:"recipient_email"`
	UserID              uint    `json:"user_id"`
	User                User    `json:"user" gorm:"foreignKey:UserID"`
	OrderID             uint    `json:"order_id" gorm:"foreignKey:PaymentID"`
}
