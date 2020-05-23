package xendit

import (
	"context"

	"github.com/imrenagi/go-payment/util/localconfig"
	"github.com/xendit/xendit-go"
	xgo "github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/client"
	"github.com/xendit/xendit-go/ewallet"
	"github.com/xendit/xendit-go/invoice"
	recurring "github.com/xendit/xendit-go/recurringpayment"
)

// NewGateway creates new midtrans payment gateway
func NewGateway(creds localconfig.APICredential) *Gateway {

	c := client.New(creds.SecretKey)

	gateway := Gateway{
		client:    c,
		Ewallet:   c.EWallet,
		Invoice:   c.Invoice,
		Recurring: c.RecurringPayment,
	}

	return &gateway
}

// Gateway ...
type Gateway struct {
	client    *client.API
	Ewallet   Ewallet
	Invoice   Invoice
	Recurring Recurring
}

type Ewallet interface {
	CreatePayment(data *ewallet.CreatePaymentParams) (*xgo.EWallet, *xgo.Error)
}

type Invoice interface {
	CreateWithContext(ctx context.Context, data *invoice.CreateParams) (*xgo.Invoice, *xgo.Error)
}

type Recurring interface {
	CreateWithContext(ctx context.Context, data *recurring.CreateParams) (*xgo.RecurringPayment, *xendit.Error)
}
