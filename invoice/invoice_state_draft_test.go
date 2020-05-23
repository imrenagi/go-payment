package invoice_test

import (
	"github.com/imrenagi/go-payment"
	. "github.com/imrenagi/go-payment/invoice"
)

func draftInvoice() *Invoice {
	i := emptyInvoice()
	i.SetBillingAddress("Foo", "foo@bar.com", "08123")
	i.SetPaymentMethod(&Payment{
		PaymentType: payment.SourceBNIVA,
	})
	return i
}
