package virtualaccount

import "time"

// CreateFixedVAParams contains parameters for CreateFixedVA
type CreateFixedVAParams struct {
	ForUserID            string     `json:"-"`
	ExternalID           string     `json:"external_id" validate:"required"`
	BankCode             string     `json:"bank_code" validate:"required"`
	Name                 string     `json:"name" validate:"required"`
	VirtualAccountNumber string     `json:"virtual_account_number,omitempty"`
	IsClosed             *bool      `json:"is_closed,omitempty"`
	IsSingleUse          *bool      `json:"is_single_use,omitempty"`
	ExpirationDate       *time.Time `json:"expiration_date,omitempty"`
	SuggestedAmount      float64    `json:"suggested_amount,omitempty"`
	ExpectedAmount       float64    `json:"expected_amount,omitempty"`
	Description          string     `json:"description,omitempty"`
}

// GetFixedVAParams contains parameters for GetFixedVA
type GetFixedVAParams struct {
	ID string `json:"id" validate:"required"`
}

// UpdateFixedVAParams contains parameters for UpdateFixedVA
type UpdateFixedVAParams struct {
	ID              string     `json:"-" validate:"required"`
	IsSingleUse     *bool      `json:"is_single_use,omitempty"`
	ExpirationDate  *time.Time `json:"expiration_date,omitempty"`
	SuggestedAmount float64    `json:"suggested_amount,omitempty"`
	ExpectedAmount  float64    `json:"expected_amount,omitempty"`
	Description     string     `json:"description,omitempty"`
}

// GetPaymentParams contains parameters for GetPayment
type GetPaymentParams struct {
	PaymentID string `json:"payment_id" validate:"required"`
}
