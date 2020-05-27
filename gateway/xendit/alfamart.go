package xendit

import (
	xinvoice "github.com/xendit/xendit-go/invoice"
)

func NewAlfamartInvoice(rb *InvoiceRequestBuilder) (*AlfamartInvoice, error) {
	return &AlfamartInvoice{
		rb: rb,
	}, nil
}

type AlfamartInvoice struct {
	rb *InvoiceRequestBuilder
}

func (o *AlfamartInvoice) Build() (*xinvoice.CreateParams, error) {
	o.rb.AddPaymentMethod("ALFAMART")
	req, err := o.rb.Build()
	if err != nil {
		return nil, err
	}
	return req, nil
}
