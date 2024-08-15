package schemas

import "time"

type OrderItemSchema struct {
	ProductID  uint                     `json:"product_id" binding:"required"`
	Quantity   uint                     `json:"quantity" binding:"required"`
	Addons     []AddonSchema            `json:"addons"`
	Variations []ProductVariationSchema `json:"variation"`
	CouponCode string                   `json:"coupon_code"`
}

type ShippingAddressSchema struct {
	AddressLine1 string `json:"address_line_1" binding:"required"`
	AddressLine2 string `json:"address_line_2"  binding:"required"`
	City         string `json:"city"  binding:"required"`
	Country      string `json:"country"  binding:"required"`
	Postcode     string `json:"postcode" binding:"required"`
	State        string `json:"state" binding:"required"`
}

type OrderCreationSchema struct {
	Products        []OrderItemSchema     `json:"products" binding:"required"`
	OrderType       string                `json:"order_type" binding:"required"`
	BranchID        uint                  `json:"branch_id" binding:"required"`
	IsScheduled     bool                  `json:"is_scheduled"`
	ScheduleAt      time.Time             `json:"schedule_time"`
	ShippingAddress ShippingAddressSchema `json:"shipping_address"`
	Payment         NewPaymentSchema      `json:"payment" binding:"required"`
}

type OrderResponseSchema struct {
	ID           uint              `json:"id"`
	Status       string            `json:"status"`
	Type         string            `json:"type"`
	Total        float64           `json:"total"`
	SubTotal     float64           `json:"sub_total"`
	IsPaid       bool              `json:"is_paid"`
	IsScheduled  bool              `json:"is_scheduled"`
	ScheduleTime time.Time         `json:"schedule_time"`
	UserID       uint              `json:"user_id"`
	BranchID     uint              `json:"branch_id"`
	Products     []OrderItemSchema `json:"products"`
}

type UpdateOrderStatusSchema struct {
	Status  string `json:"status" binding:"required"`
	OrderID uint   `json:"order_id" binding:"required"`
}
