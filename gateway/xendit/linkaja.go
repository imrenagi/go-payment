package xendit

import (
  "os"

  goxendit "github.com/xendit/xendit-go"
  "github.com/xendit/xendit-go/ewallet"
  xinvoice "github.com/xendit/xendit-go/invoice"
)

// NewLinkAja create xendit payment request for LinkAja
func NewLinkAja(rb *EWalletRequestBuilder) (*LinkAja, error) {
  return &LinkAja{
    rb: rb,
  }, nil
}

// LinkAja ...
type LinkAja struct {
  rb *EWalletRequestBuilder
}

// Build ...
func (o *LinkAja) Build() (*ewallet.CreatePaymentParams, error) {

  o.rb.SetPaymentMethod(goxendit.EWalletTypeLINKAJA).
    SetCallback(os.Getenv("LINKAJA_LEGACY_CALLBACK_URL")).
    SetRedirect(os.Getenv("LINKAJA_LEGACY_REDIRECT_URL"))

  req, err := o.rb.Build()
  if err != nil {
    return nil, err
  }

  return req, nil
}

func NewLinkAjaInvoice(rb *InvoiceRequestBuilder) (*LinkAjaInvoice, error) {
  return &LinkAjaInvoice{
    rb: rb,
  }, nil
}

type LinkAjaInvoice struct {
  rb *InvoiceRequestBuilder
}

func (o *LinkAjaInvoice) Build() (*xinvoice.CreateParams, error) {

  o.rb.AddPaymentMethod("LINKAJA")
  req, err := o.rb.Build()
  if err != nil {
    return nil, err
  }
  return req, nil
}
