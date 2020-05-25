package manage

import (
	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/config"
)

func paymentMethodListFromConfig(cfg *config.PaymentConfig, subtotal *payment.Money) *PaymentMethodList {

	if cfg == nil {
		return nil
	}

	var bankTransfers []config.NonCardPayment
	for _, bt := range cfg.BankTransfers {
		payment := config.NewNonCardPayment(bt, subtotal)
		if payment != nil {
			bankTransfers = append(bankTransfers, *payment)
		}
	}

	var ewallets []config.NonCardPayment
	for _, ew := range cfg.EWallets {
		payment := config.NewNonCardPayment(ew, subtotal)
		if payment != nil {
			ewallets = append(ewallets, *payment)
		}
	}

	var cstores []config.NonCardPayment
	for _, cs := range cfg.CStores {
		payment := config.NewNonCardPayment(cs, subtotal)
		if payment != nil {
			cstores = append(cstores, *payment)
		}
	}

	var cardlessCredits []config.NonCardPayment
	for _, cl := range cfg.CardlessCredits {
		payment := config.NewNonCardPayment(cl, subtotal)
		if payment != nil {
			cardlessCredits = append(cardlessCredits, *payment)
		}
	}

	return &PaymentMethodList{
		CardPayment:     config.NewCardPayment(cfg.CardPayment, subtotal),
		BankTransfers:   bankTransfers,
		EWallets:        ewallets,
		CStores:         cstores,
		CardlessCredits: cardlessCredits,
	}

}

// PaymentMethodList is the payment method list showed to the user
type PaymentMethodList struct {
	CardPayment     *config.CardPayment     `json:"card_payment"`
	BankTransfers   []config.NonCardPayment `json:"bank_transfers"`
	EWallets        []config.NonCardPayment `json:"ewallets"`
	CStores         []config.NonCardPayment `json:"cstores"`
	CardlessCredits []config.NonCardPayment `json:"cardless_credits"`
}
