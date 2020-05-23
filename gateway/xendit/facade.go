package xendit

import (
	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/invoice"
	"github.com/xendit/xendit-go/ewallet"
	xinvoice "github.com/xendit/xendit-go/invoice"
)

func NewEwalletRequestFromInvoice(inv *invoice.Invoice) (*ewallet.CreatePaymentParams, error) {

	var reqBuilder ewalletRequestBuilder
	var err error

	req := NewEWalletRequest(inv)

	switch inv.Payment.PaymentType {
	case payment.SourceOvo:
		reqBuilder, err = NewOVO(req)
	case payment.SourceDana:
		reqBuilder, err = NewDana(req)
	case payment.SourceLinkAja:
		reqBuilder, err = NewLinkAja(req)
	}
	if err != nil {
		return nil, err
	}

	return reqBuilder.Build()
}

func NewInvoiceRequestFromInvoice(inv *invoice.Invoice) (*xinvoice.CreateParams, error) {

	var reqBuilder invoiceRequestBuilder
	var err error

	req := NewInvoiceRequestBuilder(inv)

	switch inv.Payment.PaymentType {
	case payment.SourceOvo:
		reqBuilder, err = NewOVOInvoice(req)
	}
	if err != nil {
		return nil, err
	}

	return reqBuilder.Build()
}
