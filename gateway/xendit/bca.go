package xendit

import (
	xinvoice "github.com/xendit/xendit-go/invoice"
)

func NewBCAVAInvoice(rb *InvoiceRequestBuilder) (*BCAVAInvoice, error) {
	return &BCAVAInvoice{
		rb: rb,
	}, nil
}

type BCAVAInvoice struct {
	rb *InvoiceRequestBuilder
}

func (o *BCAVAInvoice) Build() (*xinvoice.CreateParams, error) {
	o.rb.AddPaymentMethod("BCA")
	req, err := o.rb.Build()
	if err != nil {
		return nil, err
	}
	return req, nil
}
