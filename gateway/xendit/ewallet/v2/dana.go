package ewallet

import (
	"os"

	"github.com/xendit/xendit-go/ewallet"

	"github.com/imrenagi/go-payment/invoice"
)

// NewDana is factory for Dana payment with xendit latest charge API
func NewDana(inv *invoice.Invoice) (*ewallet.CreateEWalletChargeParams, error) {

	successRedirectURL := os.Getenv("DANA_SUCCESS_REDIRECT_URL")
	if inv.SuccessRedirectURL != "" {
		successRedirectURL = inv.SuccessRedirectURL
	}

	props := map[string]string{
		"success_redirect_url": successRedirectURL,
	}

	return newBuilder(inv).
		SetPaymentMethod(EWalletIDDana).
		SetChannelProperties(props).
		Build()
}
