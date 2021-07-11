package xendit

import (
  "fmt"

  goxendit "github.com/xendit/xendit-go"
  "github.com/xendit/xendit-go/ewallet"
  xinvoice "github.com/xendit/xendit-go/invoice"
)

// NewOVO create xendit payment request for ovo
func NewOVO(rb *EWalletRequestBuilder) (*OVO, error) {
  return &OVO{
    rb: rb,
  }, nil
}

// OVO ...
type OVO struct {
  rb *EWalletRequestBuilder
}

// Build ...
func (o *OVO) Build() (*ewallet.CreatePaymentParams, error) {
  o.rb.SetPaymentMethod(goxendit.EWalletTypeOVO)
  req, err := o.rb.Build()
  if err != nil {
    return nil, err
  }

  if !OvoPhoneValidator.IsValid(req.Phone) {
    return nil, fmt.Errorf("invalid phone number. must be in 08xxxx format")
  }

  return req, nil
}

func NewOVOInvoice(rb *InvoiceRequestBuilder) (*OVOInvoice, error) {
  return &OVOInvoice{
    rb: rb,
  }, nil
}

type OVOInvoice struct {
  rb *InvoiceRequestBuilder
}

func (o *OVOInvoice) Build() (*xinvoice.CreateParams, error) {

  o.rb.AddPaymentMethod("OVO")
  req, err := o.rb.Build()
  if err != nil {
    return nil, err
  }
  return req, nil
}


