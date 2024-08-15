package models

import (
	"ecommerce/app/schemas"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"time"
)

type Product struct {
	gorm.Model
	Price                float64            `json:"price"`
	Image                string             `json:"image"`
	Description          string             `json:"description"`
	Tags                 pq.StringArray     `gorm:"type:text[]" json:"tags"`
	IsActive             bool               `json:"is_active"`
	StockType            string             `json:"stock_type"`
	DailyStock           uint               `json:"daily_stock"`
	Stock                uint               `json:"stock"`
	LastDailyStockUpdate time.Time          `json:"last_daily_stock_update"`
	DiscountType         string             `json:"discount_type"`
	DiscountValue        float64            `json:"discount_value"`
	TotalSales           uint               `json:"total_sales" gorm:"default:0"`
	Variations           []ProductVariation `json:"variations" gorm:"foreignKey:ProductID"`
	Addons               []Addon            `json:"addons" gorm:"many2many:product_addons;"`
	CategoryID           uint               `json:"category_id"`
	Category             Category           `json:"category" gorm:"foreignKey:CategoryID"`
	BranchID             uint               `json:"branch_id"`
	Branch               Branch             `json:"branch" gorm:"foreignKey:BranchID"`
}

type ProductVariation struct {
	gorm.Model
	Title         string            `json:"title"`
	Type          string            `json:"type"`
	MinSelections uint              `json:"min_selections" gorm:"default:0"`
	MaxSelections uint              `json:"max_selections" gorm:"default:0"`
	Required      bool              `json:"required" gorm:"default:false"`
	ProductID     uint              `json:"product_id" gorm:"foreignKey:ProductID"`
	Product       Product           `json:"product"`
	Options       []VariationOption `gorm:"many2many:product_variation_options;"`
}

type VariationOption struct {
	gorm.Model
	Title      string             `json:"title"`
	Price      float64            `json:"price"`
	Variations []ProductVariation `gorm:"many2many:product_variation_options;"`
}

type Addon struct {
	gorm.Model
	Title    string  `json:"title"`
	Price    float64 `json:"price"`
	Tax      float64
	Products []Product `json:"products" gorm:"many2many:product_addons;"`
}

type Review struct {
	gorm.Model
	Rating  uint   `json:"rating"`
	Comment string `json:"comment"`

	ProductID uint    `json:"product_id"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`

	UserID uint `json:"user_id"`
	User   User `json:"user" gorm:"foreignKey:UserID"`
}

func (p *Product) ToResponse() schemas.ProductResponseSchema {
	variationSchemas := make([]schemas.ProductVariationResponse, len(p.Variations))
	for i, v := range p.Variations {
		variationSchemas[i] = v.ToResponse()
	}

	addonSchemas := make([]schemas.AddonResponse, len(p.Addons))
	for i, a := range p.Addons {
		addonSchemas[i] = a.ToResponse()
	}
	return schemas.ProductResponseSchema{
		ID:            p.ID,
		Price:         p.Price,
		Image:         p.Image,
		Description:   p.Description,
		Tags:          p.Tags,
		Stock:         p.Stock,
		DiscountType:  p.DiscountType,
		DiscountValue: p.DiscountValue,
		TotalSales:    p.TotalSales,
		Variations:    variationSchemas,
		Addons:        addonSchemas,
		CategoryID:    p.CategoryID,
		BranchID:      p.BranchID,
	}
}
func (v *ProductVariation) ToResponse() schemas.ProductVariationResponse {
	return schemas.ProductVariationResponse{
		ProductVariationID: v.ProductID,
	}
}

func (a *Addon) ToResponse() schemas.AddonResponse {
	return schemas.AddonResponse{
		AddonID: a.ID,
		Price:   a.Price,
		Tax:     a.Tax,
	}
}
