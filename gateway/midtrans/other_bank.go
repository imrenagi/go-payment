package midtrans

import (
	gomidtrans "github.com/veritrans/go-midtrans"
)

// NewOtherBank create snaprequest for other_va payment source
func NewOtherBank(srb *SnapRequestBuilder) (*OtherBank, error) {
	return &OtherBank{
		srb: srb,
	}, nil
}

// OtherBank used for creating snap request for other_va
type OtherBank struct {
	srb *SnapRequestBuilder
}

// Build ...
func (b *OtherBank) Build() (*gomidtrans.SnapReq, error) {
	req, err := b.srb.Build()
	if err != nil {
		return nil, err
	}

	req.EnabledPayments = []gomidtrans.PaymentType{
		gomidtrans.SourceOtherVA,
	}

	return req, nil
}
