package manage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/imrenagi/go-payment"
	midfactory "github.com/imrenagi/go-payment/gateway/midtrans"
	midgateway "github.com/imrenagi/go-payment/gateway/midtrans"
	"github.com/imrenagi/go-payment/gateway/xendit"
	factory "github.com/imrenagi/go-payment/gateway/xendit"
	"github.com/imrenagi/go-payment/invoice"
	"github.com/imrenagi/go-payment/util/localconfig"

	goxendit "github.com/xendit/xendit-go"
)

type midtransCharger struct {
	MidtransGateway *midgateway.Gateway
}

func (c midtransCharger) Create(ctx context.Context, inv *invoice.Invoice) (*invoice.ChargeResponse, error) {

	l := log.Ctx(ctx).With().
		Str("function", "midtransCharger.Create").
		Str("invoice_number", inv.Number).
		Logger()

	snapRequest, err := midfactory.NewSnapFromInvoice(inv)
	if err != nil {
		return nil, err
	}

	bytes, err := json.MarshalIndent(snapRequest, "", "\t")
	if err != nil {
		return nil, err
	}
	l.Debug().
		RawJSON("payload", bytes).
		Msg("snap request is created")

	resp, mErr := c.MidtransGateway.SnapV2Gateway.CreateTransaction(snapRequest)
	if mErr != nil {
		return nil, err
	}

	return &invoice.ChargeResponse{
		PaymentToken: resp.Token,
		PaymentURL:   resp.RedirectURL,
	}, nil
}

func (c midtransCharger) Gateway() payment.Gateway {
	return payment.GatewayMidtrans
}

type xenditCharger struct {
	config        localconfig.Xendit
	XenditGateway *xendit.Gateway
}

func (c xenditCharger) Create(ctx context.Context, inv *invoice.Invoice) (*invoice.ChargeResponse, error) {

	switch inv.Payment.PaymentType {
	case payment.SourceOvo:
		fn := c.eWalletChargeMethod(c.config.EWallet.OVO)
		return fn(ctx, inv)
	case payment.SourceLinkAja:
		fn := c.eWalletChargeMethod(c.config.EWallet.LinkAja)
		return fn(ctx, inv)
	case payment.SourceDana:
		fn := c.eWalletChargeMethod(c.config.EWallet.Dana)
		return fn(ctx, inv)
	default:
		return c.createXenInvoice(ctx, inv)
	}

	return nil, fmt.Errorf("payment method is not recognized")
}

type chargeFn func(ctx context.Context, inv *invoice.Invoice) (*invoice.ChargeResponse, error)

func (c xenditCharger) eWalletChargeMethod(cfg localconfig.EWalletConfig) chargeFn {
	if cfg.UseInvoice {
		return c.createXenInvoice
	} else if cfg.UseLegacy {
		return c.createLegacyEWalletCharge
	} else {
		return c.createEWalletCharge
	}
}

// Deprecated: createLegacyEWalletCharge ...
func (c xenditCharger) createLegacyEWalletCharge(ctx context.Context, inv *invoice.Invoice) (*invoice.ChargeResponse, error) {

	l := log.Ctx(ctx).With().
		Str("function", "xenditCharger.createLegacyEWalletCharge").
		Logger()

	l.Debug().Msg("generating ewallet request")

	ewalletRequest, err := factory.NewEwalletRequestFromInvoice(inv)
	if err != nil {
		return nil, err
	}

	bytes, err := json.MarshalIndent(ewalletRequest, "", "\t")
	if err != nil {
		return nil, err
	}

	l.Debug().
		RawJSON("payload", bytes).
		Msg("ewallet request is created")

	xres, err := c.XenditGateway.Ewallet.CreatePayment(ewalletRequest)
	if err != nil {
		var xError *goxendit.Error
		if ok := errors.As(err, &xError); ok && xError != nil {
			l.Error().Err(xError).Msg("unable to create payment")
			return nil, xError
		}
	}

	if xres.Status == "PENDING" {
		l.Info().Msg("set invoice to pending")
		if err := inv.Process(ctx); err != nil {
			l.Error().Err(err).Msg("unable to set invoice to pending")
			return nil, err
		}
	}

	l.Info().Msg("ewallet request is created")

	return &invoice.ChargeResponse{
		PaymentURL:    xres.CheckoutURL,
		TransactionID: xres.EWalletTransactionID,
	}, nil
}

func (c xenditCharger) createEWalletCharge(ctx context.Context, inv *invoice.Invoice) (*invoice.ChargeResponse, error) {

	l := log.Ctx(ctx).With().
		Str("function", "xenditCharger.createEWalletCharge").
		Logger()

	l.Debug().Msg("generating ewallet request")

	ewChargeParams, err := factory.NewEWalletChargeRequestFromInvoice(inv)
	if err != nil {
		return nil, err
	}

	bytes, err := json.MarshalIndent(ewChargeParams, "", "\t")
	if err != nil {
		return nil, err
	}

	l.Debug().
		RawJSON("payload", bytes).
		Msg("ewallet request is created")

	chargeRes, err := c.XenditGateway.Ewallet.CreateEWalletChargeWithContext(ctx, ewChargeParams)
	if err != nil {
		var xError *goxendit.Error
		if ok := errors.As(err, &xError); ok && xError != nil {
			l.Error().Err(xError).Msg("unable to create payment")
			return nil, xError
		}
	}

	if chargeRes.Status == "PENDING" {
		l.Info().Msg("set invoice to pending")
		if err := inv.Process(ctx); err != nil {
			l.Error().Err(err).Msg("unable to set invoice to pending")
			return nil, err
		}
	}

	// Need to take care of other URLs later
	var paymentURL string
	if url, ok := chargeRes.Actions["desktop_web_checkout_url"]; ok {
		paymentURL = url
	}

	l.Info().Msg("ewallet request is created")

	return &invoice.ChargeResponse{
		PaymentURL:    paymentURL, //TODO(imre) handle mobile url etc
		TransactionID: chargeRes.ID,
	}, nil
}

func (c xenditCharger) createXenInvoice(ctx context.Context, inv *invoice.Invoice) (*invoice.ChargeResponse, error) {

	l := log.Ctx(ctx).With().
		Str("function", "xenditCharger.createXenInvoice").
		Logger()

	l.Debug().Msg("generating xeninvoice request")

	invoiceRequest, err := factory.NewInvoiceRequestFromInvoice(inv)
	if err != nil {
		return nil, err
	}

	bytes, err := json.MarshalIndent(invoiceRequest, "", "\t")
	if err != nil {
		return nil, err
	}
	l.Debug().
		RawJSON("payload", bytes).
		Msg("xeninvoice request is created")

	xres, err := c.XenditGateway.Invoice.CreateWithContext(ctx, invoiceRequest)
	if err != nil {
		var xError *goxendit.Error
		if ok := errors.As(err, &xError); ok && xError != nil {
			l.Error().Err(xError).Msg("unable to create payment")
			return nil, xError
		}
	}

	if xres.Status == "PENDING" {
		l.Info().Msg("set invoice to pending")
		if err := inv.Process(ctx); err != nil {
			l.Error().Err(err).Msg("unable to set invoice to pending")
			return nil, err
		}
	}

	l.Info().Msg("xeninvoice request is created")

	return &invoice.ChargeResponse{
		PaymentURL:    xres.InvoiceURL,
		TransactionID: xres.ID,
	}, nil
}

func (c xenditCharger) Gateway() payment.Gateway {
	return payment.GatewayXendit
}
