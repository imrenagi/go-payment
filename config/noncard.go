package config

import (
	"time"

	"github.com/imrenagi/go-payment"
)

// NonCard represent the configuration for non cards payment (ewallet, retail outlet, cardless credit,
// virtual account).
type NonCard struct {
	PaymentType payment.PaymentType `yaml:"payment_type" json:"payment_type"`
	IconURLs    []string            `yaml:"icon_urls" json:"icon_urls"`
	Gateway     payment.Gateway     `yaml:"gateway" json:"-"`
	DisplayName string              `yaml:"display_name" json:"display_name"`
	AdminFee    map[string]Fee      `yaml:"admin_fee,omitempty" json:"-"`
	WaitingTime *waitingTime        `yaml:"waiting_time" json:"-"`
}

// GetGateway returns the payment gateway used for this payment methods
func (p *NonCard) GetGateway() payment.Gateway {
	return p.Gateway
}

// GetAdminFeeConfig returns the fee configuration for a given currency.
func (p *NonCard) GetAdminFeeConfig(currency string) *Fee {
	if f, ok := p.AdminFee[currency]; ok {
		return &f
	}
	return nil
}

// GetInstallmentFeeConfig returns nil since the non card payment method doesn't
// has installment feature
func (p *NonCard) GetInstallmentFeeConfig(currency string) *Fee {
	return nil
}

// GetPaymentWaitingTime is the max waiting time for payment completion after
// customer initiate the payment
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
