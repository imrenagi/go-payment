package payment

import "time"

type Model struct {
	ID        uint64     `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `json:"deleted_at" sql:"index"`
}

// PaymentType value
type PaymentType string

const (
	SourceCreditCard PaymentType = "credit_card"
	SourceBNIVA      PaymentType = "bni_va"
	SourcePermataVA  PaymentType = "permata_va"
	SourceBCAVA      PaymentType = "bca_va"
	SourceOtherVA    PaymentType = "other_va"
	SourceEchannel   PaymentType = "echannel"
	SourceAlfamart   PaymentType = "alfamart"
	SourceGopay      PaymentType = "gopay"
	SourceAkulaku    PaymentType = "akulaku"
	SourceOvo        PaymentType = "ovo"
	SourceDana       PaymentType = "dana"
	SourceLinkAja    PaymentType = "linkaja"
)

type Bank string

const (
	BankBCA Bank = "bca"
)

type InstallmentType string

const (
	InstallmentOnline  InstallmentType = "online"
	InstallmentOffline InstallmentType = "offline"
)

type Money struct {
	Value         float64 `json:"value"`
	ValuePerMonth float64 `json:"value_per_month,omitempty"`
	Currency      string  `json:"curency"`
}

func NewIDR(val float64) *Money {
	return &Money{
		Value:    val,
		Currency: "IDR",
	}
}
