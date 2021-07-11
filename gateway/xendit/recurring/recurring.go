package recurring

import (
	"fmt"
	"os"

	xgo "github.com/xendit/xendit-go"

	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/subscription"

	xrp "github.com/xendit/xendit-go/recurringpayment"
)

// New creates recurring charge params
func New(s *subscription.Subscription) (*xrp.CreateParams, error) {
	return newRecurringChargeRequestBuilder(s).Build()
}

// newRecurringChargeRequestBuilder builder for building the recurring charge request
func newRecurringChargeRequestBuilder(s *subscription.Subscription) *RecurringChargeRequestBuilder {

	b := &RecurringChargeRequestBuilder{
		request: &xrp.CreateParams{
			ExternalID:          s.Number,
			ShouldSendEmail:     &s.ShouldSendEmail,
			MissedPaymentAction: missedPaymentAction(s.MissedPaymentAction),
			Recharge:            &s.Recharge,
			ChargeImmediately:   &s.ChargeImmediately,
			SuccessRedirectURL:  os.Getenv("RECURRING_SUCCESS_REDIRECT_URL"),
			FailureRedirectURL:  os.Getenv("RECURRING_FAILED_REDIRECT_URL"),
		},
	}

	return b.SetSchedule(s).
		SetPrice(s).
		SetBasicInfo(s).
		SetCustomerData(s)
}

type RecurringChargeRequestBuilder struct {
	request *xrp.CreateParams
}

func (b *RecurringChargeRequestBuilder) SetSchedule(s *subscription.Subscription) *RecurringChargeRequestBuilder {
	b.request.StartDate = s.Schedule.StartAt
	b.request.Interval = paymentIntervalUnit(s.Schedule.IntervalUnit)
	b.request.IntervalCount = s.Schedule.Interval
	b.request.InvoiceDuration = int(s.InvoiceDuration.Seconds())
	b.request.TotalRecurrence = s.TotalReccurence

	return b
}

func (b *RecurringChargeRequestBuilder) SetPrice(s *subscription.Subscription) *RecurringChargeRequestBuilder {
	b.request.Amount = s.Amount
	return b
}

func (b *RecurringChargeRequestBuilder) SetCustomerData(s *subscription.Subscription) *RecurringChargeRequestBuilder {
	// TODO change this
	b.request.PayerEmail = s.UserID
	return b
}

func (b *RecurringChargeRequestBuilder) SetBasicInfo(s *subscription.Subscription) *RecurringChargeRequestBuilder {
	b.request.Description = fmt.Sprintf("%s: %s", s.Name, s.Description)
	return b
}

func (b *RecurringChargeRequestBuilder) Build() (*xrp.CreateParams, error) {

	if b.request.MissedPaymentAction == "" {
		return nil, fmt.Errorf("unkown missed payment action %w", payment.ErrBadRequest)
	}
	if b.request.Interval == "" {
		return nil, fmt.Errorf("unkown recurring interval unit %w", payment.ErrBadRequest)
	}

	return b.request, nil
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