package ewallet

import (
	"os"

	goxendit "github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/ewallet"

	"github.com/imrenagi/go-payment/invoice"
)

// NewDana create xendit payment request for Dana
func NewDana(inv *invoice.Invoice) (*ewallet.CreatePaymentParams, error) {
	return newBuilder(inv).
		SetPaymentMethod(goxendit.EWalletTypeDANA).
		SetCallback(os.Getenv("DANA_LEGACY_CALLBACK_URL")).
		SetRedirect(os.Getenv("DANA_LEGACY_REDIRECT_URL")).
		Build()
}