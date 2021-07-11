package ewallet_test

import (
  "time"

  "github.com/imrenagi/go-payment"
  "github.com/imrenagi/go-payment/invoice"
)

var fakeDueDate = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC )

var dummyInv = &invoice.Invoice{
  Model:       payment.Model{},
  Number:      "a-random-invoice-number",
  InvoiceDate: time.Time{},
  DueDate:     fakeDueDate,
  PaidAt:      nil,
  Currency:    "IDR",
  SubTotal:    15000,
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
    PhoneNumber: "08111231234",
  },
  SubscriptionID:  nil,
}


var incorrectPhoneDummyInv = &invoice.Invoice{
  Model:       payment.Model{},
  Number:      "a-random-invoice-number",
  InvoiceDate: time.Time{},
  DueDate:     fakeDueDate,
  PaidAt:      nil,
  Currency:    "IDR",
  SubTotal:    15000,
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
