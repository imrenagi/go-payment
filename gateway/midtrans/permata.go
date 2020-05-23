package midtrans

import (
	gomidtrans "github.com/veritrans/go-midtrans"
)

// NewPermataVA create snaprequest for permata_va payment source
func NewPermataVA(srb *SnapRequestBuilder) (*PermataVA, error) {
	return &PermataVA{
		srb: srb,
	}, nil
}

// PermataVA used for creating snap request for permata_va
type PermataVA struct {
	srb *SnapRequestBuilder
}

// Build ...
func (b *PermataVA) Build() (*gomidtrans.SnapReq, error) {
	req, err := b.srb.Build()
	if err != nil {
		return nil, err
	}

	req.EnabledPayments = []gomidtrans.PaymentType{
		gomidtrans.SourcePermataVA,
	}

	return req, nil
}
