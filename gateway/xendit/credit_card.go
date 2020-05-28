package xendit

import (
	xinvoice "github.com/xendit/xendit-go/invoice"
)

func NewCreditCardInvoice(rb *InvoiceRequestBuilder) (*CreditCardInvoice, error) {
	return &CreditCardInvoice{
		rb: rb,
	}, nil
}

type CreditCardInvoice struct {
	rb *InvoiceRequestBuilder
}

func (o *CreditCardInvoice) Build() (*xinvoice.CreateParams, error) {
	o.rb.AddPaymentMethod("CREDIT_CARD")
	req, err := o.rb.Build()
	if err != nil {
		return nil, err
	}
	return req, nil
}
