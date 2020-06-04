package xendit

import (
	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/subscription"
	xgo "github.com/xendit/xendit-go"
)

// NewStatus convert xendit status string to subscripiton status
func NewStatus(s string) subscription.Status {
	switch s {
	case "ACTIVE":
		return subscription.StatusActive
	case "PAUSED":
		return subscription.StatusPaused
	default:
		return subscription.StatusStop
	}
}

// NewPaymentSource converts xendit payment method to payment.PaymentType
func NewPaymentSource(s string) payment.PaymentType {
	switch s {
	case "BCA":
		return payment.SourceBCAVA
	case "BRI":
		return payment.SourceBRIVA
	case "MANDIRI":
		return payment.SourceMandiriVA
	case "BNI":
		return payment.SourceBNIVA
	case "PERMATA":
		return payment.SourcePermataVA
	case "ALFAMART":
		return payment.SourceAlfamart
	case "CREDIT_CARD":
		return payment.SourceCreditCard
	case "OVO":
		return payment.SourceOvo
	default:
		return ""
	}
}

func missedPaymentAction(enum subscription.MissedPaymentAction) xgo.MissedPaymentActionEnum {
	switch enum {
	case subscription.MissedPaymentActionIgnore:
		return xgo.MissedPaymentActionIgnore
	case subscription.MissedPaymentActionStop:
		return xgo.MissedPaymentActionStop
	default:
		return ""
	}
}

func paymentIntervalUnit(enum subscription.IntervalUnit) xgo.RecurringPaymentIntervalEnum {
	switch enum {
	case subscription.IntervalUnitDay:
		return xgo.RecurringPaymentIntervalDay
	case subscription.IntervalUnitWeek:
		return xgo.RecurringPaymentIntervalWeek
	case subscription.IntervalUnitMonth:
		return xgo.RecurringPaymentIntervalMonth
	default:
		return ""
	}
}
