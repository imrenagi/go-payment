package midtrans

import (
	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/invoice"
	gomidtrans "github.com/veritrans/go-midtrans"
)

func NewSnapRequestFromInvoice(inv *invoice.Invoice) (*gomidtrans.SnapReq, error) {

	var reqBuilder requestBuilder
	var err error

	snapRequestBuilder := NewSnapRequestBuilder(inv)

	switch inv.Payment.PaymentType {
	case payment.SourceBCAVA:
		reqBuilder, err = NewBCAVA(snapRequestBuilder)
	case payment.SourcePermataVA:
		reqBuilder, err = NewPermataVA(snapRequestBuilder)
	case payment.SourceEchannel:
		reqBuilder, err = NewMandiriBill(snapRequestBuilder)
	case payment.SourceBNIVA:
		reqBuilder, err = NewBNIVA(snapRequestBuilder)
	case payment.SourceOtherVA:
		reqBuilder, err = NewOtherBank(snapRequestBuilder)
	case payment.SourceGopay:
		reqBuilder, err = NewGopay(snapRequestBuilder)
	case payment.SourceAlfamart:
		reqBuilder, err = NewAlfamart(snapRequestBuilder)
	case payment.SourceAkulaku:
		reqBuilder, err = NewAkulaku(snapRequestBuilder)
	case payment.SourceCreditCard:
		reqBuilder, err = NewCreditCard(snapRequestBuilder, inv.Payment.CreditCardDetail)
	}
	if err != nil {
		return nil, err
	}

	return reqBuilder.Build()
}
