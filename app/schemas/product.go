package schemas

type AddonSchema struct {
	AddonID  uint `json:"id" binding:"required"`
	Quantity uint `json:"quantity" binding:"required"`
}

type AddonResponse struct {
	AddonID uint    `json:"id"`
	Price   float64 `json:"price"`
	Tax     float64 `json:"tax"`
}

type VariationOptionSchema struct {
	VariationOptionID uint `json:"id" binding:"required"`
}

type ProductVariationSchema struct {
	ProductVariationID uint                    `json:"id" binding:"required"`
	Options            []VariationOptionSchema `json:"options"`
}

type ProductVariationResponse struct {
	ProductVariationID uint `json:"id"`
}

type ProductResponseSchema struct {
	ID            uint                       `json:"id"`
	Price         float64                    `json:"price"`
	Image         string                     `json:"image"`
	Description   string                     `json:"description"`
	Tags          []string                   `json:"tags"`
	Stock         uint                       `json:"stock"`
	DiscountType  string                     `json:"discount_type"`
	DiscountValue float64                    `json:"discount_value"`
	TotalSales    uint                       `json:"total_sales"`
	Variations    []ProductVariationResponse `json:"variations"`
	Addons        []AddonResponse            `json:"addons"`
	CategoryID    uint                       `json:"category_id"`
	BranchID      uint                       `json:"branch_id"`
}
