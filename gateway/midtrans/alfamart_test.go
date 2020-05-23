package midtrans_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/imrenagi/go-payment/gateway/midtrans"
	gomidtrans "github.com/veritrans/go-midtrans"
)

func TestAlfamart(t *testing.T) {

	builder := NewSnapRequestBuilder(dummyInv())
	alfamart, _ := NewAlfamart(builder)

	req, err := alfamart.Build()
	if err != nil {
		t.Logf("expect no error, got %v", err)
		t.Fail()
	}

	assert.Len(t, req.EnabledPayments, 1)
	assert.Contains(t, req.EnabledPayments, gomidtrans.SourceAlfamart)
}
