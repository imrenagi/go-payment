package ewallet_test

import (
  "time"

  "github.com/imrenagi/go-payment"
  "github.com/imrenagi/go-payment/invoice"
)

var dummyInv = &invoice.Invoice{
  Model:           payment.Model{},
  Number:          "a-random-invoice-number",
  InvoiceDate:     time.Time{},
  DueDate:         time.Time{},
  PaidAt:          nil,
  Currency:        "IDR",
  SubTotal:        15000,
  LineItems:       []invoice.LineItem{
    {
      Model:        payment.Model{
        ID: 1,
      },
      InvoiceID:    1,
      Name:         "random-item",
      Description:  "just description",
      Category:     "HOME",
      Currency:     "IDR",
      UnitPrice:    15000,
      Qty:          1,
    },
  },
  Payment:         nil,
  BillingAddress: &invoice.BillingAddress{
    PhoneNumber: "+628111231234",
  },
  SubscriptionID:  nil,
}
