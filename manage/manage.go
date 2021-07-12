package manage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/midtrans/midtrans-go/coreapi"

	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/gateway/xendit"
	"github.com/imrenagi/go-payment/invoice"
	"github.com/imrenagi/go-payment/subscription"
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
	Items []struct {
		Name         string  `json:"name"`
		Category     string  `json:"category"`
		MerchantName string  `json:"merchant"`
		Description  string  `json:"description"`
		Qty          int     `json:"qty"`
		Price        float64 `json:"price"`
		Currency     string  `json:"currency"`
	} `json:"items"`
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

// CreateSubscriptionRequest contains data for creating subscription
type CreateSubscriptionRequest struct {
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	Amount            float64 `json:"amount"`
	UserID            string  `json:"user_id"`
	Currency          string  `json:"currency"`
	TotalReccurence   int     `json:"total_recurrence"`
	CardToken         string  `json:"card_token"`
	ChargeImmediately bool    `json:"charge_immediately"`
	Schedule          struct {
		Interval     int        `json:"interval"`
		IntervalUnit string     `json:"interval_unit"`
		StartAt      *time.Time `json:"start_at"`
	} `json:"schedule"`
}

// UnmarshalJSON validates some data
func (csr *CreateSubscriptionRequest) UnmarshalJSON(data []byte) error {
	type Alias CreateSubscriptionRequest
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(csr),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if !aux.ChargeImmediately && aux.Schedule.StartAt == nil {
		return fmt.Errorf("either charge it immediately, or set a start_at to send the invoice later: %w", payment.ErrBadRequest)
	}

	return nil
}

// ToSubscription creates new subscription instance
func (csr CreateSubscriptionRequest) ToSubscription() *subscription.Subscription {
	s := subscription.New()
	s.Name = csr.Name
	s.Description = csr.Description
	s.Amount = csr.Amount
	s.Currency = csr.Currency
	s.UserID = csr.UserID
	s.TotalReccurence = csr.TotalReccurence
	s.CardToken = csr.CardToken
	s.ChargeImmediately = csr.ChargeImmediately
	schedule := subscription.NewSchedule(
		csr.Schedule.Interval,
		subscription.NewIntervalUnit(csr.Schedule.IntervalUnit),
		csr.StartAt(),
	)
	s.Schedule = *schedule
	return s
}

// StartAt return the first time for generating the invoice. If it is charged immediately, start at will be
// now, otherwise it will be the start at send by user
func (csr CreateSubscriptionRequest) StartAt() *time.Time {
	now := time.Now()
	if csr.ChargeImmediately {
		csr.Schedule.StartAt = &now
	}
	return csr.Schedule.StartAt
}

// Interface payment management interface
type Interface interface {
	invoiceI
	subscriptionI

	// return the payment methods available in payment service
	GetPaymentMethods(ctx context.Context, opts ...payment.Option) (*PaymentMethodList, error)
}

type invoiceI interface {
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

type subscriptionI interface {
	// CreateSubscription creates new subscription
	CreateSubscription(ctx context.Context, csr *CreateSubscriptionRequest) (*subscription.Subscription, error)

	// PauseSubscription pause active subscription
	PauseSubscription(ctx context.Context, subsNumber string) (*subscription.Subscription, error)

	// ResumeSubscription resume paused subscription
	ResumeSubscription(ctx context.Context, subsNumber string) (*subscription.Subscription, error)

	// StopSubscription stop subscription
	StopSubscription(ctx context.Context, subsNumber string) (*subscription.Subscription, error)
}

// XenditProcessor callback handler for xendit
type XenditProcessor interface {
	ProcessDANACallback(ctx context.Context, dps *xendit.DANAPaymentStatus) error
	ProcessLinkAjaCallback(ctx context.Context, lps *xendit.LinkAjaPaymentStatus) error
	ProcessOVOCallback(ctx context.Context, ops *xendit.OVOPaymentStatus) error
	ProcessXenditInvoicesCallback(ctx context.Context, ips *xendit.InvoicePaymentStatus) error
	ProcessXenditEWalletCallback(ctx context.Context, status *xendit.EWalletPaymentStatus) error
}

// MidtransProcessor callback handler for midtrans
type MidtransProcessor interface {
	ProcessMidtransCallback(ctx context.Context, mr *coreapi.TransactionStatusResponse) error
}

// Payment combines all interface used for payment manager
type Payment interface {
	Interface
	XenditProcessor
	MidtransProcessor
}
