package v1

import (
	"ecommerce/app/core"
	"ecommerce/app/crud"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckCoupon(c *gin.Context) {
	type request struct {
		CouponCode string `json:"coupon_code" binding:"required"`
	}
	var requestData request
	if err := c.ShouldBindJSON(&requestData); err != nil {
		core.HandleValidationErrors(c, err)
		return
	}
	db := core.GetDB()
	coupon, err := crud.CouponIsValid(db, requestData.CouponCode)
	if err != nil {
		core.CustomErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "coupon is valid to use",
		"coupon":  coupon,
	})
}

func CouponsRouter(router *gin.Engine) {
	public := router.Group("/api/v1/coupons")
	{
		public.GET("/check", CheckCoupon)
	}
}
