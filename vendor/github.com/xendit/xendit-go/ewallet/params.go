package ewallet

import (
	"net/url"
	"time"

	"github.com/xendit/xendit-go"
)

// Item is data that contained in CreatePaymentParams at Items
type Item struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

// CreatePaymentParams contains parameters for CreatePayment
type CreatePaymentParams struct {
	XApiVersion    string                 `json:"-"`
	ForUserID      string                 `json:"-"`
	EWalletType    xendit.EWalletTypeEnum `json:"ewallet_type" validate:"required"`
	ExternalID     string                 `json:"external_id" validate:"required"`
	Amount         float64                `json:"amount" validate:"required"`
	Phone          string                 `json:"phone,omitempty"`
	ExpirationDate *time.Time             `json:"expiration_date,omitempty"`
	CallbackURL    string                 `json:"callback_url,omitempty"`
	RedirectURL    string                 `json:"redirect_url,omitempty"`
	Items          []Item                 `json:"items,omitempty"`
}

// GetPaymentStatusParams contains parameters for GetPaymentStatus
type GetPaymentStatusParams struct {
	ExternalID  string                 `json:"external_id" validate:"required"`
	EWalletType xendit.EWalletTypeEnum `json:"ewallet_type" validate:"required"`
}

// QueryString creates query string from GetPaymentStatusParams, ignores nil values
func (p *GetPaymentStatusParams) QueryString() string {
	urlValues := &url.Values{}

	urlValues.Add("external_id", p.ExternalID)
	urlValues.Add("ewallet_type", string(p.EWalletType))

	return urlValues.Encode()
}
