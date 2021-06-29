package midtrans

import (
	"crypto/sha512"
	"fmt"
	"io"
	"time"

	"github.com/imrenagi/go-payment"
)

// TransactionStatus is object used to store notification from midtrans
type TransactionStatus struct {
	ID                     uint64    `json:"id" gorm:"primary_key"`
	CreatedAt              time.Time `json:"created_at" gorm:"not null;"`
	UpdatedAt              time.Time `json:"updated_at" gorm:"not null;"`
	StatusCode             string    `json:"status_code" gorm:"not null"`
	StatusMessage          string    `json:"status_message" gorm:"type:text;not null"`
	SignKey                string    `json:"signature_key" gorm:"type:text;column:signature_key;not null"`
	Bank                   string    `json:"bank"`
	FraudStatus            string    `json:"fraud_status" gorm:"not null"`
	PaymentType            string    `json:"payment_type" gorm:"not null"`
	OrderID                string    `json:"order_id" gorm:"not null;unique_index:order_id_k"`
	TransactionID          string    `json:"transaction_id"  gorm:"not null;unique_index:transaction_id_k"`
	TransactionTime        time.Time `json:"-" gorm:"not null"`
	TransactionStatus      string    `json:"transaction_status" gorm:"not null"`
	GrossAmount            string    `json:"gross_amount" gorm:"not null"`
	MaskedCard             string    `json:"masked_card"`
	Currency               string    `json:"currency" gorm:"not null"`
	CardType               string    `json:"card_type"`
	ChannelResponseCode    string    `json:"channel_response_code" gorm:"not null"`
	ChannelResponseMessage string    `json:"channel_response_message"`
	ApprovalCode           string    `json:"approval_code"`
}

// TableName returns the gorm table name
func (TransactionStatus) TableName() string {
	return "midtrans_transaction_status"
}

// IsValid checks whether the status sent is indeed sent by midtrans by validating the
// data against its authentication key.
// See https://snap-docs.midtrans.com/#handling-notifications
func (m TransactionStatus) IsValid(authKey string) error {
	key := fmt.Sprintf("%s%s%s%s", m.OrderID, m.StatusCode, m.GrossAmount, authKey)
	h512 := sha512.New()
	io.WriteString(h512, key)
	if fmt.Sprintf("%x", h512.Sum(nil)) != m.SignKey {
		return fmt.Errorf("%w: Invalid sign key", payment.ErrBadRequest)
	}
	return nil
}
