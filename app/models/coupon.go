package models

import (
	"gorm.io/gorm"
	"time"
)

type Coupon struct {
	gorm.Model
	Code         string    `gorm:"unique;not null"`
	Discount     float64   `gorm:"not null"`
	DiscountType string    `gorm:"not null"`
	ExpireDate   time.Time `gorm:"not null"`
	IsActive     bool      `gorm:"default:true"`
	UsageCount   int       `gorm:"default:0"`
	MaxUsage     int       `gorm:"default:1"`
}

func (c *Coupon) ApplyDiscount(total float64) float64 {
	if c.DiscountType == "percentage" {
		return total * (1 - c.Discount/100)
	} else if c.DiscountType == "fixed" {
		return total - c.Discount
	}
	return total
}
