package crud

import (
	"ecommerce/app/core"
	"ecommerce/app/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func CouponIsValid(db *gorm.DB, couponCode string) (models.Coupon, error) {
	var dbCoupon models.Coupon
	if err := db.
		Where("code = ?", couponCode).
		Where("is_active = ?", true).
		Where("expire_date > ?", time.Now()).
		First(&dbCoupon).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dbCoupon, &core.HTTPError{
				StatusCode: http.StatusNotFound,
				Message:    fmt.Sprintf("Coupon %s not found", couponCode),
			}
		}
		return dbCoupon, &core.HTTPError{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("Error getting coupon: %s", err),
		}
	}
	if dbCoupon.MaxUsage != 0 && dbCoupon.MaxUsage <= dbCoupon.UsageCount {
		return dbCoupon, &core.HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    fmt.Sprintf("Coupon %s has reached maximum usage", couponCode),
		}
	}
	return dbCoupon, nil
}
