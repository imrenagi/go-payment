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
		callbackToken: creds.CallbackToken,
		client:        c,
		Ewallet:       c.EWallet,
		Invoice:       c.Invoice,
		Recurring:     c.RecurringPayment,
	}

	return &gateway
}

// Gateway ...
type Gateway struct {
	callbackToken string
	client        *client.API
	Ewallet       xEwallet
	Invoice       xInvoice
	Recurring     xRecurring
}

// NotificationValidationKey returns xendit callback authentication token
func (g Gateway) NotificationValidationKey() string {
	return g.callbackToken
}

type xEwallet interface {
	CreatePayment(data *ewallet.CreatePaymentParams) (*xgo.EWallet, *xgo.Error)
}

type xInvoice interface {
	CreateWithContext(ctx context.Context, data *invoice.CreateParams) (*xgo.Invoice, *xgo.Error)
}

type xRecurring interface {
	CreateWithContext(ctx context.Context, data *recurring.CreateParams) (*xgo.RecurringPayment, *xendit.Error)
	PauseWithContext(ctx context.Context, data *recurring.PauseParams) (*xgo.RecurringPayment, *xendit.Error)
	ResumeWithContext(ctx context.Context, data *recurring.ResumeParams) (*xgo.RecurringPayment, *xendit.Error)
	StopWithContext(ctx context.Context, data *recurring.StopParams) (*xgo.RecurringPayment, *xendit.Error)
}
