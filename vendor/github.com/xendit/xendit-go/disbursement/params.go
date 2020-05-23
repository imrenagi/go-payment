package disbursement

import (
	"net/url"
)

// CreateParams contains parameters for Create
type CreateParams struct {
	IdempotencyKey    string   `json:"-"`
	ForUserID         string   `json:"-"`
	ExternalID        string   `json:"external_id" validate:"required"`
	BankCode          string   `json:"bank_code" validate:"required"`
	AccountHolderName string   `json:"account_holder_name" validate:"required"`
	AccountNumber     string   `json:"account_number" validate:"required"`
	Description       string   `json:"description" validate:"required"`
	Amount            float64  `json:"amount" validate:"required"`
	EmailTo           []string `json:"email_to,omitempty"`
	EmailCC           []string `json:"email_cc,omitempty"`
	EmailBCC          []string `json:"email_bcc,omitempty"`
}

// GetByIDParams contains parameters for GetByID
type GetByIDParams struct {
	DisbursementID string `json:"disbursement_id" validate:"required"`
	ForUserID      string `json:"-"`
}

// GetByExternalIDParams contains parameters for GetByExternalID
type GetByExternalIDParams struct {
	ExternalID string `json:"external_id" validate:"required"`
	ForUserID  string `json:"-"`
}

// QueryString creates query string from GetByExternalIDParams, ignores nil values
func (p *GetByExternalIDParams) QueryString() string {
	urlValues := &url.Values{}

	urlValues.Add("external_id", p.ExternalID)

	return urlValues.Encode()
}

// CreateBatchParams contains parameters for CreateBatch
type CreateBatchParams struct {
	IdempotencyKey string             `json:"-"`
	ForUserID      string             `json:"-"`
	Reference      string             `json:"reference" validate:"required"`
	Disbursements  []DisbursementItem `json:"disbursements" validate:"required"`
}

// DisbursementItem is data that contained in CreateBatch at Disbursements
type DisbursementItem struct {
	Amount            float64  `json:"amount" validate:"required"`
	BankCode          string   `json:"bank_code" validate:"required"`
	BankAccountName   string   `json:"bank_account_name" validate:"required"`
	BankAccountNumber string   `json:"bank_account_number" validate:"required"`
	Description       string   `json:"description" validate:"required"`
	ExternalID        string   `json:"external_id,omitempty"`
	EmailTo           []string `json:"email_to,omitempty"`
	EmailCC           []string `json:"email_cc,omitempty"`
	EmailBCC          []string `json:"email_bcc,omitempty"`
}
