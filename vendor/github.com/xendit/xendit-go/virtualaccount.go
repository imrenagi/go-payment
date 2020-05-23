package xendit

import (
	"time"
)

// VirtualAccount contains data from Xendit's API response of virtual account related requests.
// For more details see https://xendit.github.io/apireference/?bash#virtual-accounts.
// For documentation of subpackage virtualaccount, checkout https://pkg.go.dev/github.com/xendit/xendit-go/virtualaccount
type VirtualAccount struct {
	OwnerID         string     `json:"owner_id"`
	ExternalID      string     `json:"external_id"`
	BankCode        string     `json:"bank_code"`
	MerchantCode    string     `json:"merchant_code"`
	Name            string     `json:"name"`
	AccountNumber   string     `json:"account_number"`
	IsClosed        *bool      `json:"is_closed"`
	ID              string     `json:"id"`
	IsSingleUse     *bool      `json:"is_single_use"`
	Status          string     `json:"status"`
	Currency        string     `json:"currency"`
	ExpirationDate  *time.Time `json:"expiration_date"`
	SuggestedAmount float64    `json:"suggested_amount,omitempty"`
	ExpectedAmount  float64    `json:"expected_amount,omitempty"`
	Description     string     `json:"description,omitempty"`
}

// VirtualAccountBank contains data from Xendit's API response of Get Virtual Account Banks.
type VirtualAccountBank struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// VirtualAccountPayment contains data from Xendit's API response of Get Fixed Virtual Account Payment.
type VirtualAccountPayment struct {
	ID                       string     `json:"id"`
	PaymentID                string     `json:"payment_id"`
	CallbackVirtualAccountID string     `json:"callback_virtual_account_id"`
	ExternalID               string     `json:"external_id"`
	AccountNumber            string     `json:"account_number"`
	BankCode                 string     `json:"bank_code"`
	Amount                   float64    `json:"amount"`
	TransactionTimestamp     *time.Time `json:"transaction_timestamp"`
	MerchantCode             string     `json:"merchant_code"`
	Currency                 string     `json:"currency"`
}
