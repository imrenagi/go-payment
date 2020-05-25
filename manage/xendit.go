package manage

import (
	"context"
	"errors"

	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/gateway/xendit"
)

// ProcessDANACallback process dana payment notification from xendit
func (m *Manager) ProcessDANACallback(ctx context.Context, dps *xendit.DANAPaymentStatus) error {

	if err := dps.IsValid(m.xenditGateway.NotificationValidationKey()); err != nil {
		return err
	}

	inv, err := m.GetInvoice(ctx, dps.ExternalID)
	if err != nil && !errors.Is(err, payment.ErrNotFound) {
		return err
	}

	if dps.PaymentStatus == "EXPIRED" {
		if _, err := m.FailInvoice(ctx, &FailInvoiceRequest{
			InvoiceNumber: inv.Number,
			TransactionID: dps.ExternalID,
			Reason:        dps.PaymentStatus,
		}); err != nil {
			return err
		}
	}

	if dps.PaymentStatus == "PAID" {
		if _, err := m.PayInvoice(ctx, &PayInvoiceRequest{
			InvoiceNumber: inv.Number,
			TransactionID: dps.ExternalID,
		}); err != nil {
			return err
		}
	}

	return nil
}

// ProcessLinkAjaCallback process linkaja payment notification from xendit
func (m *Manager) ProcessLinkAjaCallback(ctx context.Context, lps *xendit.LinkAjaPaymentStatus) error {

	if err := lps.IsValid(m.xenditGateway.NotificationValidationKey()); err != nil {
		return err
	}

	inv, err := m.GetInvoice(ctx, lps.ExternalID)
	if err != nil && !errors.Is(err, payment.ErrNotFound) {
		return err
	}

	if lps.Status == "FAILED" {
		if _, err := m.FailInvoice(ctx, &FailInvoiceRequest{
			InvoiceNumber: inv.Number,
			TransactionID: lps.ExternalID,
			Reason:        lps.Status,
		}); err != nil {
			return err
		}
	}

	if lps.Status == "SUCCESS_COMPLETED" {
		if _, err := m.PayInvoice(ctx, &PayInvoiceRequest{
			InvoiceNumber: inv.Number,
			TransactionID: lps.ExternalID,
		}); err != nil {
			return err
		}
	}

	return nil
}

// ProcessOVOCallback process ovo payment notification from xendit
func (m *Manager) ProcessOVOCallback(ctx context.Context, ops *xendit.OVOPaymentStatus) error {

	if err := ops.IsValid(m.xenditGateway.NotificationValidationKey()); err != nil {
		return err
	}

	inv, err := m.GetInvoice(ctx, ops.ExternalID)
	if err != nil && !errors.Is(err, payment.ErrNotFound) {
		return err
	}

	if errors.Is(err, payment.ErrNotFound) {
		return nil
	}

	if ops.Status == "FAILED" {
		if _, err := m.FailInvoice(ctx, &FailInvoiceRequest{
			InvoiceNumber: inv.Number,
			TransactionID: ops.ExternalID,
			Reason:        ops.FailureCode,
		}); err != nil {
			return err
		}
	}

	if ops.Status == "COMPLETED" {
		if _, err := m.PayInvoice(ctx, &PayInvoiceRequest{
			InvoiceNumber: inv.Number,
			TransactionID: ops.ExternalID,
		}); err != nil {
			return err
		}
	}

	return nil
}

// ProcessXenditInvoicesCallback process xendit invoice payment notification from xendit
func (m *Manager) ProcessXenditInvoicesCallback(ctx context.Context, ips *xendit.InvoicePaymentStatus) error {

	if err := ips.IsValid(m.xenditGateway.NotificationValidationKey()); err != nil {
		return err
	}

	inv, err := m.GetInvoice(ctx, ips.ExternalID)
	if err != nil && !errors.Is(err, payment.ErrNotFound) {
		return err
	}

	// Need to return no error if callback
	if errors.Is(err, payment.ErrNotFound) {
		return nil
	}

	if ips.Status == "EXPIRED" {
		if _, err := m.FailInvoice(ctx, &FailInvoiceRequest{
			InvoiceNumber: inv.Number,
			TransactionID: ips.ExternalID,
			Reason:        ips.Status,
		}); err != nil {
			return err
		}
	}

	if ips.Status == "PAID" {
		if _, err := m.PayInvoice(ctx, &PayInvoiceRequest{
			InvoiceNumber: inv.Number,
			TransactionID: ips.ExternalID,
		}); err != nil {
			return err
		}
	}

	return nil
}
