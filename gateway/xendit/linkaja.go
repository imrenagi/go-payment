package xendit

import (
	"fmt"
	"os"

	goxendit "github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/ewallet"
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
		SetCallback(fmt.Sprintf("%s/payment/xendit/linkaja/callback", os.Getenv("SERVER_BASE_URL"))).
		SetRedirect(fmt.Sprintf("%s/donate/thanks", os.Getenv("WEB_BASE_URL")))

	req, err := o.rb.Build()
	if err != nil {
		return nil, err
	}

	return req, nil
}
