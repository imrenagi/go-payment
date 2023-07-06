package invoice

import (
	"fmt"

	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/util/validator"
)

var (
	emailValidator = validator.EmailValidator{}
	phoneValidator = validator.PhoneNumberValidator{}
)

// NewBillingAddress ...
func NewBillingAddress(fullName, email, phoneNumber string) (*BillingAddress, error) {
	ba := &BillingAddress{}
	err := ba.Update(fullName, email, phoneNumber)
	if err != nil {
		return nil, err
	}
	return ba, nil
}

// BillingAddress stores information about account making the payment
type BillingAddress struct {
	payment.Model
	FullName    string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	InvoiceID   uint64 `json:"-" gorm:"index:bil_addr_invoice_id_k;not null"`
}

func (b *BillingAddress) Update(name, email, phoneNumber string) error {
	err := b.setName(name)
	if err != nil {
		return err
	}

	err = b.setEmail(email)
	if err != nil {
		return err
	}

	err = b.setPhoneNumber(phoneNumber)
	if err != nil {
		return err
	}
	return nil
}

func (b *BillingAddress) setName(name string) error {
	if name == "" {
		return fmt.Errorf("%w: name must not be empty", payment.ErrBadRequest)
	}
	b.FullName = name
	return nil
}

func (b *BillingAddress) setEmail(email string) error {
	if !emailValidator.IsValid(email) {
		return fmt.Errorf("%w: email must be valid email address", payment.ErrBadRequest)
	}
	b.Email = email
	return nil
}

func (b *BillingAddress) setPhoneNumber(number string) error {
	if len(number) > 0 {
		if !phoneValidator.IsValid(number) {
			return fmt.Errorf("%w: phone number must be valid phone number", payment.ErrBadRequest)
		}
		b.PhoneNumber = number
	}
	return nil
}

func (BillingAddress) TableName() string {
	return "goldfish_invoice_billing_addresses"
}
