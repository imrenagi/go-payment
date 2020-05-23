package midtrans

import (
	gomidtrans "github.com/veritrans/go-midtrans"
)

// NewAkulaku create snaprequest for akulaku payment source
func NewAkulaku(srb *SnapRequestBuilder) (*Akulaku, error) {
	return &Akulaku{
		srb: srb,
	}, nil
}

// Akulaku used for creating snap request for akulaku
type Akulaku struct {
	srb *SnapRequestBuilder
}

// Build ...
func (b *Akulaku) Build() (*gomidtrans.SnapReq, error) {
	req, err := b.srb.Build()
	if err != nil {
		return nil, err
	}

	req.EnabledPayments = []gomidtrans.PaymentType{
		gomidtrans.SourceAkulaku,
	}

	return req, nil
}
