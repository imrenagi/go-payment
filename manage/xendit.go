package manage

import (
	"context"
	"errors"
	"strings"

	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/config"
	"github.com/imrenagi/go-payment/gateway/xendit"
	"github.com/imrenagi/go-payment/invoice"
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

	if ips.RecurringPaymentID != "" {
		return m.processXenditRecurringTransactionCallback(ctx, ips)
	}
	return m.processXenditNonRecurringTransactionCallback(ctx, ips)
}

func (m Manager) processXenditNonRecurringTransactionCallback(ctx context.Context, ips *xendit.InvoicePaymentStatus) error {
	inv, err := m.GetInvoice(ctx, ips.ExternalID)
	if err != nil && !errors.Is(err, payment.ErrNotFound) {
		return err
	}

	// Need to return no error if invoice is not found in this case
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

func (m Manager) processXenditRecurringTransactionCallback(ctx context.Context, ips *xendit.InvoicePaymentStatus) error {

	invoiceIDs := strings.Split(ips.ExternalID, "-")

	// external_id from xendit invoice for recurring transaction has format:
	// <subscription-number>_<timestamp>
	subscriptionNumber := strings.Join(invoiceIDs[:len(invoiceIDs)-1], "-")
	subs, err := m.subscriptionRepository.FindByNumber(ctx, subscriptionNumber)
	if err != nil {
		return err
	}

	inv, err := m.invoiceRepository.FindByNumber(ctx, ips.ExternalID)
	if err != nil && !errors.Is(err, payment.ErrNotFound) {
		return err
	}

	if errors.Is(err, payment.ErrNotFound) {
		fee := noAdminFee{payment.GatewayXendit}
		cfg, _ := fee.FindByPaymentType(ctx, "")
		payment, err := invoice.NewPayment(cfg, xendit.NewPaymentSource(ips.PaymentMethod), nil)

		inv = invoice.NewDefault()
		inv.Number = ips.ExternalID

		if err = inv.SetItems(ctx, []invoice.LineItem{
			*invoice.NewLineItem(
				"",
				ips.Description,
				ips.MerchantName,
				ips.Description,
				ips.Amount,
				1,
				ips.Currency,
			),
		}); err != nil {
			return err
		}

		if err = inv.UpsertBillingAddress(ips.PayerEmail, ips.PayerEmail, ""); err != nil {
			return err
		}

		if err = inv.UpdatePaymentMethod(ctx, payment, fee); err != nil {
			return err
		}

		if err = inv.Publish(ctx); err != nil {
			return err
		}
	}

	if ips.Status == "EXPIRED" {
		if err := inv.Fail(ctx); err != nil {
			return err
		}
	}

	if ips.Status == "PAID" {
		if err := inv.Pay(ctx, ips.ExternalID); err != nil {
			return err
		}
	}

	if err := subs.Record(inv); err != nil {
		return err
	}

	if err := m.subscriptionRepository.Save(ctx, subs); err != nil {
		return err
	}

	return nil
}

type noAdminFee struct {
	gateway payment.Gateway
}

func (n noAdminFee) FindByPaymentType(ctx context.Context, paymentType payment.PaymentType, opts ...payment.Option) (config.FeeConfigReader, error) {
	return config.NewFreeFee(n.gateway), nil
}
