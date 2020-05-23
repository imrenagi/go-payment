package config

type PaymentConfig struct {
	CardPayment     Card      `yaml:"card_payment" json:"card_payment"`
	BankTransfers   []NonCard `yaml:"bank_transfers" json:"bank_transfers"`
	EWallets        []NonCard `yaml:"ewallets" json:"ewallets"`
	CStores         []NonCard `yaml:"cstores" json:"cstores"`
	CardlessCredits []NonCard `yaml:"cardless_credits" json:"cardless_credits"`
}
