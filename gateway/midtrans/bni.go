package midtrans

import (
	gomidtrans "github.com/veritrans/go-midtrans"
)

// NewBNIVA create snaprequest for bni_va payment source
func NewBNIVA(srb *SnapRequestBuilder) (*BNIVA, error) {
	return &BNIVA{
		srb: srb,
	}, nil
}

// BNIVA used for creating snap request for bni_va
type BNIVA struct {
	srb *SnapRequestBuilder
}

// Build ...
func (b *BNIVA) Build() (*gomidtrans.SnapReq, error) {
	req, err := b.srb.Build()
	if err != nil {
		return nil, err
	}

	req.EnabledPayments = []gomidtrans.PaymentType{
		gomidtrans.SourceBNIVA,
	}

	return req, nil
}
