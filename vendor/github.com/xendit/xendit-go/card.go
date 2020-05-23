package xendit

import "time"

// CardCharge contains data from Xendit's API response of card's charge related requests
// and Create Authorization request.
// For more details see https://xendit.github.io/apireference/?bash#create-charge
// and https://xendit.github.io/apireference/?bash#create-authorization.
// For documentation of subpackage card, checkout https://pkg.go.dev/github.com/xendit/xendit-go/card
type CardCharge struct {
	ID                    string     `json:"id"`
	Status                string     `json:"status"`
	MerchantID            string     `json:"merchant_id"`
	Created               *time.Time `json:"created"`
	BusinessID            string     `json:"business_id"`
	AuthorizedAmount      float64    `json:"authorized_amount"`
	ExternalID            string     `json:"external_id"`
	MerchantReferenceCode string     `json:"merchant_reference_code"`
	ChargeType            string     `json:"charge_type"`
	CardBrand             string     `json:"card_brand"`
	MaskedCardNumber      string     `json:"masked_card_number"`
	CaptureAmount         float64    `json:"capture_amount,omitempty"`
	ECI                   string     `json:"eci,omitempty"`
	FailureReason         string     `json:"failure_reason,omitempty"`
	CardType              string     `json:"card_type,omitempty"`
	BankReconciliationID  string     `json:"bank_reconciliation_id,omitempty"`
	Descriptor            string     `json:"descriptor,omitempty"`
	MidLabel              string     `json:"mid_label,omitempty"`
	Currency              string     `json:"currency,omitempty"`
}

// CardRefund contains data from Xendit's API response of card's Create Refund request.
// For more details see https://xendit.github.io/apireference/?bash#CreateRefund.
// For documentation of subpackage card, checkout https://pkg.go.dev/github.com/xendit/xendit-go/card
type CardRefund struct {
	ID                 string     `json:"id"`
	Updated            *time.Time `json:"updated"`
	Created            *time.Time `json:"created"`
	CreditCardChargeID string     `json:"credit_card_charge_id"`
	UserID             string     `json:"user_id"`
	Amount             float64    `json:"amount"`
	ExternalID         string     `json:"external_id"`
	Currency           string     `json:"currency"`
	Status             string     `json:"status"`
	FeeRefundAmount    float64    `json:"fee_refund_amount"`
	FailureReason      string     `json:"failure_reason"`
}

// CardReverseAuthorization contains data from Xendit's API response of card's Reverse Authorization request.
// For more details see https://xendit.github.io/apireference/?bash#reverse-authorization.
// For documentation of subpackage card, checkout https://pkg.go.dev/github.com/xendit/xendit-go/card
type CardReverseAuthorization struct {
	ID                 string     `json:"id"`
	ExternalID         string     `json:"external_id"`
	CreditCardChargeID string     `json:"credit_card_charge_id"`
	BusinessID         string     `json:"business_id"`
	Amount             float64    `json:"amount"`
	Status             string     `json:"status"`
	Created            *time.Time `json:"created"`
	Currency           string     `json:"currency,omitempty"`
}
