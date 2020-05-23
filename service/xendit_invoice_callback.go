package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/imrenagi/go-payment"
)

// XenditInvoiceCallbackRequest
// https://xendit.github.io/apireference/#invoice-callback
type XenditInvoiceCallbackRequest struct {
	ID                     string  `json:"id"`
	ExternalID             string  `json:"external_id"`
	UserID                 string  `json:"user_id"`
	PaymentMethod          string  `json:"payment_method"`
	Status                 string  `json:"status"`
	MerchantName           string  `json:"merchant_name"`
	Amount                 float64 `json:"amount"`
	PaidAmount             float64 `json:"paid_amount"`
	BankCode               string  `json:"bank_code"`
	RetailOutletName       string  `json:"retail_outlet_name"`
	EwalletType            string  `json:"ewallet_type"`
	OnDemandLink           string  `json:"on_demand_link"`
	RecurringPaymentID     string  `json:"recurring_payment_id"`
	PaidAt                 string  `json:"paid_at"`
	PayerEmail             string  `json:"payer_email"`
	Description            string  `json:"description"`
	AdjustedReceivedAmount float64 `json:"adjusted_received_amount"`
	FeesPaidAmount         float64 `json:"fees_paid_amount"`
	CreatedAt              string  `json:"created"`
	UpdatedAt              string  `json:"updated"`
	Currency               string  `json:"currency"`
	PaymentChannel         string  `json:"payment_channel"`
	PaymentDestination     string  `json:"payment_destination"`
}

func (p *Service) HandleXenditInvoicesCallback(ctx context.Context, cbReq *XenditInvoiceCallbackRequest) error {

	fmt.Println(cbReq)
	fmt.Println(cbReq.RecurringPaymentID)

	_, err := p.GetInvoice(ctx, cbReq.ExternalID)
	if err != nil && !errors.Is(err, payment.ErrNotFound) {
		return err
	}

	// Need to return no error if callback
	if errors.Is(err, payment.ErrNotFound) {
		return nil
	}

	if cbReq.Status == "EXPIRED" {
		if _, err := p.FailInvoice(ctx, cbReq.ExternalID); err != nil {
			return err
		}
	}

	if cbReq.Status == "PAID" {
		if _, err := p.PayInvoice(ctx, cbReq.ExternalID, PayInvoiceCommand{TransactionID: cbReq.ID}); err != nil {
			return err
		}
	}

	return nil
}
