package service

import (
	"context"
	"errors"

	"github.com/imrenagi/go-payment"
)

type DANAPaymentStatus struct {
	ExternalID        string  `json:"external_id"`
	Amount            float64 `json:"amount"`
	BusinessID        string  `json:"business_id"`
	EWalletType       string  `json:"ewallet_type"`
	PaymentStatus     string  `json:"payment_status"`
	TransactionDate   string  `json:"transaction_date"`
	CallbackAuthToken string  `json:"callback_authentication_token"`
}

func (p *Service) HandleDanaStatusCallback(ctx context.Context, dana *DANAPaymentStatus) error {

	// TODO validate sender by checking the token

	_, err := p.GetInvoice(ctx, dana.ExternalID)
	if err != nil && !errors.Is(err, payment.ErrNotFound) {
		return err
	}

	if dana.PaymentStatus == "EXPIRED" {
		if _, err := p.FailInvoice(ctx, dana.ExternalID); err != nil {
			return err
		}
	}

	if dana.PaymentStatus == "PAID" {
		if _, err := p.PayInvoice(ctx, dana.ExternalID, PayInvoiceCommand{TransactionID: dana.ExternalID}); err != nil {
			return err
		}
	}

	return nil
}
