package midtrans

// ItemDetail : Represent the transaction details
type ItemDetail struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Price        int64  `json:"price"`
	Qty          int32  `json:"quantity"`
	Brand        string `json:"brand,omitempty"`
	Category     string `json:"category,omitempty"`
	MerchantName string `json:"merchant_name,omitempty"`
}

// CustAddress : Represent the customer address
type CustAddress struct {
	FName       string `json:"first_name"`
	LName       string `json:"last_name"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	City        string `json:"city"`
	Postcode    string `json:"postal_code"`
	CountryCode string `json:"country_code"`
}

// CustDetail : Represent the customer detail
type CustDetail struct {
	// first name
	FName string `json:"first_name,omitempty"`

	// last name
	LName string `json:"last_name,omitempty"`

	Email    string       `json:"email,omitempty"`
	Phone    string       `json:"phone,omitempty"`
	BillAddr *CustAddress `json:"billing_address,omitempty"`
	ShipAddr *CustAddress `json:"customer_address,omitempty"`
}

// TransactionDetails : Represent transaction details
type TransactionDetails struct {
	OrderID  string `json:"order_id"`
	GrossAmt int64  `json:"gross_amount"`
}

// ExpiryDetail : Represent SNAP expiry details
type ExpiryDetail struct {
	StartTime string `json:"start_time,omitempty"`
	Unit      string `json:"unit"`
	Duration  int64  `json:"duration"`
}

// InstallmentTermsDetail : Represent installment available banks
type InstallmentTermsDetail struct {
	Bni     []int8 `json:"bni,omitempty"`
	Mandiri []int8 `json:"mandiri,omitempty"`
	Cimb    []int8 `json:"cimb,omitempty"`
	Mega    []int8 `json:"mega,omitempty"`
	Bca     []int8 `json:"bca,omitempty"`
	Bri     []int8 `json:"bri,omitempty"`
	Maybank []int8 `json:"maybank,omitempty"`
	Offline []int8 `json:"offline,omitempty"`
}

// InstallmentDetail : Represent installment detail
type InstallmentDetail struct {
	Required bool                    `json:"required"`
	Terms    *InstallmentTermsDetail `json:"terms"`
}

// CreditCardDetail : Represent credit card detail
type CreditCardDetail struct {
	Secure          bool               `json:"secure,omitempty"`
	TokenID         string             `json:"token_id"`
	Bank            string             `json:"bank,omitempty"`
	Bins            []string           `json:"bins,omitempty"`
	Installment     *InstallmentDetail `json:"installment,omitempty"`
	InstallmentTerm int8               `json:"installment_term,omitempty"`
	Type            string             `json:"type,omitempty"`
	// indicate if generated token should be saved for next charge
	SaveTokenID          bool   `json:"save_token_id,omitempty"`
	SavedTokenIDExpireAt string `json:"saved_token_id_expired_at,omitempty"`
	Authentication       string `json:"authentication,omitempty"`
}

// PermataBankTransferDetail : Represent Permata bank_transfer detail
type PermataBankTransferDetail struct {
	Bank Bank `json:"bank"`
}

// BCABankTransferLangDetail : Represent BCA bank_transfer lang detail
type BCABankTransferLangDetail struct {
	LangID string `json:"id,omitempty"`
	LangEN string `json:"en,omitempty"`
}

/*
   Example of usage syntax:
   midtrans.BCABankTransferDetail{
       FreeText: {
           Inquiry: []midtrans.BCABankTransferLangDetail{
               {
                   LangEN: "Test",
                   LangID: "Coba",
               },
           },
       },
   }
*/

// BCABankTransferDetailFreeText : Represent BCA bank_transfer detail free_text
type BCABankTransferDetailFreeText struct {
	Inquiry []BCABankTransferLangDetail `json:"inquiry,omitempty"`
	Payment []BCABankTransferLangDetail `json:"payment,omitempty"`
}

// BCABankTransferDetail : Represent BCA bank_transfer detail
type BCABankTransferDetail struct {
	Bank     Bank                          `json:"bank"`
	VaNumber string                        `json:"va_number"`
	FreeText BCABankTransferDetailFreeText `json:"free_text"`
}

// MandiriBillBankTransferDetail : Represent Mandiri Bill bank_transfer detail
type MandiriBillBankTransferDetail struct {
	BillInfo1 string `json:"bill_info1,omitempty"`
	BillInfo2 string `json:"bill_info2,omitempty"`
}

// BankTransferDetail : Represent bank_transfer detail
type BankTransferDetail struct {
	Bank     Bank                           `json:"bank,omitempty"`
	VaNumber string                         `json:"va_number,omitempty"`
	FreeText *BCABankTransferDetailFreeText `json:"free_text,omitempty"`
	*MandiriBillBankTransferDetail
}

// BCAKlikPayDetail : Represent Internet Banking for BCA KlikPay
type BCAKlikPayDetail struct {
	// 1 = normal, 2 = installment, 3 = normal + installment
	Type    string `json:"type"`
	Desc    string `json:"description"`
	MiscFee int64  `json:"misc_fee,omitempty"`
}

// BCAKlikBCADetail : Represent BCA KlikBCA detail
type BCAKlikBCADetail struct {
	Desc   string `json:"description"`
	UserID string `json:"user_id"`
}

// MandiriClickPayDetail : Represent Mandiri ClickPay detail
type MandiriClickPayDetail struct {
	TokenID string `json:"token_id"`
	Input1  string `json:"input1"`
	Input2  string `json:"input2"`
	Input3  string `json:"input3"`
	Token   string `json:"token"`
}

// CIMBClicksDetail : Represent CIMB Clicks detail
type CIMBClicksDetail struct {
	Desc string `json:"description"`
}

// TelkomselCashDetail : Represent Telkomsel Cash detail
type TelkomselCashDetail struct {
	Promo      bool   `json:"promo"`
	IsReversal int8   `json:"is_reversal"`
	Customer   string `json:"customer"`
}

// IndosatDompetkuDetail : Represent Indosat Dompetku detail
type IndosatDompetkuDetail struct {
	MSISDN string `json:"msisdn"`
}

// MandiriEcashDetail : Represent Mandiri e-Cash detail
type MandiriEcashDetail struct {
	Desc string `json:"description"`
}

// ConvStoreDetail : Represent cstore detail
type ConvStoreDetail struct {
	Store   string `json:"store"`
	Message string `json:"message"`
}

// GopayDetail : Represent gopay detail
type GopayDetail struct {
	EnableCallback bool   `json:"enable_callback"`
	CallbackUrl    string `json:"callback_url"`
}

// ChargeReq : Represent Charge request payload
type ChargeReq struct {
	PaymentType        PaymentType        `json:"payment_type"`
	TransactionDetails TransactionDetails `json:"transaction_details"`

	CreditCard                    *CreditCardDetail              `json:"credit_card,omitempty"`
	BankTransfer                  *BankTransferDetail            `json:"bank_transfer,omitempty"`
	MandiriBillBankTransferDetail *MandiriBillBankTransferDetail `json:"echannel,omitempty"`
	BCAKlikPay                    *BCAKlikPayDetail              `json:"bca_klikpay,omitempty"`
	BCAKlikBCA                    *BCAKlikBCADetail              `json:"bca_klikbca,omitempty"`
	MandiriClickPay               *MandiriClickPayDetail         `json:"mandiri_clickpay,omitempty"`
	MandiriEcash                  *MandiriEcashDetail            `json:"mandiri_ecash,omitempty"`
	CIMBClicks                    *CIMBClicksDetail              `json:"cimb_clicks,omitempty"`
	TelkomselCash                 *TelkomselCashDetail           `json:"telkomsel_cash,omitempty"`
	IndosatDompetku               *IndosatDompetkuDetail         `json:"indosat_dompetku,omitempty"`
	CustomerDetail                *CustDetail                    `json:"customer_details,omitempty"`
	ConvStore                     *ConvStoreDetail               `json:"cstore,omitempty"`
	Gopay                         *GopayDetail                   `json:"gopay,omitempty"`

	Items        *[]ItemDetail `json:"item_details,omitempty"`
	CustField1   string        `json:"custom_field1,omitempty"`
	CustField2   string        `json:"custom_field2,omitempty"`
	CustField3   string        `json:"custom_field3,omitempty"`
	CustomExpiry *CustomExpiry `json:"custom_expiry,omitempty"`
}

// ChargeReqWithMap : Represent Charge request with map payload
type ChargeReqWithMap map[string]interface{}

// SnapReq : Represent SNAP API request payload
type SnapReq struct {
	TransactionDetails TransactionDetails `json:"transaction_details"`
	EnabledPayments    []PaymentType      `json:"enabled_payments"`
	Items              *[]ItemDetail      `json:"item_details,omitempty"`
	CustomerDetail     *CustDetail        `json:"customer_details,omitempty"`
	Expiry             *ExpiryDetail      `json:"expiry,omitempty"`
	CreditCard         *CreditCardDetail  `json:"credit_card,omitempty"`
	Gopay              *GopayDetail       `json:"gopay,omitempty"`
	CustomField1       string             `json:"custom_field1"`
	CustomField2       string             `json:"custom_field2"`
	CustomField3       string             `json:"custom_field3"`
}

// CustomExpiry : Represent Core API custom_expiry
type CustomExpiry struct {
	OrderTime      string `json:"order_time,omitempty"`
	ExpiryDuration int    `json:"expiry_duration,omitempty"`
	Unit           string `json:"unit,omitempty"`
}

// SnapReqWithMap : Represent snap request with map payload
type SnapReqWithMap map[string]interface{}

// CaptureReq : Represent Capture request payload
type CaptureReq struct {
	TransactionID string  `json:"transaction_id"`
	GrossAmt      float64 `json:"gross_amount"`
}

// IrisCreatePayoutReq : Represent Create Payout request payload
type IrisCreatePayoutReq struct {
	Payouts []IrisCreatePayoutDetailReq `json:"payouts"`
}

// IrisCreatePayoutDetailReq : Represent Create Payout detail payload
type IrisCreatePayoutDetailReq struct {
	BeneficiaryName    string `json:"beneficiary_name"`
	BeneficiaryAccount string `json:"beneficiary_account"`
	BeneficiaryBank    string `json:"beneficiary_bank"`
	BeneficiaryEmail   string `json:"beneficiary_email"`
	Amount             string `json:"amount"`
	Notes              string `json:"notes"`
}

// IrisApprovePayoutReq : Represent Approve Payout payload
type IrisApprovePayoutReq struct {
	ReferenceNo []string `json:"reference_nos"`
	OTP         string   `json:"otp"`
}

// IrisRejectPayoutReq : Represent Reject Payout payload
type IrisRejectPayoutReq struct {
	ReferenceNo  []string `json:"reference_nos"`
	RejectReason string   `json:"reject_reason"`
}

// RefundReq : Represent Refund request payload
type RefundReq struct {
	RefundKey string `json:"refund_key"`
	Amount    int64  `json:"amount"`
	Reason    string `json:"reason"`
}

// SubscribeReq : Represent Subscribe object payload (request and response)
type SubscribeReq struct {
	Name        string            `json:"name"`
	Amount      string            `json:"amount"`
	Currency    string            `json:"currency"`
	Token       string            `json:"token"`
	PaymentType PaymentType       `json:"payment_type"`
	Schedule    ScheduleDetailReq `json:"schedule"`
}

// ScheduleDetailReq : Represent Schedule object payload
type ScheduleDetailReq struct {
	Interval     int    `json:"interval"`
	IntervalUnit string `json:"interval_unit"`
}
