package xendit

import (
	"os"

	goxendit "github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/ewallet"
	xinvoice "github.com/xendit/xendit-go/invoice"

	"github.com/imrenagi/go-payment/invoice"
)

// NewLinkAja create xendit payment request for LinkAja
func NewLinkAja(rb *EWalletRequestBuilder) (*LinkAja, error) {
	return &LinkAja{
		rb: rb,
	}, nil
}

// LinkAja ...
type LinkAja struct {
	rb *EWalletRequestBuilder
}

// Build ...
func (o *LinkAja) Build() (*ewallet.CreatePaymentParams, error) {

	o.rb.SetPaymentMethod(goxendit.EWalletTypeLINKAJA).
		SetCallback(os.Getenv("LINKAJA_LEGACY_CALLBACK_URL")).
		SetRedirect(os.Getenv("LINKAJA_LEGACY_REDIRECT_URL"))

	req, err := o.rb.Build()
	if err != nil {
		return nil, err
	}

	return req, nil
}

func NewLinkAjaInvoice(rb *InvoiceRequestBuilder) (*LinkAjaInvoice, error) {
	return &LinkAjaInvoice{
		rb: rb,
	}, nil
}

type LinkAjaInvoice struct {
	rb *InvoiceRequestBuilder
}

func (o *LinkAjaInvoice) Build() (*xinvoice.CreateParams, error) {

	o.rb.AddPaymentMethod("LINKAJA")
	req, err := o.rb.Build()
	if err != nil {
		return nil, err
	}
	return req, nil
}

// NewLinkAjaCharge is factory for LinkAja payment with xendit latest charge API
func NewLinkAjaCharge(inv *invoice.Invoice) (*ewallet.CreateEWalletChargeParams, error) {

	props := map[string]string{
		"success_redirect_url": os.Getenv("LINKAJA_SUCCESS_REDIRECT_URL"),
	}

	return newEWalletChargeRequestBuilder(inv).
		SetPaymentMethod(EWalletIDLinkAja).
		SetChannelProperties(props).
		Build()
}
