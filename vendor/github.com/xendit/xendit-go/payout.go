package xendit

import "time"

// Payout contains data from Xendit's API response of invoice related request.
// For more details see https://xendit.github.io/apireference/?bash#payouts.
// For documentation of subpackage payout, checkout https://pkg.go.dev/github.com/xendit/xendit-go/payout
type Payout struct {
	ID                  string     `json:"id"`
	ExternalID          string     `json:"external_id"`
	Amount              float64    `json:"amount"`
	Status              string     `json:"status"`
	Email               string     `json:"email,omitempty"`
	PaymentID           string     `json:"payment_id,omitempty"`
	BankCode            string     `json:"bank_code,omitempty"`
	AccountHolderName   string     `json:"account_holder_name,omitempty"`
	AccountNumber       string     `json:"account_number,omitempty"`
	DisbursementID      string     `json:"disbursement_id,omitempty"`
	FailureReason       string     `json:"failure_reason,omitempty"`
	Created             *time.Time `json:"created,omitempty"`
	ExpirationTimestamp *time.Time `json:"expiration_timestamp,omitempty"`
	ClaimedTimestamp    *time.Time `json:"claimed_timestamp,omitempty"`
	FailedTimestamp     *time.Time `json:"failed_timestamp,omitempty"`
	MerchantName        string     `json:"merchant_name,omitempty"`
	PayoutURL           string     `json:"payout_url,omitempty"`
}
