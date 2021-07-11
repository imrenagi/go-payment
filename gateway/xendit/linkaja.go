package xendit

import xinvoice "github.com/xendit/xendit-go/invoice"

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
