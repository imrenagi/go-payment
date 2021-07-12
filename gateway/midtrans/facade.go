package midtrans

import (
	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/gateway/midtrans/snap"
	"github.com/imrenagi/go-payment/invoice"

	"fmt"

	midsnap "github.com/midtrans/midtrans-go/snap"
)

// NewSnapFromInvoice create snap charge request
func NewSnapFromInvoice(inv *invoice.Invoice) (*midsnap.Request, error) {

	switch inv.Payment.PaymentType {
	case payment.SourceBCAVA:
		return snap.NewBCAVA(inv)
	case payment.SourcePermataVA:
		return snap.NewPermataVA(inv)
	case payment.SourceMandiriVA:
		return snap.NewMandiriVA(inv)
	case payment.SourceBNIVA:
		return snap.NewBNIVA(inv)
	case payment.SourceOtherVA:
		return snap.NewOtherBankVA(inv)
	case payment.SourceAlfamart:
		return snap.NewAlfamart(inv)
	case payment.SourceAkulaku:
		return snap.NewAkulaku(inv)
	case payment.SourceGopay:
		return snap.NewGopay(inv)
	case payment.SourceCreditCard:
		return snap.NewCreditCard(inv)
	case payment.SourceShopeePay:
		return snap.NewShopeePay(inv)
	default:
		return nil, fmt.Errorf("payment type not known")
	}
}
