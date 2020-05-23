package invoice

import (
	"net/url"
	"strconv"
	"time"

	"github.com/xendit/xendit-go/utils/urlvalues"
)

// CreateParams contains parameters for Create
type CreateParams struct {
	ForUserID                string   `json:"-"`
	ExternalID               string   `json:"external_id" validate:"required"`
	PayerEmail               string   `json:"payer_email" validate:"required"`
	Description              string   `json:"description" validate:"required"`
	Amount                   float64  `json:"amount" validate:"required"`
	ShouldSendEmail          *bool    `json:"should_send_email,omitempty"`
	CallbackVirtualAccountID string   `json:"callback_virtual_account_id,omitempty"`
	InvoiceDuration          int      `json:"invoice_duration,omitempty"`
	SuccessRedirectURL       string   `json:"success_redirect_url,omitempty"`
	FailureRedirectURL       string   `json:"failure_redirect_url,omitempty"`
	PaymentMethods           []string `json:"payment_methods,omitempty"`
	MidLabel                 string   `json:"mid_label,omitempty"`
	Currency                 string   `json:"currency,omitempty"`
	FixedVA                  *bool    `json:"fixed_va,omitempty"`
}

// GetParams contains parameters for Get
type GetParams struct {
	ID string `json:"id" validate:"required"`
}

// GetAllParams contains parameters for GetAll
type GetAllParams struct {
	Statuses           []string  `json:"statuses,omitempty"`
	Limit              int       `json:"limit,omitempty"`
	CreatedAfter       time.Time `json:"created_after,omitempty"`
	CreatedBefore      time.Time `json:"created_before,omitempty"`
	PaidAfter          time.Time `json:"paid_after,omitempty"`
	PaidBefore         time.Time `json:"paid_before,omitempty"`
	ExpiredAfter       time.Time `json:"expired_after,omitempty"`
	ExpiredBefore      time.Time `json:"expired_before,omitempty"`
	LastInvoiceID      string    `json:"last_invoice_id,omitempty"`
	ClientTypes        []string  `json:"client_types,omitempty"`
	PaymentChannels    []string  `json:"payment_channels,omitempty"`
	OnDemandLink       string    `json:"on_demand_link,omitempty"`
	RecurringPaymentID string    `json:"recurring_payment_id,omitempty"`
}

// QueryString creates query string from GetAllParams, ignores nil values
func (p *GetAllParams) QueryString() string {
	urlValues := &url.Values{}

	urlvalues.AddStringSliceToURLValues(urlValues, p.Statuses, "statuses")
	if p.Limit > 0 {
		urlValues.Add("limit", strconv.Itoa(p.Limit))
	}
	urlvalues.AddTimeToURLValues(urlValues, p.CreatedAfter, "created_after")
	urlvalues.AddTimeToURLValues(urlValues, p.CreatedBefore, "created_before")
	urlvalues.AddTimeToURLValues(urlValues, p.PaidAfter, "paid_after")
	urlvalues.AddTimeToURLValues(urlValues, p.PaidBefore, "paid_before")
	urlvalues.AddTimeToURLValues(urlValues, p.PaidBefore, "paid_before")
	urlvalues.AddTimeToURLValues(urlValues, p.ExpiredAfter, "expired_after")
	urlvalues.AddTimeToURLValues(urlValues, p.ExpiredBefore, "expired_before")
	urlvalues.AddStringSliceToURLValues(urlValues, p.ClientTypes, "client_types")
	urlvalues.AddStringSliceToURLValues(urlValues, p.PaymentChannels, "payment_channels")
	if p.OnDemandLink != "" {
		urlValues.Add("on_demand", p.OnDemandLink)
	}
	if p.RecurringPaymentID != "" {
		urlValues.Add("recurring_payment_id", p.RecurringPaymentID)
	}

	return urlValues.Encode()
}

// ExpireParams contains parameters for Expire
type ExpireParams struct {
	ID string `json:"id" validate:"required"`
}
