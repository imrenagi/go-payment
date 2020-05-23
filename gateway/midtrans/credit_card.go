package midtrans

import (
	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/invoice"
	gomidtrans "github.com/veritrans/go-midtrans"
)

// NewCreditCard create snaprequest for bca_va payment source
func NewCreditCard(srb *SnapRequestBuilder, creditCardDetail *invoice.CreditCardDetail) (*CreditCard, error) {
	return &CreditCard{
		srb:              srb,
		creditCardDetail: creditCardDetail,
	}, nil
}

// CreditCard used for creating snap request for bca_va
type CreditCard struct {
	srb              *SnapRequestBuilder
	creditCardDetail *invoice.CreditCardDetail
}

// Build ...
func (b *CreditCard) Build() (*gomidtrans.SnapReq, error) {
	req, err := b.srb.Build()
	if err != nil {
		return nil, err
	}

	req.EnabledPayments = []gomidtrans.PaymentType{
		gomidtrans.SourceCreditCard,
	}

	ccDetail := &gomidtrans.CreditCardDetail{
		Secure:         true,
		Bank:           string(gomidtrans.BankBca),
		Authentication: "3ds",
	}

	if b.creditCardDetail != nil {
		var installmentTermsDetail gomidtrans.InstallmentTermsDetail
		switch b.creditCardDetail.Installment.Type {
		case payment.InstallmentOffline:
			if b.creditCardDetail.Installment.Term > 0 {
				installmentTermsDetail = gomidtrans.InstallmentTermsDetail{
					Offline: []int8{
						int8(b.creditCardDetail.Installment.Term),
					},
				}

				ccDetail.Installment = &gomidtrans.InstallmentDetail{
					Required: true,
					Terms:    &installmentTermsDetail,
				}
			}
		}
	}
	req.CreditCard = ccDetail

	return req, nil
}
