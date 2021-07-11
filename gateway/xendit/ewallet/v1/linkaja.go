package ewallet

import (
  "os"

  goxendit "github.com/xendit/xendit-go"
  "github.com/xendit/xendit-go/ewallet"

  "github.com/imrenagi/go-payment/invoice"
)

// NewLinkAja create xendit payment request for LinkAja
func NewLinkAja(inv *invoice.Invoice) (*ewallet.CreatePaymentParams, error) {

  return newBuilder(inv).
    SetPaymentMethod(goxendit.EWalletTypeLINKAJA).
    SetCallback(os.Getenv("LINKAJA_LEGACY_CALLBACK_URL")).
    SetRedirect(os.Getenv("LINKAJA_LEGACY_REDIRECT_URL")).
    Build()
}

