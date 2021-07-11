package xendit

import xinvoice "github.com/xendit/xendit-go/invoice"

func NewDanaInvoice(rb *InvoiceRequestBuilder) (*DanaInvoice, error) {
  return &DanaInvoice{
    rb: rb,
  }, nil
}

type DanaInvoice struct {
  rb *InvoiceRequestBuilder
}

func (o *DanaInvoice) Build() (*xinvoice.CreateParams, error) {

  o.rb.AddPaymentMethod("DANA")
  req, err := o.rb.Build()
  if err != nil {
    return nil, err
  }
  return req, nil
}
