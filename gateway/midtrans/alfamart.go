package midtrans

import (
	gomidtrans "github.com/veritrans/go-midtrans"
)

// NewAlfamart create snaprequest for alfamart payment source
func NewAlfamart(srb *SnapRequestBuilder) (*Alfamart, error) {
	return &Alfamart{
		srb: srb,
	}, nil
}

// Alfamart used for creating snap request for alfamart
type Alfamart struct {
	srb *SnapRequestBuilder
}

// Build ...
func (b *Alfamart) Build() (*gomidtrans.SnapReq, error) {
	req, err := b.srb.Build()
	if err != nil {
		return nil, err
	}

	req.EnabledPayments = []gomidtrans.PaymentType{
		gomidtrans.SourceAlfamart,
	}

	return req, nil
}
