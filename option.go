package payment

import (
	"net/http"
	"strconv"
)

// Option is type closure accepting a Options
type Option func(*Options)

// CreditCard tells information about the acquire bank and
// installment used
type CreditCard struct {
	Bank        Bank
	Installment Installment
}

// Installment tells installent type and term
type Installment struct {
	Type InstallmentType
	Term int
}

// Options stores all optional properties for payment purposes
type Options struct {
	Price      *Money
	CreditCard *CreditCard
}

// WithPrice can be used if user want to add optional price information
// used for estimating the admin/installment fee
func WithPrice(price float64, currency string) Option {
	return func(o *Options) {
		o.Price = &Money{
			Value:    price,
			Currency: currency,
		}
	}
}

// WithCreditCard can be used if user want to use the installment feature. It accepts the acquire bank,
// installment type and term
func WithCreditCard(bank Bank, installmentType InstallmentType, installmentTerm int) Option {
	if bank == "" {
		bank = BankBCA
	}
	if installmentType == "" {
		installmentType = InstallmentOffline
	}
	return func(o *Options) {
		o.CreditCard = &CreditCard{
			Bank: bank,
			Installment: Installment{
				Type: installmentType,
				Term: installmentTerm,
			},
		}
	}
}

// NewPaymentMethodListOptions accepts http.Request and returns set of option containing the price and its currency.
func NewPaymentMethodListOptions(r *http.Request) ([]Option, error) {
	r.ParseForm()
	var options []Option
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
