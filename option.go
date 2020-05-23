package payment

import (
	"net/http"
	"strconv"
)

type PaymentOption func(*PaymentOptions)

type CreditCard struct {
	Bank        Bank
	Installment Installment
}

type Installment struct {
	Type InstallmentType
	Term int
}

type PaymentOptions struct {
	Price      *Money
	CreditCard *CreditCard
}

func WithPrice(price float64, currency string) PaymentOption {
	return func(o *PaymentOptions) {
		o.Price = &Money{
			Value:    price,
			Currency: currency,
		}
	}
}

func WithCreditCard(bank Bank, installmentType InstallmentType, installmentTerm int) PaymentOption {
	if bank == "" {
		bank = BankBCA
	}
	if installmentType == "" {
		installmentType = InstallmentOffline
	}
	return func(o *PaymentOptions) {
		o.CreditCard = &CreditCard{
			Bank: bank,
			Installment: Installment{
				Type: installmentType,
				Term: installmentTerm,
			},
		}
	}
}

func NewPaymentMethodListOptions(r *http.Request) ([]PaymentOption, error) {
	r.ParseForm()
	var options []PaymentOption
	var price float64
	var currency string
	var err error
	if len(r.Form["price"]) > 0 {
		price, err = strconv.ParseFloat(r.Form["price"][0], 64)
		if err != nil {
			return nil, err
		}
	}
	if len(r.Form["currency"]) > 0 {
		currency = r.Form["currency"][0]
	}
	if price > 0 && currency != "" {
		options = append(options, WithPrice(price, currency))
	}

	return options, nil
}
