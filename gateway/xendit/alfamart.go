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

func (ai *AlfamartInvoice) Build() (*xinvoice.CreateParams, error) {
	ai.rb.AddPaymentMethod("ALFAMART")
	req, err := ai.rb.Build()
	if err != nil {
		return nil, err
	}
	return req, nil
}
