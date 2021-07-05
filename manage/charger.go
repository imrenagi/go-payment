package manage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/imrenagi/go-payment"
	midfactory "github.com/imrenagi/go-payment/gateway/midtrans"
	midgateway "github.com/imrenagi/go-payment/gateway/midtrans"
	"github.com/imrenagi/go-payment/gateway/xendit"
	factory "github.com/imrenagi/go-payment/gateway/xendit"
	"github.com/imrenagi/go-payment/invoice"
	goxendit "github.com/xendit/xendit-go"
)

type midtransCharger struct {
	MidtransGateway *midgateway.Gateway
}

func (c midtransCharger) Create(ctx context.Context, inv *invoice.Invoice) (*invoice.ChargeResponse, error) {

	snapRequest, err := midfactory.NewSnapRequestFromInvoice(inv)
	if err != nil {
		return nil, err
	}

	// bytes, err := json.MarshalIndent(snapRequest, "", "\t")
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Println(string(bytes))

	resp, err := c.MidtransGateway.SnapGateway.GetToken(snapRequest)
	if err != nil {
		return nil, err
	}

	return &invoice.ChargeResponse{
		PaymentToken: resp.Token,
		PaymentURL:   resp.RedirectURL,
		// TransactionID: resp.

	}, nil
}

func (c midtransCharger) Gateway() payment.Gateway {
	return payment.GatewayMidtrans
}

type xenditCharger struct {
	XenditGateway *xendit.Gateway
}

func (c xenditCharger) Create(ctx context.Context, inv *invoice.Invoice) (*invoice.ChargeResponse, error) {

	switch inv.Payment.PaymentType {
	case payment.SourceOvo:
		ewChargeParams, err := factory.NewEWalletChargeRequestFromInvoice(inv)
		if err != nil {
			return nil, err
		}

			bytes, err := json.MarshalIndent(ewChargeParams, "", "\t")
			if err != nil {
				return nil, err
			}
			fmt.Println(string(bytes))

		chargeRes, err := c.XenditGateway.Ewallet.CreateEWalletChargeWithContext(ctx, ewChargeParams)
		var xError *goxendit.Error
		if ok := errors.As(err, &xError); ok && xError != nil {
			return nil, xError
		}

		if chargeRes.Status == "PENDING" {
			if err := inv.Process(ctx); err != nil {
				return nil, err
			}
		}

		// Need to take care of other URLs later
		var paymentURL string
		if url, ok := chargeRes.Actions["desktop_web_checkout_url"]; ok {
			paymentURL = url
		}

		return &invoice.ChargeResponse{
			PaymentURL:    paymentURL, // TODO handle mobile url etc
			TransactionID: chargeRes.ID,
		}, nil
	// case payment.SourceLinkAja,
	// 	payment.SourceDana:
	// 	ewalletRequest, err := factory.NewEwalletRequestFromInvoice(inv)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	bytes, err := json.MarshalIndent(ewalletRequest, "", "\t")
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	fmt.Println(string(bytes))

	// 	xres, err := c.XenditGateway.Ewallet.CreatePayment(ewalletRequest)
	// 	var xError *goxendit.Error
	// 	if ok := errors.As(err, &xError); ok && xError != nil {
	// 		return nil, xError
	// 	}

	// 	if xres.Status == "PENDING" {
	// 		if err := inv.Process(ctx); err != nil {
	// 			return nil, err
	// 		}
	// 	}

	// 	return &invoice.ChargeResponse{
	// 		PaymentURL:    xres.CheckoutURL,
	// 		TransactionID: xres.EWalletTransactionID,
	// 	}, nil
	default:
		invoiceRequest, err := factory.NewInvoiceRequestFromInvoice(inv)
		if err != nil {
			return nil, err
		}

		bytes, err := json.MarshalIndent(invoiceRequest, "", "\t")
		if err != nil {
			return nil, err
		}
		fmt.Println(string(bytes))

		xres, err := c.XenditGateway.Invoice.CreateWithContext(ctx, invoiceRequest)
		var xError *goxendit.Error
		if ok := errors.As(err, &xError); ok && xError != nil {
			return nil, xError
		}

		if xres.Status == "PENDING" {
			if err := inv.Process(ctx); err != nil {
				return nil, err
			}
		}

		return &invoice.ChargeResponse{
			PaymentURL:    xres.InvoiceURL,
			TransactionID: xres.ID,
		}, nil
	}

	return nil, fmt.Errorf("payment method is not recognized")
}

func (c xenditCharger) Gateway() payment.Gateway {
	return payment.GatewayXendit
}
