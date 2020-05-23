package midtrans_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/imrenagi/go-payment/gateway/midtrans"
)

func TestNewSnapRequestBuilder(t *testing.T) {

	builder := NewSnapRequestBuilder(dummyInv())
	req, err := builder.Build()
	if err != nil {
		t.Logf("expect no error, got %v", err)
		t.Fail()
	}

	assert.NotEmpty(t, req.TransactionDetails.OrderID)
	assert.Equal(t, int64(5000), req.TransactionDetails.GrossAmt)
	assert.Equal(t, "2020-08-01 08:00:00 +0700", req.Expiry.StartTime)
	assert.Equal(t, "hour", req.Expiry.Unit)
	assert.Equal(t, int64(24), req.Expiry.Duration)
	assert.Equal(t, "Foo", req.CustomerDetail.FName)
	assert.Empty(t, req.CustomerDetail.LName)
	assert.Equal(t, "0812312412", req.CustomerDetail.Phone)
	assert.Equal(t, "Foo", req.CustomerDetail.BillAddr.FName)
	assert.Equal(t, "0812312412", req.CustomerDetail.BillAddr.Phone)
	assert.NotNil(t, 1, len(*req.Items))
}

func TestNewSnapRequestBuilder_WithDiscount(t *testing.T) {

	inv := dummyInv()
	inv.Discount = 500

	builder := NewSnapRequestBuilder(inv)
	req, err := builder.Build()
	if err != nil {
		t.Logf("expect no error, got %v", err)
		t.Fail()
	}

	assert.NotEmpty(t, req.TransactionDetails.OrderID)
	assert.Equal(t, int64(4500), req.TransactionDetails.GrossAmt)
	assert.Equal(t, "2020-08-01 08:00:00 +0700", req.Expiry.StartTime)
	assert.Equal(t, "hour", req.Expiry.Unit)
	assert.Equal(t, int64(24), req.Expiry.Duration)
	assert.Equal(t, "Foo", req.CustomerDetail.FName)
	assert.Empty(t, req.CustomerDetail.LName)
	assert.Equal(t, "0812312412", req.CustomerDetail.Phone)
	assert.Equal(t, "Foo", req.CustomerDetail.BillAddr.FName)
	assert.Equal(t, "0812312412", req.CustomerDetail.BillAddr.Phone)
	assert.Equal(t, 2, len(*req.Items))

	items := *req.Items

	assert.Equal(t, int64(5000), items[0].Price)
	assert.Equal(t, int64(-500), items[1].Price)
}
