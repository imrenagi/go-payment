package midtrans

import (
	gomidtrans "github.com/veritrans/go-midtrans"
)

// NewBCAVA create snaprequest for bca_va payment source
func NewBCAVA(srb *SnapRequestBuilder) (*BCAVA, error) {
	return &BCAVA{
		srb: srb,
	}, nil
}

// BCAVA used for creating snap request for bca_va
type BCAVA struct {
	srb *SnapRequestBuilder
}

// Build ...
func (b *BCAVA) Build() (*gomidtrans.SnapReq, error) {
	req, err := b.srb.Build()
	if err != nil {
		return nil, err
	}

	req.EnabledPayments = []gomidtrans.PaymentType{
		gomidtrans.SourceBCAVA,
	}

	return req, nil
}
