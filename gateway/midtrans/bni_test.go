package midtrans_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/imrenagi/go-payment/gateway/midtrans"
	gomidtrans "github.com/veritrans/go-midtrans"
)

func TestBNIVA(t *testing.T) {

	builder := NewSnapRequestBuilder(dummyInv())
	bni, _ := NewBNIVA(builder)

	req, err := bni.Build()
	if err != nil {
		t.Logf("expect no error, got %v", err)
		t.Fail()
	}

	assert.Len(t, req.EnabledPayments, 1)
	assert.Contains(t, req.EnabledPayments, gomidtrans.SourceBNIVA)
}
