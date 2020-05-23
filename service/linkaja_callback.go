package service

import (
	"context"
	"errors"

	"github.com/imrenagi/go-payment"
)

type LinkAjaPaymentStatus struct {
	ExternalID        string  `json:"external_id"`
	Amount            float64 `json:"amount"`
	Status            string  `json:"status"`
	EWalletType       string  `json:"ewallet_type"`
	CallbackAuthToken string  `json:"callback_authentication_token"`
}

func (p *Service) HandleLinkAjaStatusCallback(ctx context.Context, linkaja *LinkAjaPaymentStatus) error {

	// TODO verify the sender by checking the token

	_, err := p.GetInvoice(ctx, linkaja.ExternalID)
	if err != nil && !errors.Is(err, payment.ErrNotFound) {
		return err
	}

	if linkaja.Status == "FAILED" {
		if _, err := p.FailInvoice(ctx, linkaja.ExternalID); err != nil {
			return err
		}
	}

	if linkaja.Status == "SUCCESS_COMPLETED" {
		if _, err := p.PayInvoice(ctx, linkaja.ExternalID, PayInvoiceCommand{TransactionID: linkaja.ExternalID}); err != nil {
			return err
		}
	}

	return nil
}
