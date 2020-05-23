package config

import (
	"math"
	"time"

	"github.com/imrenagi/go-payment"
)

type FeeConfigReader interface {
	GetGateway() payment.Gateway
	GetPaymentWaitingTime() *time.Duration
	GetAdminFeeConfig(currency string) *Fee
	GetInstallmentFeeConfig(currency string) *Fee
}

type Fee struct {
	PercentageVal float64 `yaml:"val_percentage"`
	CurrencyVal   float64 `yaml:"val_currency"`
	Currency      string  `yaml:"currency" json:"currency"`
}

// Estimate estimates the fee for a given value. This value will be ceiled.
func (f Fee) Estimate(val float64) float64 {
	return math.Ceil((((f.PercentageVal/100)*val + f.CurrencyVal) * 100) / 100)
}
