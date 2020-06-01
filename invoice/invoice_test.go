package invoice_test

import (
	"context"
	"testing"
	"time"

	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/config"
	cfgm "github.com/imrenagi/go-payment/config/mocks"
	dsm "github.com/imrenagi/go-payment/datastore/mocks"
	. "github.com/imrenagi/go-payment/invoice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func emptyInvoice() *Invoice {

	date := time.Date(2020, 8, 1, 1, 0, 0, 0, time.UTC)
	dueDate := date.Add(24 * time.Hour)
	i := New(date, dueDate)
	i.ID = uint64(1)

	return i
}

func TestNewInvoice(t *testing.T) {

	i := emptyInvoice()

	assert.NotEmpty(t, i.Number)
	assert.Equal(t, "IDR", i.Currency)
	assert.Equal(t, float64(0), i.SubTotal)
	assert.Equal(t, float64(0), i.Tax)
	assert.Equal(t, float64(0), i.Discount)
	assert.Equal(t, float64(0), i.ServiceFee)
	assert.Equal(t, float64(0), i.InstallmentFee)
	assert.Equal(t, Draft, i.State)
	assert.NotNil(t, i.StateController)
	assert.Nil(t, i.LineItems)
	assert.Empty(t, i.BillingAddress)

}

func TestTitle(t *testing.T) {

	t.Run("no installment information shown up", func(t *testing.T) {
		inv := emptyInvoice()
		inv.Title = "A Title"

		assert.Equal(t, "A Title", inv.GetTitle())
	})
}

func TestInvoice_GetTotal(t *testing.T) {

	i := emptyInvoice()
	i.SubTotal = 1000
	i.Tax = 50
	i.ServiceFee = 200
	i.InstallmentFee = 300
	i.Discount = 100

	assert.Equal(t, float64(1450), i.GetTotal())
}

func TestInvoice_Clear(t *testing.T) {

	i := emptyInvoice()
	i.SubTotal = 1000
	i.Tax = 50
	i.ServiceFee = 200
	i.InstallmentFee = 300
	i.Discount = 100

	i.Clear()

	assert.Equal(t, float64(0), i.GetTotal())
	assert.Equal(t, float64(0), i.SubTotal)
	assert.Equal(t, float64(0), i.Discount)
	assert.Equal(t, float64(0), i.Tax)
	assert.Equal(t, float64(0), i.ServiceFee)
	assert.Equal(t, float64(0), i.InstallmentFee)
	assert.Equal(t, Draft, i.State)
	assert.Empty(t, i.LineItems)
	assert.Nil(t, i.Payment)
}

func TestInvoice_UpsertBillingAddress(t *testing.T) {

	t.Run("should create new billing address", func(t *testing.T) {
		i := emptyInvoice()
		i.BillingAddress = &BillingAddress{
			FullName:    "John",
			Email:       "example@example.com",
			PhoneNumber: "021123123",
		}

		assert.NotNil(t, i.BillingAddress)

		i.UpsertBillingAddress("Foo", "foo@bar.com", "08123")

		assert.Equal(t, "Foo", i.BillingAddress.FullName)
		assert.Equal(t, "foo@bar.com", i.BillingAddress.Email)
		assert.Equal(t, "08123", i.BillingAddress.PhoneNumber)
	})

	t.Run("should overwrite billing address", func(t *testing.T) {
		i := emptyInvoice()
		assert.Nil(t, i.BillingAddress)

		i.UpsertBillingAddress("Foo", "foo@bar.com", "08123")

		assert.Equal(t, "Foo", i.BillingAddress.FullName)
		assert.Equal(t, "foo@bar.com", i.BillingAddress.Email)
		assert.Equal(t, "08123", i.BillingAddress.PhoneNumber)
	})

}

type mockFeeReader struct {
	AdminFee       *config.Fee
	InstallmentFee *config.Fee
}

func (m *mockFeeReader) GetAdminFeeConfig(currency string) *config.Fee {
	return m.AdminFee
}

func (m *mockFeeReader) GetInstallmentFeeConfig(currency string) *config.Fee {
	return m.InstallmentFee
}

func (m *mockFeeReader) GetPaymentWaitingTime() *time.Duration {
	return nil
}

func (m *mockFeeReader) GetGateway() payment.Gateway {
	return payment.GatewayMidtrans
}

type mockPaymentMethodFinder struct {
	AdminFee       *config.Fee
	InstallmentFee *config.Fee
	Error          error
}

func (f mockPaymentMethodFinder) FindByPaymentType(ctx context.Context, paymentType payment.PaymentType, opts ...payment.Option) (config.FeeConfigReader, error) {
	if f.Error != nil {
		return nil, f.Error
	}

	m := mockFeeReader{
		AdminFee:       f.AdminFee,
		InstallmentFee: f.InstallmentFee,
	}

	return &m, nil
}

func TestInvoice_UpdatePaymentMethod(t *testing.T) {

	t.Run("should error if payment method is empty", func(t *testing.T) {
		i := emptyInvoice()
		err := i.UpdatePaymentMethod(context.TODO(), nil, nil)
		assert.NotNil(t, err)
		assert.Error(t, InvoiceError{InvoiceErrorPaymentMethodNotSet})
	})

	t.Run("should update fee with installment and admin fee", func(t *testing.T) {

		i := emptyInvoice()

		feeMock := &cfgm.FeeConfigReader{}
		readerMock := &dsm.PaymentConfigReader{}

		readerMock.On("FindByPaymentType", mock.Anything, mock.Anything, mock.Anything).
			Return(feeMock, nil)

		feeMock.On("GetAdminFeeConfig", mock.Anything).Return(
			&config.Fee{
				PercentageVal: 0,
				CurrencyVal:   1000,
				Currency:      "IDR",
			}, nil)
		feeMock.On("GetInstallmentFeeConfig", mock.Anything).Return(
			&config.Fee{
				PercentageVal: 0,
				CurrencyVal:   2000,
				Currency:      "IDR",
			}, nil)

		i.UpdatePaymentMethod(context.TODO(), &Payment{
			PaymentType: payment.SourceBNIVA,
		}, readerMock)
		assert.NotNil(t, i.Payment)
		assert.Equal(t, float64(1000), i.ServiceFee)
		assert.Equal(t, float64(2000), i.InstallmentFee)
	})

}

func TestInvoice_SetItem(t *testing.T) {

	item := LineItem{
		Name:         "",
		Category:     "COURSE",
		MerchantName: "Collegos",
		Currency:     "IDR",
		UnitPrice:    10000,
		Qty:          1,
	}

	t.Run("add new item", func(t *testing.T) {
		i := emptyInvoice()
		assert.Empty(t, i.LineItems)

		err := i.SetItems(context.TODO(), []LineItem{item})
		assert.Nil(t, err)
		assert.NotNil(t, i.LineItems)
		assert.Equal(t, LineItem{
			Name:         "",
			Category:     "COURSE",
			MerchantName: "Collegos",
			Currency:     "IDR",
			UnitPrice:    10000,
			Qty:          1,
		}, i.LineItems[0])
		assert.Equal(t, float64(10000), i.GetSubTotal())
	})

	t.Run("add multiple items", func(t *testing.T) {
		i := emptyInvoice()
		assert.Empty(t, i.LineItems)

		err := i.SetItems(context.TODO(), []LineItem{item, item})
		assert.Nil(t, err)
		assert.NotNil(t, i.LineItems)
		assert.Equal(t, LineItem{
			Name:         "",
			Category:     "COURSE",
			MerchantName: "Collegos",
			Currency:     "IDR",
			UnitPrice:    10000,
			Qty:          1,
		}, i.LineItems[0])
		assert.Len(t, i.LineItems, 2)
		assert.Equal(t, float64(20000), i.GetSubTotal())
	})

}

func TestInvoice_Publish(t *testing.T) {

	t.Run("can't published because payment is not set", func(t *testing.T) {
		i := emptyInvoice()
		err := i.Publish(context.TODO())
		assert.NotNil(t, err)
		assert.Equal(t, InvoiceError{InvoiceErrorPaymentMethodNotSet}, err)
	})

	t.Run("can't published because billing address is not set", func(t *testing.T) {

		i := emptyInvoice()

		feeMock := &cfgm.FeeConfigReader{}
		readerMock := &dsm.PaymentConfigReader{}

		readerMock.On("FindByPaymentType", mock.Anything, mock.Anything, mock.Anything).
			Return(feeMock, nil)

		feeMock.On("GetAdminFeeConfig", mock.Anything).Return(
			&config.Fee{
				PercentageVal: 0,
				CurrencyVal:   1000,
				Currency:      "IDR",
			}, nil)
		feeMock.On("GetInstallmentFeeConfig", mock.Anything).Return(
			&config.Fee{
				PercentageVal: 0,
				CurrencyVal:   2000,
				Currency:      "IDR",
			}, nil)

		i.UpdatePaymentMethod(context.TODO(), &Payment{
			PaymentType: payment.SourceBNIVA,
		}, readerMock)

		err := i.Publish(context.TODO())
		assert.NotNil(t, err)
		assert.Equal(t, InvoiceError{InvoiceErrorBillingAddressNotSet}, err)
	})

	t.Run("published should make invoice available for another day", func(t *testing.T) {
		i := draftInvoice()

		err := i.Publish(context.TODO())
		assert.Nil(t, err)
		assert.Equal(t, Published, i.State)

		assert.Equal(t, 24*time.Hour, i.DueDate.Sub(i.InvoiceDate))

	})
}
