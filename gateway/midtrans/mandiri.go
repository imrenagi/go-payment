package midtrans

import (
	gomidtrans "github.com/veritrans/go-midtrans"
)

// NewMandiriBill create snaprequest for mandiri echannel payment source
func NewMandiriBill(srb *SnapRequestBuilder) (*MandiriBill, error) {
	return &MandiriBill{
		srb: srb,
	}, nil
}

// MandiriBill used for creating snap request for mandiri echannel
type MandiriBill struct {
	srb *SnapRequestBuilder
}

// Build ...
func (b *MandiriBill) Build() (*gomidtrans.SnapReq, error) {
	req, err := b.srb.Build()
	if err != nil {
		return nil, err
	}

	req.EnabledPayments = []gomidtrans.PaymentType{
		gomidtrans.SourceEchannel,
	}

	return req, nil
}
