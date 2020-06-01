package invoice_test

import (
	"context"

	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/config"
	cfgm "github.com/imrenagi/go-payment/config/mocks"
	dsm "github.com/imrenagi/go-payment/datastore/mocks"
	"github.com/stretchr/testify/mock"

	. "github.com/imrenagi/go-payment/invoice"
)

func draftInvoice() *Invoice {
	i := emptyInvoice()
	i.UpsertBillingAddress("Foo", "foo@bar.com", "08123")

	feeMock := &cfgm.FeeConfigReader{}
	readerMock := &dsm.PaymentConfigReader{}

	readerMock.On("FindByPaymentType", mock.Anything, mock.Anything, mock.Anything).
		Return(feeMock, nil)

	feeMock.On("GetAdminFeeConfig", mock.Anything).Return(
		&config.Fee{
			PercentageVal: 0,
			CurrencyVal:   0,
			Currency:      "IDR",
		}, nil)
	feeMock.On("GetInstallmentFeeConfig", mock.Anything).Return(
		&config.Fee{
			PercentageVal: 0,
			CurrencyVal:   0,
			Currency:      "IDR",
		}, nil)

	i.UpdatePaymentMethod(context.TODO(), &Payment{
		PaymentType: payment.SourceBNIVA,
	}, readerMock)
	return i
}
