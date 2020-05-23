package config

import (
	"log"

	"gopkg.in/yaml.v2"
)

func LoadPaymentConfigs(data []byte) (*PaymentConfig, error) {
	var cfg PaymentConfig
	err := yaml.Unmarshal([]byte(data), &cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return &cfg, nil
}
