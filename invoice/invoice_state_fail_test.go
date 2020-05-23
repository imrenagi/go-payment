package invoice_test

import (
	"context"
	"testing"

	. "github.com/imrenagi/go-payment/invoice"
	"github.com/stretchr/testify/assert"
)

// func TestFailedState_State(t *testing.T) {
// 	s := FailedState{}
// 	assert.Equal(t, Failed, s.State())
// }

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
		i.SetPaymentMethod(&Payment{})
		i.SetBillingAddress("foo", "foo@bar.com", "0811")

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
