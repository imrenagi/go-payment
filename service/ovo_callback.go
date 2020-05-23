package service

import (
	"context"
	"errors"

	"github.com/imrenagi/go-payment"
)

type OVOPaymentStatus struct {
	Event       string  `json:"event"`
	ID          string  `json:"id"`
	ExternalID  string  `json:"external_id"`
	BusinessID  string  `json:"business_id"`
	Phone       string  `json:"phone"`
	EWalletType string  `json:"ewallet_type"`
	Amount      float64 `json:"amount"`
	FailureCode string  `json:"failure_code"`
	Status      string  `json:"status"`
}

func (p *Service) HandleOVOStatusCallback(ctx context.Context, ovo *OVOPaymentStatus) error {

	_, err := p.GetInvoice(ctx, ovo.ExternalID)
	if err != nil && !errors.Is(err, payment.ErrNotFound) {
		return err
	}

	// Need to return no error if callback
	if errors.Is(err, payment.ErrNotFound) {
		return nil
	}

	if ovo.Status == "FAILED" {
		if _, err := p.FailInvoice(ctx, ovo.ExternalID); err != nil {
			return err
		}
	}

	if ovo.Status == "COMPLETED" {
		if _, err := p.PayInvoice(ctx, ovo.ExternalID, PayInvoiceCommand{TransactionID: ovo.ID}); err != nil {
			return err
		}
	}

	return nil
}
