package midtrans_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/imrenagi/go-payment/gateway/midtrans"
	gomidtrans "github.com/veritrans/go-midtrans"
)

func TestGopay(t *testing.T) {

	builder := NewSnapRequestBuilder(dummyInv())
	gopay, _ := NewGopay(builder)

	req, err := gopay.Build()
	if err != nil {
		t.Logf("expect no error, got %v", err)
		t.Fail()
	}

	assert.Len(t, req.EnabledPayments, 1)
	assert.Contains(t, req.EnabledPayments, gomidtrans.SourceGopay)
}
