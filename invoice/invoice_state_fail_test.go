package invoice_test

import (
	"context"
	"testing"

	"github.com/imrenagi/go-payment/config"
	cfgm "github.com/imrenagi/go-payment/config/mocks"
	dsm "github.com/imrenagi/go-payment/datastore/mocks"
	. "github.com/imrenagi/go-payment/invoice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFailedState_Pay(t *testing.T) {

	t.Run("should change to paid if the invoice is paid", func(t *testing.T) {
		i := emptyInvoice()
		i.SetState(&FailedState{})

		err := i.Pay(context.TODO(), "")
		assert.Nil(t, err)
		assert.Equal(t, Paid, i.State)
	})

	t.Run("can't change state to published", func(t *testing.T) {
		i := emptyInvoice()
		i.SetState(&FailedState{})

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

		i.UpdatePaymentMethod(context.TODO(), &Payment{}, readerMock)
		i.UpsertBillingAddress("foo", "foo@bar.com", "0811")

		err := i.Publish(context.TODO())
		assert.NotNil(t, err)
		assert.Equal(t, InvoiceError{InvoiceErrorInvalidStateTransition}, err)
	})

	t.Run("can't change state to process", func(t *testing.T) {
		i := emptyInvoice()
		i.SetState(&FailedState{})

		err := i.Process(context.TODO())
		assert.NotNil(t, err)
		assert.Equal(t, InvoiceError{InvoiceErrorInvalidStateTransition}, err)
	})

	t.Run("fail do nothing", func(t *testing.T) {
		i := emptyInvoice()
		i.SetState(&FailedState{})

		err := i.Fail(context.TODO())
		assert.Nil(t, err)
	})

}
