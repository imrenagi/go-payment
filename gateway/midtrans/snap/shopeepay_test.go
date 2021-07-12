package snap_test

import (
	"testing"

	midsnap "github.com/midtrans/midtrans-go/snap"
	"github.com/stretchr/testify/assert"

	"github.com/imrenagi/go-payment/gateway/midtrans/snap"
	"github.com/imrenagi/go-payment/invoice"
)

func TestNewShopeePay(t *testing.T) {
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
			name: "standard shopeepay request",
			args: args{inv: dummyInv},
			want: &midsnap.Request{
				EnabledPayments: []midsnap.SnapPaymentType{
					midsnap.PaymentTypeShopeepay,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := snap.NewShopeePay(tt.args.inv)
			assert.Equal(t, tt.wantErr, err)
			assert.EqualValues(t, tt.want.ShopeePay, got.ShopeePay)
			assert.Contains(t, got.EnabledPayments, midsnap.PaymentTypeShopeepay)
		})
	}
}
