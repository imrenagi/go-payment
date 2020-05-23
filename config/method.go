package config

import (
	"encoding/json"

	"github.com/imrenagi/go-payment"
)

func NewNonCardPayment(cfg NonCard, value *payment.Money) *NonCardPayment {
	return &NonCardPayment{
		NonCard: cfg,
		value:   value,
	}
}

type NonCardPayment struct {
	NonCard
	value *payment.Money
}

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

func NewCardPayment(cfg Card, value *payment.Money) *CardPayment {
	x := &CardPayment{
		Card:  cfg,
		value: value,
	}

	var installments []Installment
	for _, i := range x.Installments {
		i.SetValue(value)
		installments = append(installments, i)
	}
	x.Installments = installments
	return x
}

type CardPayment struct {
	Card
	value *payment.Money
}
