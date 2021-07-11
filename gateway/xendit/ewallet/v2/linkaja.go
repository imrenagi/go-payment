package ewallet

import (
  "os"

  "github.com/xendit/xendit-go/ewallet"

  "github.com/imrenagi/go-payment/invoice"
)

// NewLinkAja is factory for LinkAja payment with xendit latest charge API
func NewLinkAja(inv *invoice.Invoice) (*ewallet.CreateEWalletChargeParams, error) {

  props := map[string]string{
    "success_redirect_url": os.Getenv("LINKAJA_SUCCESS_REDIRECT_URL"),
  }

  return newBuilder(inv).
    SetPaymentMethod(EWalletIDLinkAja).
    SetChannelProperties(props).
    Build()
}
