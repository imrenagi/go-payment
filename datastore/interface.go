package datastore

import (
	"context"

	"github.com/imrenagi/go-payment/gateway/midtrans"
	"github.com/imrenagi/go-payment/invoice"
)

// MidtransTransactionStatusRepository is an interface for
// the storage of midtrans transaction status.
type MidtransTransactionStatusRepository interface {
	Save(ctx context.Context, status *midtrans.TransactionStatus) error
	FindByOrderID(ctx context.Context, orderID string) (*midtrans.TransactionStatus, error)
}

// InvoiceRepository is an interface for invoice storage
type InvoiceRepository interface {
	FindByNumber(ctx context.Context, number string) (*invoice.Invoice, error)
	Save(ctx context.Context, invoice *invoice.Invoice) error
}
