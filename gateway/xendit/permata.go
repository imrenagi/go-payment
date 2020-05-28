package xendit

import (
	xinvoice "github.com/xendit/xendit-go/invoice"
)

func NewPermataVAInvoice(rb *InvoiceRequestBuilder) (*PermataVAInvoice, error) {
	return &PermataVAInvoice{
		rb: rb,
	}, nil
}

type PermataVAInvoice struct {
	rb *InvoiceRequestBuilder
}

func (o *PermataVAInvoice) Build() (*xinvoice.CreateParams, error) {
	o.rb.AddPaymentMethod("PERMATA")
	req, err := o.rb.Build()
	if err != nil {
		return nil, err
	}
	return req, nil
}
