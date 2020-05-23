package invoice_test

import (
	"testing"
	"time"

	. "github.com/imrenagi/go-payment/invoice"
	"github.com/stretchr/testify/assert"
)

func TestStatePublished_State(t *testing.T) {
	i := emptyInvoice()
	state := &PublishedState{}

	t.Run("invoice before due date should be still in published state", func(t *testing.T) {
		now := time.Now()
		i.SetState(state)
		i.DueDate = now.AddDate(0, 0, 1)
		assert.Equal(t, Published, state.State(i))
	})

	t.Run("invoice after due date should be still in failed state", func(t *testing.T) {
		now := time.Now()
		i.SetState(state)
		i.DueDate = now.AddDate(0, 0, -1)
		assert.Equal(t, Failed, state.State(i))
	})
}
