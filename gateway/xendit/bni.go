package xendit

import (
	xinvoice "github.com/xendit/xendit-go/invoice"
)

func NewBNIVAInvoice(rb *InvoiceRequestBuilder) (*BNIVAInvoice, error) {
	return &BNIVAInvoice{
		rb: rb,
	}, nil
}

type BNIVAInvoice struct {
	rb *InvoiceRequestBuilder
}

func (o *BNIVAInvoice) Build() (*xinvoice.CreateParams, error) {
	o.rb.AddPaymentMethod("BNI")
	req, err := o.rb.Build()
	if err != nil {
		return nil, err
	}
	return req, nil
}
