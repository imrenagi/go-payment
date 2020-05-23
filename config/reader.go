package config

import (
	"log"

	"gopkg.in/yaml.v2"
)

// PaymentConfig stores all configuration for payment methods. This struct represent the yaml config file used
// for storing information about payment method fee (admin and installment fee), its max waiting time, also
// payment gateway information
type PaymentConfig struct {
	CardPayment     Card      `yaml:"card_payment" json:"card_payment"`
	BankTransfers   []NonCard `yaml:"bank_transfers" json:"bank_transfers"`
	EWallets        []NonCard `yaml:"ewallets" json:"ewallets"`
	CStores         []NonCard `yaml:"cstores" json:"cstores"`
	CardlessCredits []NonCard `yaml:"cardless_credits" json:"cardless_credits"`
}

// LoadPaymentConfigs reads payment yaml config file
func LoadPaymentConfigs(data []byte) (*PaymentConfig, error) {
	var cfg PaymentConfig
	err := yaml.Unmarshal([]byte(data), &cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return &cfg, nil
}
