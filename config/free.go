package config

import (
	"time"

	"github.com/imrenagi/go-payment"
)

// NewFreeFee returns a fee config reader which has no fee at all
func NewFreeFee(gateway payment.Gateway) FeeConfigReader {
	return &freeFee{
		gateway: gateway,
	}
}

type freeFee struct {
	gateway payment.Gateway
}

func (f freeFee) GetGateway() payment.Gateway {
	return f.gateway
}

func (f freeFee) GetPaymentWaitingTime() *time.Duration {
	dur := 0 * time.Hour
	return &dur
}

func (f freeFee) GetAdminFeeConfig(currency string) *Fee {
	return &Fee{}
}

func (f freeFee) GetInstallmentFeeConfig(currency string) *Fee {
	return &Fee{}
}
