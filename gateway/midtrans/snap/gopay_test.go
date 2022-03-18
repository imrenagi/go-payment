package snap_test

import (
	"testing"

	"github.com/midtrans/midtrans-go"
	midsnap "github.com/midtrans/midtrans-go/snap"
	"github.com/stretchr/testify/assert"

	. "github.com/imrenagi/go-payment/gateway/midtrans/snap"
	"github.com/imrenagi/go-payment/invoice"
)

func TestNewGopay(t *testing.T) {

	tests := []struct {
		name    string
		inv     *invoice.Invoice
		req     *midsnap.Request
		wantErr error
	}{
		{
			name: "successfully create gopay request",
			inv:  dummyInv,
			req: &midsnap.Request{
				TransactionDetails: midtrans.TransactionDetails{
					OrderID:  "a-random-invoice-number",
					GrossAmt: 15700,
				},
				Items: &[]midtrans.ItemDetails{
					{
						ID:       "1",
						Name:     "random-item",
						Price:    15000,
						Qty:      1,
						Category: "HOME",
					},
					{
						ID:       "adminfee",
						Name:     "Biaya Admin",
						Price:    500,
						Qty:      1,
						Category: "FEE",
					},
					{
						ID:       "installmentfee",
						Name:     "Installment Fee",
						Price:    1000,
						Qty:      1,
						Category: "FEE",
					},
					{
						ID:       "discount",
						Name:     "Discount",
						Price:    -1000,
						Qty:      1,
						Category: "DISCOUNT",
					},
					{
						ID:       "tax",
						Name:     "Tax",
						Price:    200,
						Qty:      1,
						Category: "TAX",
					},
				},
				CustomerDetail: &midtrans.CustomerDetails{
					FName: "John Doe",
					LName: "",
					Email: "foo@bar.com",
					Phone: "08111231234",
					BillAddr: &midtrans.CustomerAddress{
						FName: "John Doe",
						LName: "",
						Phone: "08111231234",
					},
				},
				EnabledPayments: []midsnap.SnapPaymentType{
					midsnap.PaymentTypeGopay,
				},
				Expiry: &midsnap.ExpiryDetails{
					StartTime: "2021-01-01 07:00:00 +0700",
					Unit:      "minute",
					Duration:  1440,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := NewGopay(tt.inv)
			assert.Equal(t, tt.wantErr, err)
			assert.EqualValues(t, tt.req, req)
		})
	}
}
