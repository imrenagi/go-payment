package snap

import (
	"github.com/midtrans/midtrans-go/snap"

	"github.com/imrenagi/go-payment/invoice"
)

func NewAlfamart(inv *invoice.Invoice) (*snap.Request, error) {
	return newBuilder(inv).
		AddPaymentMethods(snap.PaymentTypeAlfamart).
		Build()
}
