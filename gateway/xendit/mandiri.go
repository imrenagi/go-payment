package xendit

import (
	xinvoice "github.com/xendit/xendit-go/invoice"
)

func NewMandiriVAInvoice(rb *InvoiceRequestBuilder) (*MandiriVAInvoice, error) {
	return &MandiriVAInvoice{
		rb: rb,
	}, nil
}

type MandiriVAInvoice struct {
	rb *InvoiceRequestBuilder
}

func (o *MandiriVAInvoice) Build() (*xinvoice.CreateParams, error) {
	o.rb.AddPaymentMethod("MANDIRI")
	req, err := o.rb.Build()
	if err != nil {
		return nil, err
	}
	return req, nil
}
