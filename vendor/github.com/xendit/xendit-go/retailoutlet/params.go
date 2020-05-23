package retailoutlet

import (
	"time"

	"github.com/xendit/xendit-go"
)

// CreateFixedPaymentCodeParams contains parameters for CreateFixedPaymentCode
type CreateFixedPaymentCodeParams struct {
	ExternalID       string                      `json:"external_id" validate:"required"`
	RetailOutletName xendit.RetailOutletNameEnum `json:"retail_outlet_name" validate:"required"`
	Name             string                      `json:"name" validate:"required"`
	ExpectedAmount   float64                     `json:"expected_amount" validate:"required"`
	PaymentCode      string                      `json:"payment_code,omitempty"`
	ExpirationDate   *time.Time                  `json:"expiration_date,omitempty"`
	IsSingleUse      *bool                       `json:"is_single_use,omitempty"`
}

// GetFixedPaymentCodeParams contains parameters for GetFixedPaymentCode
type GetFixedPaymentCodeParams struct {
	FixedPaymentCodeID string `json:"fixed_payment_code_id" validate:"required"`
}

// UpdateFixedPaymentCodeParams contains parameters for UpdateFixedPaymentCode
type UpdateFixedPaymentCodeParams struct {
	FixedPaymentCodeID string     `json:"-" validate:"required"`
	Name               string     `json:"name,omitempty"`
	ExpectedAmount     float64    `json:"expected_amount,omitempty"`
	ExpirationDate     *time.Time `json:"expiration_date,omitempty"`
}
