package xendit

import (
	"fmt"

	"github.com/imrenagi/go-payment"
)

// DANAPaymentStatus stores the data sent by xendit while triggering
// any webhook for dana payment
type DANAPaymentStatus struct {
	ExternalID        string  `json:"external_id"`
	Amount            float64 `json:"amount"`
	BusinessID        string  `json:"business_id"`
	EWalletType       string  `json:"ewallet_type"`
	PaymentStatus     string  `json:"payment_status"`
	TransactionDate   string  `json:"transaction_date"`
	CallbackAuthToken string  `json:"callback_authentication_token"`
}

// IsValid checks whether the callback auth token sent by xendit matches the
// authentication token stored on the dashboard
func (s DANAPaymentStatus) IsValid(authKey string) error {
	return checkCallbackToken(authKey, s.CallbackAuthToken)
}

// LinkAjaPaymentStatus stores the data sent by xendit while triggering
// any webhook for linkaja payment
type LinkAjaPaymentStatus struct {
	ExternalID        string  `json:"external_id"`
	Amount            float64 `json:"amount"`
	Status            string  `json:"status"`
	EWalletType       string  `json:"ewallet_type"`
	CallbackAuthToken string  `json:"callback_authentication_token"`
}

// IsValid checks whether the callback auth token sent by xendit matches the
// authentication token stored on the dashboard
func (s LinkAjaPaymentStatus) IsValid(authKey string) error {
	return checkCallbackToken(authKey, s.CallbackAuthToken)
}

func checkCallbackToken(stored, given string) error {
	if stored != given {
		return fmt.Errorf("callback authentication token is invalid, %w", payment.ErrBadRequest)
	}
	return nil
}

// OVOPaymentStatus stores the data sent by xendit while triggering
// any webhook for ovo payment
type OVOPaymentStatus struct {
	Event       string  `json:"event"`
	ID          string  `json:"id"`
	ExternalID  string  `json:"external_id"`
	BusinessID  string  `json:"business_id"`
	Phone       string  `json:"phone"`
	EWalletType string  `json:"ewallet_type"`
	Amount      float64 `json:"amount"`
	FailureCode string  `json:"failure_code"`
	Status      string  `json:"status"`
}

// IsValid always returns no error at least for now since
// we have no idea why xendit is not returning the callback token
// on the notification payload
func (s OVOPaymentStatus) IsValid(authKey string) error {
	return nil
}

// InvoicePaymentStatus stores the data sent by xendit while triggering
// any webhook for xenInvoice
// https://xendit.github.io/apireference/#invoice-callback
type InvoicePaymentStatus struct {
	ID                     string  `json:"id"`
	ExternalID             string  `json:"external_id"`
	UserID                 string  `json:"user_id"`
	PaymentMethod          string  `json:"payment_method"`
	Status                 string  `json:"status"`
	MerchantName           string  `json:"merchant_name"`
	Amount                 float64 `json:"amount"`
	PaidAmount             float64 `json:"paid_amount"`
	BankCode               string  `json:"bank_code"`
	RetailOutletName       string  `json:"retail_outlet_name"`
	EwalletType            string  `json:"ewallet_type"`
	OnDemandLink           string  `json:"on_demand_link"`
	RecurringPaymentID     string  `json:"recurring_payment_id"`
	PaidAt                 string  `json:"paid_at"`
	PayerEmail             string  `json:"payer_email"`
	Description            string  `json:"description"`
	AdjustedReceivedAmount float64 `json:"adjusted_received_amount"`
	FeesPaidAmount         float64 `json:"fees_paid_amount"`
	CreatedAt              string  `json:"created"`
	UpdatedAt              string  `json:"updated"`
	Currency               string  `json:"currency"`
	PaymentChannel         string  `json:"payment_channel"`
	PaymentDestination     string  `json:"payment_destination"`
}

// IsValid always returns no error at least for now since
// we have no idea why xendit is not returning the callback token
// on the notification payload
func (s InvoicePaymentStatus) IsValid(authKey string) error {
	return nil
}
