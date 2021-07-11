package ewallet

import (
  "fmt"

  goxendit "github.com/xendit/xendit-go"
  "github.com/xendit/xendit-go/ewallet"

  "github.com/imrenagi/go-payment/invoice"
)


// NewOVO creates CreatePaymentParams from invoice data
func NewOVO(inv *invoice.Invoice) (*ewallet.CreatePaymentParams, error) {

  if inv.BillingAddress == nil {
    return nil, fmt.Errorf("phone number must be provided in billing address")
  }

  if !OvoPhoneValidator.IsValid(inv.BillingAddress.PhoneNumber) {
    return nil, fmt.Errorf("invalid phone number. must be in 08xxxx format")
  }

  return newBuilder(inv).
    SetPaymentMethod(goxendit.EWalletTypeOVO).
    Build()
}



