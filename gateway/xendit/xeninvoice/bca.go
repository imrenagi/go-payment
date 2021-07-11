package xeninvoice

import (
  xinvoice "github.com/xendit/xendit-go/invoice"

  "github.com/imrenagi/go-payment/invoice"
)

func NewBCAVA(inv *invoice.Invoice) (*xinvoice.CreateParams, error) {
  return newBuilder(inv).
    AddPaymentMethod("BCA").
    Build()
}
