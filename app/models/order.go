package models

import (
	"ecommerce/app/schemas"
	"gorm.io/gorm"
	"time"
)

type Order struct {
	gorm.Model

	Status       string    `gorm:"type:varchar(20);not null" json:"status"`
	Type         string    `gorm:"type:varchar(20);not null" json:"type"`
	Total        float64   `gorm:"type:decimal(10,2);not null" json:"total"`
	SubTotal     float64   `gorm:"type:decimal(10,2);not null" json:"sub_total"`
	IsPaid       bool      `gorm:"type:boolean;not null;default:false" json:"is_paid"`
	IsScheduled  bool      `gorm:"type:boolean;not null;default:false" json:"is_scheduled"`
	ScheduleTime time.Time `json:"schedule_time"`
	Coupon       string    `gorm:"type:varchar(20);null" json:"coupon"`
	Discount     float64   `gorm:"type:decimal(10, 2);null" json:"discount"`

	UserID uint
	User   User `gorm:"foreignkey:UserID"`

	BranchID uint `gorm:"foreignkey:OrderID"`

	Products []OrderItem `gorm:"foreignkey:OrderID"`

	ShippingAddress ShippingAddress `gorm:"foreignKey:OrderID;references:ID"`

	Payment Payment `gorm:"foreignkey:OrderID"`
}

type OrderItem struct {
	gorm.Model
	Quantity   uint    `gorm:"not null" json:"quantity"`
	TotalPrice float64 `gorm:"not null" json:"total_price"`

	ProductID uint    `gorm:"not null" json:"product_id"`
	Product   Product `gorm:"foreignkey:ProductID" json:"product"`

	OrderID uint  `gorm:"not null" json:"order_id"`
	Order   Order `gorm:"foreignkey:OrderID" json:"order"`

	SelectedVariations []OrderItemVariation `gorm:"foreignkey:OrderItemID" json:"selected_variations"`
	SelectedAddons     []OrderItemAddon     `gorm:"foreignkey:OrderItemID" json:"selected_addons"`
}

type OrderItemVariation struct {
	gorm.Model
	OrderItemID        uint              `gorm:"not null" json:"order_item_id"`
	ProductVariationID uint              `gorm:"not null" json:"product_variation_id"`
	ProductVariation   ProductVariation  `gorm:"foreignkey:ProductVariationID" json:"product_variation"`
	SelectedOptions    []VariationOption `gorm:"many2many:order_item_variation_options;" json:"selected_options"`
}

type OrderItemAddon struct {
	gorm.Model
	OrderItemID uint  `gorm:"not null" json:"order_item_id"`
	AddonID     uint  `gorm:"not null" json:"addon_id"`
	Addon       Addon `gorm:"foreignkey:AddonID" json:"addon"`
	Quantity    uint  `gorm:"not null" json:"quantity"`
}

func (o *Order) ToResponse() schemas.OrderResponseSchema {
	var productSchemas []schemas.OrderItemSchema
	for _, item := range o.Products {
		productSchema := schemas.OrderItemSchema{
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			Variations: convertVariations(item.SelectedVariations),
			Addons:     convertAddons(item.SelectedAddons),
		}
		productSchemas = append(productSchemas, productSchema)
	}

	return schemas.OrderResponseSchema{
		ID:           o.ID,
		Status:       o.Status,
		Type:         o.Type,
		Total:        o.Total,
		SubTotal:     o.SubTotal,
		IsPaid:       o.IsPaid,
		IsScheduled:  o.IsScheduled,
		ScheduleTime: o.ScheduleTime,
		UserID:       o.UserID,
		BranchID:     o.BranchID,
		Products:     productSchemas,
	}
}

// Helper Functions
func convertVariations(variations []OrderItemVariation) []schemas.ProductVariationSchema {
	var variationSchemas []schemas.ProductVariationSchema
	for _, variation := range variations {
		variationSchema := schemas.ProductVariationSchema{
			ProductVariationID: variation.ProductVariationID,
			Options:            convertVariationOptions(variation.SelectedOptions),
		}
		variationSchemas = append(variationSchemas, variationSchema)
	}
	return variationSchemas
}

func convertVariationOptions(options []VariationOption) []schemas.VariationOptionSchema {
	var optionSchemas []schemas.VariationOptionSchema
	for _, option := range options {
		optionSchema := schemas.VariationOptionSchema{
			VariationOptionID: option.ID,
		}
		optionSchemas = append(optionSchemas, optionSchema)
	}
	return optionSchemas
}

func convertAddons(addons []OrderItemAddon) []schemas.AddonSchema {
	var addonSchemas []schemas.AddonSchema
	for _, addon := range addons {
		addonSchema := schemas.AddonSchema{
			AddonID:  addon.AddonID,
			Quantity: addon.Quantity,
		}
		addonSchemas = append(addonSchemas, addonSchema)
	}
	return addonSchemas
}

func (o *Order) ValidateStatus(status string) bool {
	switch status {
	case "pending", "paid", "shipped":
		return true
	default:
		return false
	}
}
