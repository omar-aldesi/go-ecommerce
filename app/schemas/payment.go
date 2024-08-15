package schemas

type NewPaymentSchema struct {
	Amount              float64 `json:"amount" binding:"required"`
	Currency            string  `json:"currency" binding:"required"`
	Status              string  `json:"status" binding:"required"`
	Gateway             string  `json:"gateway" binding:"required"`
	PaymentIntentID     string  `json:"payment_intent_id" binding:"required"`
	PaymentClientSecret string  `json:"payment_client_secret" binding:"required"`
	ReceiptEmail        string  `json:"receipt_email" binding:"required"`
}
