package midtrans

import (
	gomidtrans "github.com/veritrans/go-midtrans"
)

// NewGopay create snaprequest for gopay payment source
func NewGopay(srb *SnapRequestBuilder) (*Gopay, error) {
	return &Gopay{
		srb: srb,
	}, nil
}

// Gopay used for creating snap request for gopay
type Gopay struct {
	srb *SnapRequestBuilder
}

// Build ...
func (b *Gopay) Build() (*gomidtrans.SnapReq, error) {
	req, err := b.srb.Build()
	if err != nil {
		return nil, err
	}

	req.EnabledPayments = []gomidtrans.PaymentType{
		gomidtrans.SourceGopay,
	}

	return req, nil
}
