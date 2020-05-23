package midtrans_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/imrenagi/go-payment/gateway/midtrans"
	gomidtrans "github.com/veritrans/go-midtrans"
)

func TestPermataVA(t *testing.T) {

	builder := NewSnapRequestBuilder(dummyInv())
	permata, _ := NewPermataVA(builder)

	req, err := permata.Build()
	if err != nil {
		t.Logf("expect no error, got %v", err)
		t.Fail()
	}

	assert.Len(t, req.EnabledPayments, 1)
	assert.Contains(t, req.EnabledPayments, gomidtrans.SourcePermataVA)
}
