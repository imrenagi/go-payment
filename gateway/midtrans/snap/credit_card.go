package snap

import (
	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/invoice"
	"github.com/midtrans/midtrans-go"
	midsnap "github.com/midtrans/midtrans-go/snap"
)

func NewCreditCard(inv *invoice.Invoice) (*midsnap.Request, error) {

	ccDetail := &midsnap.CreditCardDetails{
		Secure:      true,
		Bank:        string(midtrans.BankBca),
		Channel:     "migs",
		Type:        "",
		Installment: nil,
	}

	detail := inv.Payment.CreditCardDetail
	if detail != nil {
		switch detail.Installment.Type {
		case payment.InstallmentOffline:
			if detail.Installment.Term > 0 {
				installmentTermsDetail := midsnap.InstallmentTermsDetail{
					Offline: []int8{
						int8(detail.Installment.Term),
					},
				}
				ccDetail.Installment = &midsnap.InstallmentDetail{
					Required: true,
					Terms:    &installmentTermsDetail,
				}
			}
		}
	}

	return newBuilder(inv).
		AddPaymentMethods(midsnap.PaymentTypeCreditCard).
		SetCreditCardDetail(ccDetail).
		Build()
}
