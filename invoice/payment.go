package invoice

import (
	"fmt"
	"time"

	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/config"
)

func NewPayment(cfg config.FeeConfigReader, payType payment.PaymentType, ccDetail *CreditCardDetail) (*Payment, error) {

	if string(payType) == "" {
		return nil, fmt.Errorf("%w: payment_method must be set", payment.ErrBadRequest)
	}

	var waitTimeMS int64
	if cfg.GetPaymentWaitingTime() != nil {
		t := cfg.GetPaymentWaitingTime()
		waitTimeMS = int64(*t / time.Millisecond)
	}

	return &Payment{
		Gateway:          cfg.GetGateway().String(),
		PaymentType:      payType,
		CreditCardDetail: ccDetail,
		WaitingTimeMS:    &waitTimeMS,
	}, nil
}

type Payment struct {
	payment.Model
	Gateway          string              `json:"gateway" gorm:"not null"`
	PaymentType      payment.PaymentType `json:"payment_type"`
	Token            string              `json:"token"`
	RedirectURL      string              `json:"redirect_url" gorm:"type:text"`
	InvoiceID        uint64              `json:"-" gorm:"sql:index;"`
	TransactionID    string              `json:"transaction_id" gorm:"sql:index;"`
	WaitingTimeMS    *int64              `json:"-"`
	CreditCardDetail *CreditCardDetail   `json:"credit_card,omitempty" gorm:"ForeignKey:PaymentID"`
}

func (p *Payment) WaitingDuration() *time.Duration {
	if p.WaitingTimeMS == nil {
		return nil
	}
	dur := time.Duration(*p.WaitingTimeMS) * time.Millisecond
	return &dur
}

func (p *Payment) Reset() error {
	p.Gateway = ""
	p.Token = ""
	p.RedirectURL = ""
	p.TransactionID = ""
	p.PaymentType = payment.PaymentType("")

	if p.CreditCardDetail != nil {
		now := time.Now()
		p.CreditCardDetail.DeletedAt = &now
	}

	return nil
}

func (Payment) TableName() string {
	return "invoice_payment_info"
}

type CreditCardDetail struct {
	payment.Model
	PaymentID   uint64       `json:"-" gorm:"index:cc_payment_id_k"`
	Installment Installment  `json:"installment" gorm:"embedded;embedded_prefix:installment_"`
	Bank        payment.Bank `json:"bank"`
}

type Installment struct {
	Type payment.InstallmentType `json:"type"`
	Term int                     `json:"term"`
}

func (CreditCardDetail) TableName() string {
	return "invoice_payment_cc_details"
}
