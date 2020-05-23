package cardlesscredit

import "github.com/xendit/xendit-go"

// CreatePaymentParams contains parameters for CreatePayment
type CreatePaymentParams struct {
	CardlessCreditType xendit.CardlessCreditTypeEnum `json:"cardless_credit_type" validate:"required"`
	ExternalID         string                        `json:"external_id" validate:"required"`
	Amount             float64                       `json:"amount" validate:"required"`
	PaymentType        xendit.PaymentTypeEnum        `json:"payment_type" validate:"required"`
	Items              []Item                        `json:"items" validate:"required"`
	CustomerDetails    CustomerDetails               `json:"customer_details" validate:"required"`
	ShippingAddress    ShippingAddress               `json:"shipping_address" validate:"required"`
	RedirectURL        string                        `json:"redirect_url" validate:"required"`
	CallbackURL        string                        `json:"callback_url" validate:"required"`
}

// Item is data that contained in CreatePaymentParams at Items
type Item struct {
	ID       string  `json:"id" validate:"required"`
	Name     string  `json:"name" validate:"required"`
	Price    float64 `json:"price" validate:"required"`
	Type     string  `json:"type" validate:"required"`
	URL      string  `json:"url" validate:"required"`
	Quantity int     `json:"quantity" validate:"required"`
}

// CustomerDetails is data that contained in CreatePaymentParams at CustomerDetails
type CustomerDetails struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
}

// ShippingAddress is data that contained in CreatePaymentParams at ShippingAddress
type ShippingAddress struct {
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	Address     string `json:"address" validate:"required"`
	City        string `json:"city" validate:"required"`
	PostalCode  string `json:"postal_code" validate:"required"`
	Phone       string `json:"phone" validate:"required"`
	CountryCode string `json:"country_code" validate:"required"`
}
