package invoice

import (
	"fmt"

	"github.com/imrenagi/go-payment"
)

type InvoiceError struct {
	Code int
}

const (
	InvoiceErrorPaymentMethodNotSet = iota
	InvoiceErrorBillingAddressNotSet
	InvoiceErrorNoMoreItemExist
	InvoiceErrorInvalidStateTransition
	InvoiceErrorNoPaymentSet
	InvoiceErrorInvalidDiscountValue
)

func (e InvoiceError) Error() string {
	switch e.Code {
	case InvoiceErrorPaymentMethodNotSet:
		return "Payment method must be set before continue"
	case InvoiceErrorBillingAddressNotSet:
		return "Billing address must be set before continue"
	case InvoiceErrorNoMoreItemExist:
		return "Can't remove item from empty invoice"
	case InvoiceErrorInvalidStateTransition:
		return "Can't change invoice state. Action violations"
	case InvoiceErrorInvalidDiscountValue:
		return "Discount value be greater than 0"
	default:
		return "Unknown order error code"
	}
}

func (e InvoiceError) Unwrap() error {
	switch e.Code {
	case InvoiceErrorPaymentMethodNotSet,
		InvoiceErrorBillingAddressNotSet,
		InvoiceErrorNoMoreItemExist,
		InvoiceErrorInvalidStateTransition:
		return fmt.Errorf("%s: %w", e.Error(), payment.ErrCantProceed)
	case InvoiceErrorInvalidDiscountValue:
		return fmt.Errorf("%s: %w", e.Error(), payment.ErrBadRequest)
	default:
		return fmt.Errorf("%s: %w", e.Error(), payment.ErrInternal)
	}
}
