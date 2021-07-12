package snap_test

import (
	"time"

	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/invoice"
)

var fakeInvDate = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
var fakeDueDate = time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)

var dummyInv = &invoice.Invoice{
	Model:          payment.Model{},
	Number:         "a-random-invoice-number",
	InvoiceDate:    fakeInvDate,
	DueDate:        fakeDueDate,
	PaidAt:         nil,
	Currency:       "IDR",
	SubTotal:       15000,
	Discount:       1000,
	Tax:            200,
	ServiceFee:     500,
	InstallmentFee: 1000,
	LineItems: []invoice.LineItem{
		{
			Model: payment.Model{
				ID: 1,
			},
			InvoiceID:   1,
			Name:        "random-item",
			Description: "just description",
			Category:    "HOME",
			Currency:    "IDR",
			UnitPrice:   15000,
			Qty:         1,
		},
	},
	Payment: nil,
	BillingAddress: &invoice.BillingAddress{
		FullName:    "John Doe",
		Email:       "foo@bar.com",
		PhoneNumber: "08111231234",
	},
	SubscriptionID: nil,
}
