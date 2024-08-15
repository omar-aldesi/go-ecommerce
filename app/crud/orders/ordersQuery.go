package orders

import (
	"ecommerce/app/core"
	"ecommerce/app/models"
	"fmt"
	"gorm.io/gorm"
	"net/http"
)

func ListUserOrders(db *gorm.DB, user models.User) ([]models.Order, error) {
	var userOrders []models.Order
	if err := db.Where("user_id = ?", user.ID).
		Preload("Products.SelectedVariations.ProductVariation").
		Preload("Products.SelectedVariations.SelectedOptions").
		Preload("Products.SelectedAddons.Addon").
		Preload("ShippingAddress").
		Find(&userOrders).Error; err != nil {
		return nil, &core.HTTPError{
			Message:    fmt.Sprintf("cannot list user orders: %v", err),
			StatusCode: http.StatusInternalServerError,
		}
	}
	return userOrders, nil
}

func GetOrderByID(db *gorm.DB, user models.User, orderID uint) (models.Order, error) {
	var order models.Order
	if err := db.Preload("Products.SelectedVariations.ProductVariation").
		Preload("Products.SelectedVariations.SelectedOptions").
		Preload("Products.SelectedAddons.Addon").
		Preload("ShippingAddress").
		First(&order, orderID).Error; err != nil {
		return models.Order{}, &core.HTTPError{
			Message:    fmt.Sprintf("cannot get order by id %v", orderID),
			StatusCode: http.StatusNotFound,
		}
	}
	if order.UserID != user.ID {
		return models.Order{}, &core.HTTPError{
			Message:    fmt.Sprintf("user %v has no permission to get order %v", user.ID, orderID),
			StatusCode: http.StatusBadRequest,
		}
	}
	return order, nil
}
