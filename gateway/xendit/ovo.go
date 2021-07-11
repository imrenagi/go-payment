package xendit

import (
  xinvoice "github.com/xendit/xendit-go/invoice"
)

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

