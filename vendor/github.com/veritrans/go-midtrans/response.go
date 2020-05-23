package midtrans

import "time"

// VANumber : bank virtual account number
type VANumber struct {
	Bank     string `json:"bank"`
	VANumber string `json:"va_number"`
}

// Action represents response action
type Action struct {
	Name   string `json:"name"`
	Method string `json:"method"`
	URL    string `json:"url"`
}

// Response after calling the API
type Response struct {
	StatusCode         string     `json:"status_code"`
	StatusMessage      string     `json:"status_message"`
	PermataVaNumber    string     `json:"permata_va_number"`
	SignKey            string     `json:"signature_key"`
	CardToken          string     `json:"token_id"`
	SavedCardToken     string     `json:"saved_token_id"`
	SavedTokenExpAt    string     `json:"saved_token_id_expired_at"`
	SecureToken        bool       `json:"secure_token"`
	Bank               string     `json:"bank"`
	BillerCode         string     `json:"biller_code"`
	BillKey            string     `json:"bill_key"`
	XlTunaiOrderID     string     `json:"xl_tunai_order_id"`
	BIIVaNumber        string     `json:"bii_va_number"`
	ReURL              string     `json:"redirect_url"`
	ECI                string     `json:"eci"`
	ValMessages        []string   `json:"validation_messages"`
	Page               int        `json:"page"`
	TotalPage          int        `json:"total_page"`
	TotalRecord        int        `json:"total_record"`
	FraudStatus        string     `json:"fraud_status"`
	PaymentType        string     `json:"payment_type"`
	OrderID            string     `json:"order_id"`
	TransactionID      string     `json:"transaction_id"`
	TransactionTime    string     `json:"transaction_time"`
	TransactionStatus  string     `json:"transaction_status"`
	GrossAmount        string     `json:"gross_amount"`
	VANumbers          []VANumber `json:"va_numbers"`
	PaymentCode        string     `json:"payment_code"`
	Store              string     `json:"store"`
	MerchantID         string     `json:"merchant_id"`
	MaskedCard         string     `json:"masked_card"`
	Currency           string     `json:"currency"`
	CardType           string     `json:"card_type"`
	Actions            []Action   `json:"actions"`
	RefundChargebackID int        `json:"refund_chargeback_id"`
	RefundAmount       string     `json:"refund_amount"`
	RefundKey          string     `json:"refund_key"`
	Refunds            []Refund   `json:"refunds"`
}

// Refund Details
type Refund struct {
	RefundChargebackID int    `json:"refund_chargeback_id"`
	RefundAmount       string `json:"refund_amount"`
	Reason             string `json:"reason"`
	RefundKey          string `json:"refund_key"`
	RefundMethod       string `json:"refund_method"`
	BankConfirmedAt    string `json:"bank_confirmed_at"`
	CreatedAt          string `json:"created_at"`
}

// ResponseWithMap after calling the API
type ResponseWithMap map[string]interface{}

// SnapResponse : Response after calling the Snap API
type SnapResponse struct {
	StatusCode    string   `json:"status_code"`
	Token         string   `json:"token"`
	RedirectURL   string   `json:"redirect_url"`
	ErrorMessages []string `json:"error_messages"`
}

// IrisBeneficiaryBanksResponse : Show list of supported banks in IRIS. https://iris-docs.midtrans.com/#list-banks
type IrisBeneficiaryBanksResponse struct {
	BeneficiaryBanks []IrisBeneficiaryBankResponse `json:"beneficiary_banks"`
	StatusCode       string                        `json:"status_code"`
}

// IrisBeneficiaryBankResponse : Represent Beneficiary bank response payload
type IrisBeneficiaryBankResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// IrisBeneficiaries : Beneficiaries request (create, update, list)
// https://iris-docs.midtrans.com/#create-beneficiaries
// https://iris-docs.midtrans.com/#update-beneficiaries
// https://iris-docs.midtrans.com/#list-beneficiaries
type IrisBeneficiaries struct {
	Name      string `json:"name"`
	Account   string `json:"account"`
	Bank      string `json:"bank"`
	AliasName string `json:"alias_name"`
	Email     string `json:"email"`
}

// IrisBeneficiariesResponse : Represent Beneficiaries response payload
type IrisBeneficiariesResponse struct {
	Status     string   `json:"status"`
	StatusCode string   `json:"status_code"`
	Errors     []string `json:"errors"`
}

// IrisCreatePayoutResponse : Represent Create payout response payload
type IrisCreatePayoutResponse struct {
	Payouts      []IrisCreatePayoutDetailResponse `json:"payouts"`
	ErrorMessage string                           `json:"error_message"`
	Errors       []string                         `json:"errors"`
}

// IrisCreatePayoutDetailResponse : Represent Create payout detail response payload
type IrisCreatePayoutDetailResponse struct {
	Status      string `json:"status"`
	ReferenceNo string `json:"reference_no"`
}

// IrisApprovePayoutResponse : Represent Approve payout response payload
type IrisApprovePayoutResponse struct {
	Status       string   `json:"status"`
	ErrorMessage string   `json:"error_message"`
	Errors       []string `json:"errors"`
}

// IrisRejectPayoutResponse : Represent Reject payout response payload
type IrisRejectPayoutResponse struct {
	Status       string   `json:"status"`
	ErrorMessage string   `json:"error_message"`
	Errors       []string `json:"errors"`
}

// IrisPayoutDetailResponse : Represent Payout detail response payload
type IrisPayoutDetailResponse struct {
	Amount             string    `json:"amount"`
	BeneficiaryName    string    `json:"beneficiary_name"`
	BeneficiaryAccount string    `json:"beneficiary_account"`
	Bank               string    `json:"bank"`
	ReferenceNo        string    `json:"reference_no"`
	Notes              string    `json:"notes"`
	BeneficiaryEmail   string    `json:"beneficiary_email"`
	Status             string    `json:"status"`
	CreatedBy          string    `json:"created_by"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	ErrorMessage       string    `json:"error_message"`
	Errors             string    `json:"errors"`
}

// IrisBankAccountDetailResponse : Represent Bank account detail payload
type IrisBankAccountDetailResponse struct {
	AccountName  string                              `json:"account_name"`
	AccountNo    string                              `json:"account_no"`
	BankName     string                              `json:"bank_name"`
	ErrorMessage string                              `json:"error_message"`
	Errors       *IrisBankAccountDetailErrorResponse `json:"errors"`
}

// IrisBankAccountDetailErrorResponse : Represent Bank account detail error payload
type IrisBankAccountDetailErrorResponse struct {
	Account []string `json:"account"`
	Bank    []string `json:"bank"`
}

// IrisBalanceResponse : Represent balance detail response payload
type IrisBalanceResponse struct {
	Balance string `json:"balance"`
}

type SubscribeResponse struct {
	ID            string                 `json:"id"`
	CreatedAt     string                 `json:"created_at"`
	Status        string                 `json:"status"`
	Name          string                 `json:"name"`
	Amount        string                 `json:"amount"`
	Currency      string                 `json:"currency"`
	Token         string                 `json:"token"`
	PaymentType   PaymentType            `json:"payment_type"`
	Schedule      ScheduleDetailResponse `json:"schedule"`
	StatusMessage string                 `json:"status_message"`
}

type ScheduleDetailResponse struct {
	Interval            int    `json:"interval"`
	IntervalUnit        string `json:"interval_unit"`
	StartTime           string `json:"start_time"`
	PreviousExecutionAt string `json:"previous_execution_at"`
	NextExecutionAt     string `json:"next_execution_at"`
}
