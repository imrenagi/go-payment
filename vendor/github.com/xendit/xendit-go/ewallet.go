package xendit

import "time"

// EWalletTypeEnum constants are the available e-wallet type
type EWalletTypeEnum string

// This consists the values that EWalletTypeEnum can take
const (
	EWalletTypeOVO     EWalletTypeEnum = "OVO"
	EWalletTypeDANA    EWalletTypeEnum = "DANA"
	EWalletTypeLINKAJA EWalletTypeEnum = "LINKAJA"
)

// EWallet contains data from Xendit's API response of e-wallet related requests.
// For more details see https://xendit.github.io/apireference/?bash#ewallets.
// For documentation of subpackage ewallet, checkout https://pkg.go.dev/github.com/xendit/xendit-go/ewallet
type EWallet struct {
	EWalletType          EWalletTypeEnum `json:"ewallet_type"`
	ExternalID           string          `json:"external_id"`
	Status               string          `json:"status"`
	Amount               float64         `json:"amount"`
	TransactionDate      *time.Time      `json:"transaction_date,omitempty"`
	CheckoutURL          string          `json:"checkout_url,omitempty"`
	BusinessID           string          `json:"business_id,omitempty"`
	Created              *time.Time      `json:"created,omitempty"`
	EWalletTransactionID string          `json:"e_wallet_transaction_id,omitempty"`
}
