package snap

import (
	"github.com/midtrans/midtrans-go/snap"

	"github.com/imrenagi/go-payment/invoice"
)

func NewAkulaku(inv *invoice.Invoice) (*snap.Request, error) {
	return newBuilder(inv).
		AddPaymentMethods(snap.PaymentTypeAkulaku).
		Build()
}
