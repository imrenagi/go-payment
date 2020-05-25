package manager

import (
	"context"

	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/gateway/xendit"
	"github.com/imrenagi/go-payment/invoice"
	mgo "github.com/veritrans/go-midtrans"
)

// GenerateInvoiceRequest provide to generate new invoice
type GenerateInvoiceRequest struct {
	Payment struct {
		PaymentType      payment.PaymentType       `json:"payment_type"`
		CreditCardDetail *invoice.CreditCardDetail `json:"credit_card,omitempty"`
	} `json:"payment"`
	Customer struct {
		Name        string `json:"name"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
	} `json:"customer"`
	Item struct {
		Name         string  `json:"name"`
		Category     string  `json:"category"`
		MerchantName string  `json:"merchant"`
		Description  string  `json:"description"`
		Qty          int     `json:"qty"`
		Price        float64 `json:"price"`
		Currency     string  `json:"currency"`
	} `json:"item"`
}

// PayInvoiceRequest provide information which invoice to pay and by using what
// transactionID
type PayInvoiceRequest struct {
	InvoiceNumber string `json:"invoice_number"`
	TransactionID string `json:"transaction_id"`
}

// FailInvoiceRequest provide which invoice that is failed and its reason
type FailInvoiceRequest struct {
	InvoiceNumber string `json:"invoice_number"`
	TransactionID string `json:"transaction_id"`
	Reason        string `json:"reason"`
}

// Interface payment management interface
type Interface interface {
	// return the payment methods available in payment service
	GetPaymentMethods(ctx context.Context, opts ...payment.Option) (*PaymentMethodList, error)

	// return invoice given its invoice number
	GetInvoice(ctx context.Context, number string) (*invoice.Invoice, error)

	// generate new invoice
	GenerateInvoice(ctx context.Context, gir *GenerateInvoiceRequest) (*invoice.Invoice, error)

	// PayInvoice pays an invoice
	PayInvoice(ctx context.Context, pir *PayInvoiceRequest) (*invoice.Invoice, error)

	// ProcessInvoice ...
	ProcessInvoice(ctx context.Context, invoiceNumber string) (*invoice.Invoice, error)

	// FailInvoice make the invoice failed
	FailInvoice(ctx context.Context, fir *FailInvoiceRequest) (*invoice.Invoice, error)
}

// XenditProcessor callback handler for xendit
type XenditProcessor interface {
	ProcessDANACallback(ctx context.Context, dps *xendit.DANAPaymentStatus) error
	ProcessLinkAjaCallback(ctx context.Context, lps *xendit.LinkAjaPaymentStatus) error
	ProcessOVOCallback(ctx context.Context, ops *xendit.OVOPaymentStatus) error
	ProcessXenditInvoicesCallback(ctx context.Context, ips *xendit.InvoicePaymentStatus) error
}

// MidtransProcessor callback handler for midtrans
type MidtransProcessor interface {
	ProcessMidtransCallback(ctx context.Context, mr mgo.Response) error
}
