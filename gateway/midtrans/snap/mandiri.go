package snap

import (
	"github.com/midtrans/midtrans-go/snap"

	"github.com/imrenagi/go-payment/invoice"
)

func NewMandiriVA(inv *invoice.Invoice) (*snap.Request, error) {
	return newBuilder(inv).
		AddPaymentMethods(snap.PaymentTypeEChannel).
		Build()
}
