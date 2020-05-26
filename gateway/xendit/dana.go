package xendit

import (
	"fmt"
	"os"

	goxendit "github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/ewallet"
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
		SetCallback(fmt.Sprintf("%s/payment/xendit/dana/callback", os.Getenv("SERVER_BASE_URL"))).
		SetRedirect(fmt.Sprintf("%s%s", os.Getenv("WEB_BASE_URL"), os.Getenv("SUCCESS_REDIRECT_PATH")))

	req, err := o.rb.Build()
	if err != nil {
		return nil, err
	}

	return req, nil
}
