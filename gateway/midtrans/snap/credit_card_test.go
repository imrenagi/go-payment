package snap_test

import (
	"context"
	"testing"
	"time"

	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/gateway/midtrans/snap"
	"github.com/imrenagi/go-payment/invoice"
	midsnap "github.com/midtrans/midtrans-go/snap"
	"github.com/stretchr/testify/assert"
)

func baseCreditCardInvoice() *invoice.Invoice {
	date := time.Date(2020, 8, 1, 1, 0, 0, 0, time.UTC)
	dueDate := date.Add(24 * time.Hour)
	i := invoice.New(date, dueDate)

	i.SubTotal = 5000
	i.UpsertBillingAddress("Foo", "foo@bar.com", "0812312412")

	i.SetItems(context.TODO(),
		[]invoice.LineItem{
			{
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

func TestCreditCardWithoutInstallment(t *testing.T) {

	inv := baseCreditCardInvoice()
	inv.ServiceFee = 1000
	inv.Payment = &invoice.Payment{
		PaymentType: payment.SourceCreditCard,
	}
	req, err := snap.NewCreditCard(inv)
	assert.NoError(t, err)

	assert.Len(t, req.EnabledPayments, 1)
	assert.Equal(t, int64(6000), req.TransactionDetails.GrossAmt)
	assert.Equal(t, 2, len(*req.Items))
	assert.Contains(t, req.EnabledPayments, midsnap.PaymentTypeCreditCard)

	assert.True(t, req.CreditCard.Secure)
	assert.Equal(t, "bca", req.CreditCard.Bank)
}

func TestCreditCardWithInstallment(t *testing.T) {

	inv := baseCreditCardInvoice()
	inv.InstallmentFee = 2000
	inv.Payment = &invoice.Payment{
		PaymentType: payment.SourceCreditCard,
		CreditCardDetail: &invoice.CreditCardDetail{
			Installment: invoice.Installment{
				Type: payment.InstallmentOffline,
				Term: 3,
			},
			Bank: payment.BankBCA,
		},
	}

	req, _ := snap.NewCreditCard(inv)
	assert.Len(t, req.EnabledPayments, 1)
	assert.Equal(t, int64(7000), req.TransactionDetails.GrossAmt)
	assert.Equal(t, 2, len(*req.Items))
	assert.Contains(t, req.EnabledPayments, midsnap.PaymentTypeCreditCard)
	assert.True(t, req.CreditCard.Secure)
	assert.Equal(t, "bca", req.CreditCard.Bank)
	assert.True(t, req.CreditCard.Installment.Required)
	assert.Contains(t, req.CreditCard.Installment.Terms.Offline, int8(3))
	assert.Empty(t, req.CreditCard.Installment.Terms.Bni)
	assert.Empty(t, req.CreditCard.Installment.Terms.Mandiri)
}
