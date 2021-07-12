package xendit

import (
	"fmt"

	"github.com/xendit/xendit-go/ewallet"
	xinvoice "github.com/xendit/xendit-go/invoice"

	"github.com/imrenagi/go-payment"
	v1 "github.com/imrenagi/go-payment/gateway/xendit/ewallet/v1"
	v2 "github.com/imrenagi/go-payment/gateway/xendit/ewallet/v2"
	"github.com/imrenagi/go-payment/gateway/xendit/xeninvoice"
	"github.com/imrenagi/go-payment/invoice"
)

// NewEWalletChargeRequestFromInvoice create ewallet charge params for xendit ewallet API
func NewEWalletChargeRequestFromInvoice(inv *invoice.Invoice) (*ewallet.CreateEWalletChargeParams, error) {
	switch inv.Payment.PaymentType {
	case payment.SourceOvo:
		return v2.NewOVO(inv)
	case payment.SourceDana:
		return v2.NewDana(inv)
	case payment.SourceLinkAja:
		return v2.NewLinkAja(inv)
	default:
		return nil, fmt.Errorf("unsupported payment method")
	}
}

// Deprecated: NewEwalletRequestFromInvoice creates ewallet request for xendit
func NewEwalletRequestFromInvoice(inv *invoice.Invoice) (*ewallet.CreatePaymentParams, error) {
	switch inv.Payment.PaymentType {
	case payment.SourceOvo:
		return v1.NewOVO(inv)
	case payment.SourceDana:
		return v1.NewDana(inv)
	case payment.SourceLinkAja:
		return v1.NewLinkAja(inv)
	default:
		return nil, fmt.Errorf("payment type is not known")
	}
}

func NewInvoiceRequestFromInvoice(inv *invoice.Invoice) (*xinvoice.CreateParams, error) {
	switch inv.Payment.PaymentType {
	case payment.SourceOvo:
		return xeninvoice.NewOVO(inv)
	case payment.SourceDana:
		return xeninvoice.NewDana(inv)
	case payment.SourceLinkAja:
		return xeninvoice.NewLinkAja(inv)
	case payment.SourceAlfamart:
		return xeninvoice.NewAlfamart(inv)
	case payment.SourceBCAVA:
		return xeninvoice.NewBCAVA(inv)
	case payment.SourceBRIVA:
		return xeninvoice.NewBRIVA(inv)
	case payment.SourceBNIVA:
		return xeninvoice.NewBNIVA(inv)
	case payment.SourcePermataVA:
		return xeninvoice.NewPermataVA(inv)
	case payment.SourceMandiriVA:
		return xeninvoice.NewMandiriVA(inv)
	case payment.SourceCreditCard:
		return xeninvoice.NewCreditCard(inv)
	case payment.SourceShopeePay:
		return xeninvoice.NewShopeePay(inv)
	case payment.SourceQRIS:
		return xeninvoice.NewQRIS(inv)
	default:
		return nil, fmt.Errorf("payment type is not known")
	}
}
