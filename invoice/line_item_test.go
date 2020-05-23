package invoice_test

import (
	"testing"

	. "github.com/imrenagi/go-payment/invoice"
	"github.com/stretchr/testify/assert"
)

func TestLineItem_IncreaseQty(t *testing.T) {

	li := LineItem{}

	assert.Equal(t, 0, li.Qty)
	err := li.IncreaseQty()
	assert.Nil(t, err)
	assert.Equal(t, 1, li.Qty)

	err = li.IncreaseQty()
	assert.Nil(t, err)
	assert.Equal(t, 2, li.Qty)
}

func TestLineItem_DecreaseQty(t *testing.T) {

	t.Run("decrease to 0", func(t *testing.T) {
		li := LineItem{}
		li.IncreaseQty()
		assert.Equal(t, 1, li.Qty)

		err := li.DecreaseQty()
		assert.Nil(t, err)
		assert.Equal(t, 0, li.Qty)
	})

	t.Run("decrease to -1 should error", func(t *testing.T) {
		li := LineItem{}

		err := li.DecreaseQty()
		assert.NotNil(t, err)
		assert.Error(t, err, LineItemError{LineItemErrInvalidQty})
	})

}
