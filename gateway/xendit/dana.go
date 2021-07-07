package xendit

import (
	"os"

	goxendit "github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/ewallet"
	xinvoice "github.com/xendit/xendit-go/invoice"
)

// NewDana create xendit payment request for Dana
func NewDana(rb *EWalletRequestBuilder) (*Dana, error) {
	return &Dana{
		rb: rb,
	}, nil
}

// Dana ...
type Dana struct {
	rb *EWalletRequestBuilder
}

// Build ...
func (o *Dana) Build() (*ewallet.CreatePaymentParams, error) {

	o.rb.SetPaymentMethod(goxendit.EWalletTypeDANA).
		SetCallback(os.Getenv("DANA_LEGACY_CALLBACK_URL")).
		SetRedirect(os.Getenv("DANA_LEGACY_REDIRECT_URL"))

	req, err := o.rb.Build()
	if err != nil {
		return nil, err
	}

	return req, nil
}

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

// NewDanaCharge is factory for Dana payment with xendit latest charge API
func NewDanaCharge(rb *EWalletChargeRequestBuilder) (*DanaCharge, error) {
	return &DanaCharge{
		rb:    rb,
	}, nil
}

type DanaCharge struct {
	phone string
	rb    *EWalletChargeRequestBuilder
}

func (o *DanaCharge) Build() (*ewallet.CreateEWalletChargeParams, error) {

	props := map[string]string{
		"success_redirect_url": os.Getenv("DANA_SUCCESS_REDIRECT_URL"),
	}

	o.rb.SetPaymentMethod(EWalletIDDana).
		SetChannelProperties(props)
	return o.rb.Build()
}