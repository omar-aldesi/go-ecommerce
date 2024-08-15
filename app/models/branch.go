package models

import "gorm.io/gorm"

type Branch struct {
	gorm.Model
	Name     string    `gorm:"type:varchar(20);not null"`
	Products []Product `gorm:"foreignkey:BranchID"`
}
