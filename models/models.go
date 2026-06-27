package models

// PaymentRequest is the common payload sent from the frontend tester.
type PaymentRequest struct {
	// Shared
	OrderID     string  `json:"order_id"`
	Amount      float64 `json:"amount"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Email       string  `json:"email"`
	Phone       string  `json:"phone"`
	ProductName string  `json:"product_name"`

	// Credit Card specific
	TokenID        string `json:"token_id"`         // from Midtrans.js
	Bank           string `json:"bank"`             // e.g. "bca", "bni", "mandiri"
	Installment    bool   `json:"installment"`
	InstallmentTerm int   `json:"installment_term"` // 3, 6, 12

	// GoPay specific
	EnableCallback bool   `json:"enable_callback"`
	CallbackURL    string `json:"callback_url"`

	// OVO specific
	OVOPhone string `json:"ovo_phone"`

	// Akulaku / Kredivo
	CallbackFinish string `json:"callback_finish"`
}

// PaymentResponse is what we return to the frontend.
type PaymentResponse struct {
	Success       bool        `json:"success"`
	Message       string      `json:"message"`
	OrderID       string      `json:"order_id,omitempty"`
	PaymentType   string      `json:"payment_type,omitempty"`
	TransactionID string      `json:"transaction_id,omitempty"`
	Data          interface{} `json:"data,omitempty"`
}
