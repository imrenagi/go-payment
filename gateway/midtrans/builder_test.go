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

func TestNewSnapRequestBuilder_LongItemName(t *testing.T) {

	i := dummyInv()
	i.LineItems[0].Name = "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. "

	builder := NewSnapRequestBuilder(i)
	req, err := builder.Build()
	if err != nil {
		t.Logf("expect no error, got %v", err)
		t.Fail()
	}

	items := *req.Items
	assert.Len(t, items[0].Name, 50)
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
