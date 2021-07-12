package snap_test

import (
	"testing"

	midsnap "github.com/midtrans/midtrans-go/snap"
	"github.com/stretchr/testify/assert"

	"github.com/imrenagi/go-payment/gateway/midtrans/snap"
	"github.com/imrenagi/go-payment/invoice"
)

func TestNewMandiriVA(t *testing.T) {
	type args struct {
		inv *invoice.Invoice
	}
	tests := []struct {
		name    string
		args    args
		want    *midsnap.Request
		wantErr error
	}{
		{
			name: "standard mandiri va request",
			args: args{inv: dummyInv},
			want: &midsnap.Request{
				EnabledPayments: []midsnap.SnapPaymentType{
					midsnap.PaymentTypeEChannel,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := snap.NewMandiriVA(tt.args.inv)
			assert.Equal(t, tt.wantErr, err)
			assert.Contains(t, got.EnabledPayments, midsnap.PaymentTypeEChannel)
		})
	}
}
