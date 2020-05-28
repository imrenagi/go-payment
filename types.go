package payment

import "time"

// Model is base for database struct
type Model struct {
	ID        uint64     `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `json:"deleted_at" sql:"index"`
}

// PaymentType represent the payment method name
type PaymentType string

const (
	SourceCreditCard PaymentType = "credit_card"
	SourceBNIVA      PaymentType = "bni_va"
	SourcePermataVA  PaymentType = "permata_va"
	SourceBCAVA      PaymentType = "bca_va"
	SourceOtherVA    PaymentType = "other_va"
	SourceAlfamart   PaymentType = "alfamart"
	SourceGopay      PaymentType = "gopay"
	SourceAkulaku    PaymentType = "akulaku"
	SourceOvo        PaymentType = "ovo"
	SourceDana       PaymentType = "dana"
	SourceLinkAja    PaymentType = "linkaja"
	SourceBRIVA      PaymentType = "bri_va"
	SourceMandiriVA  PaymentType = "mandiri_va"
)

// Bank is a bank
type Bank string

const (
	BankBCA Bank = "bca"
	BankBNI Bank = "bni"
	BankBRI Bank = "bri"
)

// InstallmentType shows the type of installment.
type InstallmentType string

const (
	// InstallmentOnline used if the cardholder's card is the same as the the bank providing the installment
	InstallmentOnline InstallmentType = "online"
	// InstallmentOffline used if the cardholders's card might not be the same as the bank providing the installment
	InstallmentOffline InstallmentType = "offline"
)

// Money is just notation for showing the money value and its currency
type Money struct {
	Value         float64 `json:"value"`
	ValuePerMonth float64 `json:"value_per_month,omitempty"`
	Currency      string  `json:"curency"`
}
