package snap

import (
	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/invoice"

	"fmt"

	midsnap "github.com/midtrans/midtrans-go/snap"
)

// NewSnapRequestFromInvoice create snap charge request
func NewSnapRequestFromInvoice(inv *invoice.Invoice) (*midsnap.Request, error) {

	switch inv.Payment.PaymentType {
	// case payment.SourceBCAVA:
	//   reqBuilder, err = NewBCAVA(snapRequestBuilder)
	// case payment.SourcePermataVA:
	//   reqBuilder, err = NewPermataVA(snapRequestBuilder)
	// case payment.SourceMandiriVA:
	//   reqBuilder, err = NewMandiriBill(snapRequestBuilder)
	// case payment.SourceBNIVA:
	//   reqBuilder, err = NewBNIVA(snapRequestBuilder)
	// case payment.SourceOtherVA:
	//   reqBuilder, err = NewOtherBank(snapRequestBuilder)
	case payment.SourceGopay:
		return NewGopay(inv)
	// case payment.SourceAlfamart:
	//   reqBuilder, err = NewAlfamart(snapRequestBuilder)
	// case payment.SourceAkulaku:
	//   reqBuilder, err = NewAkulaku(snapRequestBuilder)
	// case payment.SourceCreditCard:
	//   reqBuilder, err = NewCreditCard(snapRequestBuilder, inv.Payment.CreditCardDetail)
	default:
		return nil, fmt.Errorf("payment type not known")
	}
}
