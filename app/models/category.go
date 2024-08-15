package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Title         string        `json:"title"`
	SubCategories []SubCategory `json:"sub_categories" gorm:"foreignKey:CategoryID"`
}

type SubCategory struct {
	gorm.Model
	Title      string   `json:"title"`
	CategoryID uint     `json:"category_id"`
	Category   Category `json:"-" gorm:"foreignKey:CategoryID"`
}
