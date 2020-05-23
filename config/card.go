package config

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/imrenagi/go-payment"
)

// Card represent the credit card payment config retrieved from the yaml config file
type Card struct {
	PaymentType  payment.PaymentType `yaml:"payment_type" json:"payment_type"`
	IconURLs     []string            `yaml:"icon_urls" json:"icon_urls"`
	Installments []Installment       `yaml:"installments" json:"installments"`
}

// GetInstallment returns an installment information for a given bank and its type
func (cfg Card) GetInstallment(bank payment.Bank, aType payment.InstallmentType) (*Installment, error) {
	for _, i := range cfg.Installments {
		if i.Bank == bank && i.Type == aType {
			return &i, nil
		}
	}
	return nil, fmt.Errorf("installment %w", payment.ErrNotFound)
}

// Installment contains information about the package of the installment, issuer bank, its type,
// and installment terms along with its fee information
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

// GetTerm finds an installment given its term. If it doesn't exist, it returns
// ErrNotFound error
func (i Installment) GetTerm(term int) (*InstallmentTerm, error) {

	for _, t := range i.Terms {
		if term == t.Term {
			return &t, nil
		}
	}
	return nil, fmt.Errorf("installment term %w", payment.ErrNotFound)
}

// SetValue sets the value of money which will be used for admin/installment fee
func (i *Installment) SetValue(val *payment.Money) error {
	var newTerms []InstallmentTerm
	for _, t := range i.Terms {
		t.value = val
		newTerms = append(newTerms, t)
	}
	i.Terms = newTerms
	return nil
}

// UnmarshalYAML custom unmarshall for installment. This will add payment gateway info
// to each terms in an installment
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

// InstallmentTerm stores information about the admin/installment fee applicable for a particular term
type InstallmentTerm struct {
	Gateway       payment.Gateway `yaml:"-" json:"-"`
	Term          int             `yaml:"term" json:"term"`
	AdminFee      map[string]Fee  `yaml:"admin_fee,omitempty" json:"-"`
	InstalmentFee map[string]Fee  `yaml:"installment_fee,omitempty" json:"-"`
	value         *payment.Money
}

// GetGateway returns payment gateway for current installment term
func (p *InstallmentTerm) GetGateway() payment.Gateway {
	return p.Gateway
}

// GetAdminFeeConfig returns admin fee rules of an installment term
func (p *InstallmentTerm) GetAdminFeeConfig(currency string) *Fee {
	if f, ok := p.AdminFee[currency]; ok {
		return &f
	}
	return nil
}

// GetInstallmentFeeConfig returns installment fee rules of an installment term
func (p *InstallmentTerm) GetInstallmentFeeConfig(currency string) *Fee {
	if f, ok := p.InstalmentFee[currency]; ok {
		return &f
	}
	return nil
}

// GetPaymentWaitingTime is the max waiting time for payment completion after
// customer initiate the payment
func (p *InstallmentTerm) GetPaymentWaitingTime() *time.Duration {
	dur := 24 * time.Hour
	return &dur
}

// GetAdminFee returns the admin fee of an installment term in money notation
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

// GetInstallmentFee returns the installment fee of an installment term in money notation
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

// MarshalJSON augments admin and installment fee in money notation to
// installment term struct
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
