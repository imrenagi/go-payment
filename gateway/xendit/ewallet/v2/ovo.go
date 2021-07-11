package ewallet

import (
  "fmt"

  "github.com/xendit/xendit-go/ewallet"

  "github.com/imrenagi/go-payment/invoice"
)

// NewOVO is factory for OVO payment with xendit latest charge API
func NewOVO(inv *invoice.Invoice) (*ewallet.CreateEWalletChargeParams, error) {

  if inv.BillingAddress == nil {
    return nil, fmt.Errorf("customer phone number is required")
  }

  if !OvoChargePhoneValidator.IsValid(inv.BillingAddress.PhoneNumber) {
    return nil, fmt.Errorf("invalid phone format. must be in +628xxxxxx format")
  }

  props := map[string]string{
    "mobile_number": inv.BillingAddress.PhoneNumber,
  }

  return newEWalletChargeRequestBuilder(inv).
    SetPaymentMethod(EWalletIDOVO).
    SetChannelProperties(props).
    Build()
}

