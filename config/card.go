package config

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/imrenagi/go-payment"
)

// Card is just for credit card
type Card struct {
	PaymentType  payment.PaymentType `yaml:"payment_type" json:"payment_type"`
	IconURLs     []string            `yaml:"icon_urls" json:"icon_urls"`
	Installments []Installment       `yaml:"installments" json:"installments"`
}

func (cfg Card) GetInstallment(bank payment.Bank, aType payment.InstallmentType) (*Installment, error) {
	for _, i := range cfg.Installments {
		if i.Bank == bank && i.Type == aType {
			return &i, nil
		}
	}
	return nil, fmt.Errorf("installment %w", payment.ErrNotFound)
}

type Installment struct {
	Gateway     payment.Gateway         `yaml:"gateway" json:"-"`
	DisplayName string                  `yaml:"display_name" json:"display_name"`
	Type        payment.InstallmentType `yaml:"type" json:"type"`
	Bank        payment.Bank            `yaml:"bank" json:"bank"`
	Channel     string                  `yaml:"channel" json:"-"`
	IsActive    bool                    `yaml:"active" json:"-"`
	IsDefault   bool                    `yaml:"default" json:"-"`
	Terms       []InstallmentTerm       `yaml:"terms" json:"terms"`
}

func (i Installment) GetTerm(term int) (*InstallmentTerm, error) {

	for _, t := range i.Terms {
		if term == t.Term {
			return &t, nil
		}
	}
	return nil, fmt.Errorf("installment term %w", payment.ErrNotFound)
}

func (i *Installment) SetValue(val *payment.Money) error {
	var newTerms []InstallmentTerm
	for _, t := range i.Terms {
		t.value = val
		newTerms = append(newTerms, t)
	}
	i.Terms = newTerms
	return nil
}

func (i *Installment) UnmarshalYAML(unmarshal func(interface{}) error) error {

	type Alias Installment
	var ins Alias
	if err := unmarshal(&ins); err != nil {
		return err
	}

	var tmp []InstallmentTerm
	for _, term := range ins.Terms {
		term.Gateway = ins.Gateway
		tmp = append(tmp, term)
	}

	ins.Terms = tmp
	*i = Installment(ins)

	return nil
}

type InstallmentTerm struct {
	Gateway       payment.Gateway `yaml:"-" json:"-"`
	Term          int             `yaml:"term" json:"term"`
	AdminFee      map[string]Fee  `yaml:"admin_fee,omitempty" json:"-"`
	InstalmentFee map[string]Fee  `yaml:"installment_fee,omitempty" json:"-"`
	value         *payment.Money
}

func (p *InstallmentTerm) GetGateway() payment.Gateway {
	return p.Gateway
}

func (p *InstallmentTerm) GetAdminFeeConfig(currency string) *Fee {
	if f, ok := p.AdminFee[currency]; ok {
		return &f
	}
	return nil
}

func (p *InstallmentTerm) GetInstallmentFeeConfig(currency string) *Fee {
	if f, ok := p.InstalmentFee[currency]; ok {
		return &f
	}
	return nil
}

func (p *InstallmentTerm) GetPaymentWaitingTime() *time.Duration {
	dur := 24 * time.Hour
	return &dur
}

func (p *InstallmentTerm) GetAdminFee() *payment.Money {
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

func (p *InstallmentTerm) GetInstallmentFee() *payment.Money {
	if p.value == nil {
		return nil
	}
	if f, ok := p.InstalmentFee[p.value.Currency]; ok {
		val := f.Estimate(p.value.Value)
		return &payment.Money{
			Value:         val,
			ValuePerMonth: math.Ceil((p.value.Value+val)/float64(p.Term)*100) / 100,
			Currency:      p.value.Currency,
		}
	}
	return nil
}

func (p *InstallmentTerm) MarshalJSON() ([]byte, error) {
	type Alias InstallmentTerm

	return json.Marshal(&struct {
		*Alias
		AdminFee       *payment.Money `json:"admin_fee,omitempty"`
		InstallmentFee *payment.Money `json:"installment_fee,omitempty"`
	}{
		Alias:          (*Alias)(p),
		AdminFee:       p.GetAdminFee(),
		InstallmentFee: p.GetInstallmentFee(),
	})
}
