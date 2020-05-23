package midtrans_test

import (
	"testing"

	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/invoice"

	"github.com/stretchr/testify/assert"

	. "github.com/imrenagi/go-payment/gateway/midtrans"
	gomidtrans "github.com/veritrans/go-midtrans"
)

func TestCreditCardWithoutInstallment(t *testing.T) {

	inv := dummyInv()
	inv.ServiceFee = 1000
	inv.Payment = &invoice.Payment{
		PaymentType: payment.SourceCreditCard,
	}
	builder := NewSnapRequestBuilder(inv)
	cc, _ := NewCreditCard(builder, inv.Payment.CreditCardDetail)

	req, err := cc.Build()
	if err != nil {
		t.Logf("expect no error, got %v", err)
		t.Fail()
	}

	assert.Len(t, req.EnabledPayments, 1)
	assert.Equal(t, int64(6000), req.TransactionDetails.GrossAmt)
	assert.Equal(t, 2, len(*req.Items))
	assert.Contains(t, req.EnabledPayments, gomidtrans.SourceCreditCard)

	assert.True(t, req.CreditCard.Secure)
	assert.Equal(t, "bca", req.CreditCard.Bank)
	assert.Equal(t, "3ds", req.CreditCard.Authentication)

}

func TestCreditCardWithInstallment(t *testing.T) {

	inv := dummyInv()
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
	builder := NewSnapRequestBuilder(inv)
	cc, _ := NewCreditCard(builder, inv.Payment.CreditCardDetail)

	req, err := cc.Build()
	if err != nil {
		t.Logf("expect no error, got %v", err)
		t.Fail()
	}

	assert.Len(t, req.EnabledPayments, 1)
	assert.Equal(t, int64(7000), req.TransactionDetails.GrossAmt)
	assert.Equal(t, 2, len(*req.Items))
	assert.Contains(t, req.EnabledPayments, gomidtrans.SourceCreditCard)
	assert.True(t, req.CreditCard.Secure)
	assert.Equal(t, "bca", req.CreditCard.Bank)
	assert.Equal(t, "3ds", req.CreditCard.Authentication)
	assert.True(t, req.CreditCard.Installment.Required)
	assert.Contains(t, req.CreditCard.Installment.Terms.Offline, int8(3))
	assert.Empty(t, req.CreditCard.Installment.Terms.Bni)
	assert.Empty(t, req.CreditCard.Installment.Terms.Mandiri)
}
