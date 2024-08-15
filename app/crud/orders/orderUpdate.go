package orders

import (
	"ecommerce/app/core"
	"ecommerce/app/models"
	"errors"
	"gorm.io/gorm"
	"net/http"
)

func UpdateOrderStatus(db *gorm.DB, user models.User, orderId uint, status string) error {
	var dbOrder models.Order
	if err := db.First(&dbOrder, orderId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &core.HTTPError{
				StatusCode: http.StatusNotFound,
				Message:    "Order not found",
			}
		}
		return &core.HTTPError{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}
	if dbOrder.UserID != user.ID {
		return &core.HTTPError{
			StatusCode: http.StatusUnauthorized,
			Message:    "User not authorized to update this order",
		}
	}
	if !dbOrder.ValidateStatus(status) {
		return &core.HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid status",
		}
	}
	dbOrder.Status = status
	if err := db.Save(&dbOrder).Error; err != nil {
		return &core.HTTPError{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}
	return nil
}
