package config

import (
	"encoding/json"

	"github.com/imrenagi/go-payment"
)

// NewNonCardPayment returns new NonCardPayment. if value is not nil, the admin fee of this payment method
// can be calculated.
func NewNonCardPayment(cfg NonCard, value *payment.Money) *NonCardPayment {
	return &NonCardPayment{
		NonCard: cfg,
		value:   value,
	}
}

// NonCardPayment represent all payment method other than cards. This includes ewallet, virtual account,
// retail outlet, and cardless credit. This struct might have information about
// the value of a product will be paid by using cards.
type NonCardPayment struct {
	NonCard
	value *payment.Money
}

// GetAdminFee returns the admin fee in money notation
func (p *NonCardPayment) GetAdminFee() *payment.Money {
	if p.value == nil {
		return nil
	}
	if f, ok := p.AdminFee[p.value.Currency]; ok {
		val := f.Estimate(p.value.Value)
		return &payment.Money{
			Value:    val,
			Currency: p.value.Currency,
		}
	}
	return nil
}

// GetInstallmentFee returns nil since non card payment has no installment
func (p *NonCardPayment) GetInstallmentFee() *payment.Money {
	return nil
}

func (p *NonCardPayment) MarshalJSON() ([]byte, error) {
	type Alias NonCardPayment

	return json.Marshal(&struct {
		*Alias
		AdminFee *payment.Money `json:"admin_fee,omitempty"`
	}{
		Alias:    (*Alias)(p),
		AdminFee: p.GetAdminFee(),
	})
}

// NewCardPayment creates a new CardPayment. if value is provided, each installment within this payment can
// calculate its own admin/installment fee in money notation. Otherwise, it just returns the percentage and/or the value
// of the fee.
func NewCardPayment(cfg Card, value *payment.Money) *CardPayment {
	cp := &CardPayment{
		Card:  cfg,
		value: value,
	}

	var installments []Installment
	for _, i := range cp.Installments {
		i.SetValue(value)
		installments = append(installments, i)
	}
	cp.Installments = installments
	return cp
}

// CardPayment represent credit card based payment method. This struct might have information about
// the value of a product will be paid by using cards.
type CardPayment struct {
	Card
	value *payment.Money
}
