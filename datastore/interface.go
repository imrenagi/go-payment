package datastore

import (
	"context"

	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/config"
	"github.com/imrenagi/go-payment/gateway/midtrans"
	"github.com/imrenagi/go-payment/invoice"
)

// MidtransTransactionStatusRepository ...
type MidtransTransactionStatusRepository interface {
	Save(ctx context.Context, status *midtrans.TransactionStatus) error
	FindByOrderID(ctx context.Context, orderID string) (*midtrans.TransactionStatus, error)
}

// MidtransCardTokenRepository ...
type MidtransCardTokenRepository interface {
	Save(ctx context.Context, token *midtrans.CardToken) error
	FindAllByUserID(ctx context.Context, userID string) ([]midtrans.CardToken, error)
}

type PaymentMethodRepository interface {
	FindByPaymentType(ctx context.Context, paymentType payment.PaymentType, opts ...payment.PaymentOption) (config.FeeConfigReader, error)
	FindAll(ctx context.Context) (*config.PaymentConfig, error)
}

type InvoiceRepository interface {
	FindByNumber(ctx context.Context, number string) (*invoice.Invoice, error)
	Save(ctx context.Context, invoice *invoice.Invoice) error
}
