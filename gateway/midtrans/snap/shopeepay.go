package snap

import (
	midsnap "github.com/midtrans/midtrans-go/snap"

	"github.com/imrenagi/go-payment/invoice"
)

func NewShopeePay(inv *invoice.Invoice) (*midsnap.Request, error) {
	return newBuilder(inv).
		AddPaymentMethods(midsnap.PaymentTypeShopeepay).
		Build()
}
