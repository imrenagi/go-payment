package xendit

// CardlessCreditTypeEnum constants are the available cardless credit types
type CardlessCreditTypeEnum string

// This consists the values that CardlessCreditTypeEnum can take
const (
	CardlessCreditTypeEnumKREDIVO CardlessCreditTypeEnum = "KREDIVO"
)

// PaymentTypeEnum constants are the available payment types
type PaymentTypeEnum string

// This consists the values that PaymentTypeEnum can take
const (
	PaymentTypeEnum30Days   PaymentTypeEnum = "30_days"
	PaymentTypeEnum3Months  PaymentTypeEnum = "3_months"
	PaymentTypeEnum6Months  PaymentTypeEnum = "6_months"
	PaymentTypeEnum12Months PaymentTypeEnum = "12_months"
)

// CardlessCredit contains data from Xendit's API response of cardless credit related requests.
// For more details see https://xendit.github.io/apireference/?bash#cardless-credit.
// For documentation of subpackage cardlesscredit, checkout https://pkg.go.dev/github.com/xendit/xendit-go/cardlesscredit
type CardlessCredit struct {
	RedirectURL        string                 `json:"redirect_url"`
	TransactionID      string                 `json:"transaction_id"`
	OrderID            string                 `json:"order_id"`
	ExternalID         string                 `json:"external_id"`
	CardlessCreditType CardlessCreditTypeEnum `json:"cardless_credit_type"`
}
