package config

import (
	"time"

	"github.com/imrenagi/go-payment"
)

// NonCard is base for payment method other than cc
type NonCard struct {
	PaymentType payment.PaymentType `yaml:"payment_type" json:"payment_type"`
	IconURLs    []string            `yaml:"icon_urls" json:"icon_urls"`
	Gateway     payment.Gateway     `yaml:"gateway" json:"-"`
	DisplayName string              `yaml:"display_name" json:"display_name"`
	AdminFee    map[string]Fee      `yaml:"admin_fee,omitempty" json:"-"`
	WaitingTime *waitingTime        `yaml:"waiting_time" json:"-"`
}

func (p *NonCard) GetGateway() payment.Gateway {
	return p.Gateway
}

func (p *NonCard) GetAdminFeeConfig(currency string) *Fee {
	if f, ok := p.AdminFee[currency]; ok {
		return &f
	}
	return nil
}

func (p *NonCard) GetInstallmentFeeConfig(currency string) *Fee {
	return nil
}

func (p *NonCard) GetPaymentWaitingTime() *time.Duration {
	var dur time.Duration

	switch p.WaitingTime.Unit {
	case day:
		dur = time.Duration(p.WaitingTime.Duration*24) * time.Hour
	case minute:
		dur = time.Duration(p.WaitingTime.Duration) * time.Minute
	case hour:
		dur = time.Duration(p.WaitingTime.Duration) * time.Hour
	case second:
		dur = time.Duration(p.WaitingTime.Duration) * time.Second
	}
	return &dur
}
