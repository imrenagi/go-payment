package xeninvoice

import (
  xinvoice "github.com/xendit/xendit-go/invoice"

  "github.com/imrenagi/go-payment/invoice"
)

func NewPermataVA(inv *invoice.Invoice) (*xinvoice.CreateParams, error) {
  return newBuilder(inv).
    AddPaymentMethod("PERMATA").
    Build()
}
