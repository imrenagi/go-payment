package card

// Data is credit card data and used in CardData
type Data struct {
	AccountNumber string `json:"account_number"`
	ExpMonth      string `json:"exp_month"`
	ExpYear       string `json:"exp_year"`
	CVN           string `json:"cvn,omitempty"`
}

// CreateChargeParams contains parameters for CreateCharge
type CreateChargeParams struct {
	TokenID          string  `json:"token_id" validate:"required"`
	ExternalID       string  `json:"external_id" validate:"required"`
	Amount           float64 `json:"amount" validate:"required"`
	AuthenticationID string  `json:"authentication_id,omitempty"`
	CardCVN          string  `json:"card_cvn,omitempty"`
	Capture          *bool   `json:"capture,omitempty"`
	CardData         *Data   `json:"card_data,omitempty"`
	Descriptor       string  `json:"descriptor,omitempty"`
	MidLabel         string  `json:"mid_label,omitempty"`
	Currency         string  `json:"currency,omitempty"`
	IsRecurring      *bool   `json:"is_recurring,omitempty"`
}

// CaptureChargeParams contains parameters for CaptureCharge
type CaptureChargeParams struct {
	ChargeID string  `json:"-" validate:"required"`
	Amount   float64 `json:"amount" validate:"required"`
}

// GetChargeParams contains parameters for GetCharge
type GetChargeParams struct {
	ChargeID string `json:"credit_card_charge_id" validate:"required"`
}

// CreateRefundParams contains parameters for CreateRefund
type CreateRefundParams struct {
	IdempotencyKey string  `json:"-"`
	ChargeID       string  `json:"-" validate:"required"`
	Amount         float64 `json:"amount" validate:"required"`
	ExternalID     string  `json:"external_id" validate:"required"`
}

// ReverseAuthorizationParams contains parameters for ReverseAuthorization
type ReverseAuthorizationParams struct {
	ChargeID   string `json:"-" validate:"required"`
	ExternalID string `json:"external_id" validate:"required"`
}
