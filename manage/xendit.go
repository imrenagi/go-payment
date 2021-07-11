package manage

import (
	"context"
	"errors"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/config"
	"github.com/imrenagi/go-payment/gateway/xendit"
	"github.com/imrenagi/go-payment/invoice"
)

// ProcessDANACallback process dana payment notification from xendit
func (m *Manager) ProcessDANACallback(ctx context.Context, dps *xendit.DANAPaymentStatus) error {

	l := log.Ctx(ctx).With().
		Str("function", "Manager.ProcessDANACallback").
		Str("invoice_number", dps.ExternalID).
		Str("payment_status", dps.PaymentStatus).
		Logger()

	if err := dps.IsValid(m.xenditGateway.NotificationValidationKey()); err != nil {
		l.Warn().Err(err).Msg("invalid notification callback key")
		return err
	}

	inv, err := m.GetInvoice(ctx, dps.ExternalID)
	if err != nil && !errors.Is(err, payment.ErrNotFound) {
		l.Error().Err(err).Msg("unable to get the invoice")
		return err
	}

	if dps.PaymentStatus == "EXPIRED" {
		l.Info().Msg("set invoice to failed")
		if _, err := m.FailInvoice(ctx, &FailInvoiceRequest{
			InvoiceNumber: inv.Number,
			TransactionID: dps.ExternalID,
			Reason:        dps.PaymentStatus,
		}); err != nil {
			l.Error().Err(err).Msg("unable to fail the invoice")
			return err
		}
	}

	if dps.PaymentStatus == "PAID" {
		l.Info().Msg("set invoice to paid")
		if _, err := m.PayInvoice(ctx, &PayInvoiceRequest{
			InvoiceNumber: inv.Number,
			TransactionID: dps.ExternalID,
		}); err != nil {
			l.Error().Err(err).Msg("unable to set invoice to paid")
			return err
		}
	}

	l.Info().Msg("callback is processed successfully")
	return nil
}

// ProcessLinkAjaCallback process linkaja payment notification from xendit
func (m *Manager) ProcessLinkAjaCallback(ctx context.Context, lps *xendit.LinkAjaPaymentStatus) error {

	l := log.Ctx(ctx).With().
		Str("function", "Manager.ProcessLinkAjaCallback").
		Str("invoice_number", lps.ExternalID).
		Str("payment_status", lps.Status).
		Logger()

	if err := lps.IsValid(m.xenditGateway.NotificationValidationKey()); err != nil {
		l.Warn().Err(err).Msg("invalid notification callback key")
		return err
	}

	inv, err := m.GetInvoice(ctx, lps.ExternalID)
	if err != nil && !errors.Is(err, payment.ErrNotFound) {
		l.Error().Err(err).Msg("unable to get the invoice")
		return err
	}

	if lps.Status == "FAILED" {
		l.Info().Msg("set invoice to failed")
		if _, err := m.FailInvoice(ctx, &FailInvoiceRequest{
			InvoiceNumber: inv.Number,
			TransactionID: lps.ExternalID,
			Reason:        lps.Status,
		}); err != nil {
			l.Error().Err(err).Msg("unable to fail the invoice")
			return err
		}
	}

	if lps.Status == "SUCCESS_COMPLETED" {
		l.Info().Msg("set invoice to paid")
		if _, err := m.PayInvoice(ctx, &PayInvoiceRequest{
			InvoiceNumber: inv.Number,
			TransactionID: lps.ExternalID,
		}); err != nil {
			l.Error().Err(err).Msg("unable to set invoice to paid")
			return err
		}
	}

	l.Info().Msg("callback is processed successfully")
	return nil
}

// ProcessOVOCallback process ovo payment notification from xendit
func (m *Manager) ProcessOVOCallback(ctx context.Context, ops *xendit.OVOPaymentStatus) error {

	l := log.Ctx(ctx).With().
		Str("function", "Manager.ProcessOVOCallback").
		Str("invoice_number", ops.ExternalID).
		Str("payment_status", ops.Status).
		Logger()

	if err := ops.IsValid(m.xenditGateway.NotificationValidationKey()); err != nil {
		l.Warn().Err(err).Msg("invalid notification callback key")
		return err
	}

	inv, err := m.GetInvoice(ctx, ops.ExternalID)
	if err != nil && !errors.Is(err, payment.ErrNotFound) {
		l.Error().Err(err).Msg("unable to get the invoice")
		return err
	}

	if errors.Is(err, payment.ErrNotFound) {
		return nil
	}

	if ops.Status == "FAILED" {
		l.Info().Msg("set invoice to failed")
		if _, err := m.FailInvoice(ctx, &FailInvoiceRequest{
			InvoiceNumber: inv.Number,
			TransactionID: ops.ExternalID,
			Reason:        ops.FailureCode,
		}); err != nil {
			l.Error().Err(err).Msg("unable to fail the invoice")
			return err
		}
	}

	if ops.Status == "COMPLETED" {
		l.Info().Msg("set invoice to paid")
		if _, err := m.PayInvoice(ctx, &PayInvoiceRequest{
			InvoiceNumber: inv.Number,
			TransactionID: ops.ExternalID,
		}); err != nil {
			l.Error().Err(err).Msg("unable to set invoice to paid")
			return err
		}
	}

	l.Info().Msg("callback is processed successfully")
	return nil
}

// ProcessXenditInvoicesCallback process xendit invoice payment notification from xendit
func (m *Manager) ProcessXenditInvoicesCallback(ctx context.Context, ips *xendit.InvoicePaymentStatus) error {

	l := log.Ctx(ctx).With().
		Str("function", "Manager.ProcessXenditInvoicesCallback").
		Str("invoice_number", ips.ExternalID).
		Str("payment_status", ips.Status).
		Logger()

	if err := ips.IsValid(m.xenditGateway.NotificationValidationKey()); err != nil {
		l.Warn().Err(err).Msg("invalid notification callback key")
		return err
	}

	if ips.RecurringPaymentID != "" {
		return m.processXenditRecurringTransactionCallback(ctx, ips)
	}
	return m.processXenditNonRecurringTransactionCallback(ctx, ips)
}

func (m *Manager) ProcessXenditEWalletCallback(ctx context.Context, status *xendit.EWalletPaymentStatus) error {

	l := log.Ctx(ctx).With().
		Str("function", "Manager.ProcessXenditEWalletCallback").
		Str("invoice_number", status.Data.ReferenceID).
		Str("payment_status", status.Data.Status).
		Logger()

	if err := status.IsValid(m.xenditGateway.NotificationValidationKey()); err != nil {
		l.Warn().Err(err).Msg("invalid notification callback key")
		return err
	}

	inv, err := m.GetInvoice(ctx, status.Data.ReferenceID)
	if err != nil && !errors.Is(err, payment.ErrNotFound) {
		l.Error().Err(err).Msg("unable to get the invoice")
		return err
	}

	// Need to return no error if invoice is not found in this case
	if errors.Is(err, payment.ErrNotFound) {
		return nil
	}

	if status.Data.Status == "FAILED" || status.Data.Status == "VOIDED" {
		l.Info().Msg("set invoice to failed")
		if _, err := m.FailInvoice(ctx, &FailInvoiceRequest{
			InvoiceNumber: inv.Number,
			TransactionID: status.Data.ReferenceID,
			Reason:        status.Data.Status,
		}); err != nil {
			l.Error().Err(err).Msg("unable to fail the invoice")
			return err
		}
	}

	if status.Data.Status == "SUCCEEDED" {
		l.Info().Msg("set invoice to paid")
		if _, err := m.PayInvoice(ctx, &PayInvoiceRequest{
			InvoiceNumber: inv.Number,
			TransactionID: status.Data.ReferenceID,
		}); err != nil {
			l.Error().Err(err).Msg("unable to complete the invoice payment")
			return err
		}
	}

	l.Info().Msg("callback is processed successfully")
	return nil
}

func (m Manager) processXenditNonRecurringTransactionCallback(ctx context.Context, ips *xendit.InvoicePaymentStatus) error {

	l := log.Ctx(ctx).With().
		Str("function", "Manager.processXenditNonRecurringTransactionCallback").
		Str("invoice_number", ips.ExternalID).
		Str("payment_status", ips.Status).
		Logger()

	inv, err := m.GetInvoice(ctx, ips.ExternalID)
	if err != nil && !errors.Is(err, payment.ErrNotFound) {
		l.Error().Err(err).Msg("unable to get the invoice")
		return err
	}

	// Need to return no error if invoice is not found in this case
	if errors.Is(err, payment.ErrNotFound) {
		return nil
	}

	if ips.Status == "EXPIRED" {
		l.Info().Msg("set invoice to fail")
		if _, err := m.FailInvoice(ctx, &FailInvoiceRequest{
			InvoiceNumber: inv.Number,
			TransactionID: ips.ExternalID,
			Reason:        ips.Status,
		}); err != nil {
			l.Error().Err(err).Msg("unable to set invoice to fail")
			return err
		}
	}

	if ips.Status == "PAID" {
		l.Info().Msg("set invoice to success")
		if _, err := m.PayInvoice(ctx, &PayInvoiceRequest{
			InvoiceNumber: inv.Number,
			TransactionID: ips.ExternalID,
		}); err != nil {
			l.Error().Err(err).Msg("unable to set invoice to paid")
			return err
		}
	}

	l.Info().Msg("callback is successfully processed")
	return nil
}

func (m Manager) processXenditRecurringTransactionCallback(ctx context.Context, ips *xendit.InvoicePaymentStatus) error {

	l := log.Ctx(ctx).With().
		Str("function", "Manager.processXenditRecurringTransactionCallback").
		Str("external_id", ips.ExternalID).
		Str("payment_status", ips.Status).
		Logger()

	invoiceIDs := strings.Split(ips.ExternalID, "-")

	// external_id from xendit invoice for recurring transaction has format:
	// <subscription-number>_<timestamp>
	subscriptionNumber := strings.Join(invoiceIDs[:len(invoiceIDs)-1], "-")
	subs, err := m.subscriptionRepository.FindByNumber(ctx, subscriptionNumber)
	if err != nil {
		return err
	}

	l = l.With().Str("subscription_number", subscriptionNumber).Logger()

	inv, err := m.invoiceRepository.FindByNumber(ctx, ips.ExternalID)
	if err != nil && !errors.Is(err, payment.ErrNotFound) {
		l.Error().Err(err).Msg("unable to get the invoice")
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
		l.Info().Msg("set invoice to fail")
		if err := inv.Fail(ctx); err != nil {
			l.Error().Err(err).Msg("unable to set invoice to fail")
			return err
		}
	}

	if ips.Status == "PAID" {
		l.Info().Msg("set invoice to paid")
		if err := inv.Pay(ctx, ips.ExternalID); err != nil {
			l.Error().Err(err).Msg("unable to set invoice to paid")
			return err
		}
	}

	l.Debug().Msg("saving subscription")

	if err := subs.Save(inv); err != nil {
		return err
	}

	if err := m.subscriptionRepository.Save(ctx, subs); err != nil {
		return err
	}

	l.Info().Msg("subscription callback is done")

	return nil
}

type noAdminFee struct {
	gateway payment.Gateway
}

func (n noAdminFee) FindByPaymentType(ctx context.Context, paymentType payment.PaymentType, opts ...payment.Option) (config.FeeConfigReader, error) {
	return config.NewFreeFee(n.gateway), nil
}
