package service

import (
	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/invoice"
)

type CreateDonationCommand struct {
	Payment struct {
		PaymentType      payment.PaymentType       `json:"payment_type"`
		CreditCardDetail *invoice.CreditCardDetail `json:"credit_card,omitempty"`
	} `json:"payment"`
	DonaturName string  `json:"name"`
	Message     string  `json:"message"`
	Value       float64 `json:"value"`
	Email       string  `json:"email"`
	PhoneNumber string  `json:"phone_number"`
}

type PayInvoiceCommand struct {
	TransactionID string `json:"transaction_id"`
}
