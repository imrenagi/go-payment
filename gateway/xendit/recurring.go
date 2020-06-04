package xendit

import (
	"fmt"
	"os"

	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/subscription"

	xrp "github.com/xendit/xendit-go/recurringpayment"
)

// NewRecurringChargeRequestBuilder builder for building the recurring charge request
func NewRecurringChargeRequestBuilder(s *subscription.Subscription) *RecurringChargeRequestBuilder {

	b := &RecurringChargeRequestBuilder{
		request: &xrp.CreateParams{
			ExternalID:          s.Number,
			ShouldSendEmail:     &s.ShouldSendEmail,
			MissedPaymentAction: missedPaymentAction(s.MissedPaymentAction),
			Recharge:            &s.Recharge,
			ChargeImmediately:   &s.ChargeImmediately,
			SuccessRedirectURL:  fmt.Sprintf("%s%s", os.Getenv("WEB_BASE_URL"), os.Getenv("SUCCESS_REDIRECT_PATH")),
			FailureRedirectURL:  fmt.Sprintf("%s%s", os.Getenv("WEB_BASE_URL"), os.Getenv("FAILED_REDIRECT_PATH")),
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
