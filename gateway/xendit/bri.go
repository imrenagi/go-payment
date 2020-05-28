package xendit

import (
	xinvoice "github.com/xendit/xendit-go/invoice"
)

func NewBRIVAInvoice(rb *InvoiceRequestBuilder) (*BRIVAInvoice, error) {
	return &BRIVAInvoice{
		rb: rb,
	}, nil
}

type BRIVAInvoice struct {
	rb *InvoiceRequestBuilder
}

func (o *BRIVAInvoice) Build() (*xinvoice.CreateParams, error) {
	o.rb.AddPaymentMethod("BRI")
	req, err := o.rb.Build()
	if err != nil {
		return nil, err
	}
	return req, nil
}
