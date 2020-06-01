package midtrans_test

import (
	"context"
	"time"

	"github.com/imrenagi/go-payment/invoice"
)

func dummyInv() *invoice.Invoice {
	date := time.Date(2020, 8, 1, 1, 0, 0, 0, time.UTC)
	dueDate := date.Add(24 * time.Hour)
	i := invoice.New(date, dueDate)

	i.SubTotal = 5000
	i.UpsertBillingAddress("Foo", "foo@bar.com", "0812312412")

	i.SetItems(context.TODO(),
		[]invoice.LineItem{
			invoice.LineItem{
				InvoiceID:    1,
				Name:         "Terjemahan B",
				Category:     "TRANSLATION",
				MerchantName: "Collegos",
				Currency:     "IDR",
				UnitPrice:    5000,
				Qty:          1,
			}},
	)

	return i
}
