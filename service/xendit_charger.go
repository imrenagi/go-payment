package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/imrenagi/go-payment"
	factory "github.com/imrenagi/go-payment/gateway/xendit"
	"github.com/imrenagi/go-payment/gateway/xendit"
	"github.com/imrenagi/go-payment/invoice"
	goxendit "github.com/xendit/xendit-go"
)

type xenditCharger struct {
	XenditGateway *xendit.Gateway
}

func (c xenditCharger) Create(ctx context.Context, inv *invoice.Invoice) (*invoice.ChargeResponse, error) {

	// recurringRequest, err := factory.NewRecurringChargeRequestBuilder(inv).Build()
	// if err != nil {
	// 	return nil, err
	// }

	// bytes, err := json.MarshalIndent(recurringRequest, "", "\t")
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Println(string(bytes))

	// xres, err := c.XenditGateway.Recurring.CreateWithContext(ctx, recurringRequest)
	// var xError *goxendit.Error
	// if ok := errors.As(err, &xError); ok && xError != nil {
	// 	return nil, xError
	// }

	// // TODO change this to active instead of pending.
	// if xres.Status == "PENDING" {
	// 	if err := inv.Process(ctx); err != nil {
	// 		return nil, err
	// 	}
	// }

	// return &invoice.ChargeResponse{
	// 	PaymentURL:    xres.LastCreatedInvoiceURL,
	// 	TransactionID: xres.ID,
	// }, nil

	switch inv.Payment.PaymentType {
	case payment.SourceLinkAja,
		payment.SourceDana:
		ewalletRequest, err := factory.NewEwalletRequestFromInvoice(inv)
		if err != nil {
			return nil, err
		}

		bytes, err := json.MarshalIndent(ewalletRequest, "", "\t")
		if err != nil {
			return nil, err
		}
		fmt.Println(string(bytes))

		xres, err := c.XenditGateway.Ewallet.CreatePayment(ewalletRequest)
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
			PaymentURL:    xres.CheckoutURL,
			TransactionID: xres.EWalletTransactionID,
		}, nil
	case payment.SourceOvo:
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

func (c xenditCharger) Gateway() string {
	// TODO change this later
	return "xendit"
}
