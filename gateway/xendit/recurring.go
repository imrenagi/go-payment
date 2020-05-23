package xendit

import (
	"fmt"
	"os"
	"time"

	"github.com/imrenagi/go-payment/invoice"
	xgo "github.com/xendit/xendit-go"
	xrp "github.com/xendit/xendit-go/recurringpayment"
)

func NewRecurringChargeRequestBuilder(inv *invoice.Invoice) *RecurringChargeRequestBuilder {

	var truePtr bool = true

	b := &RecurringChargeRequestBuilder{
		request: &xrp.CreateParams{
			ExternalID:          inv.Number,
			Interval:            xgo.RecurringPaymentIntervalMonth,
			ShouldSendEmail:     &truePtr,
			MissedPaymentAction: xgo.MissedPaymentActionIgnore,
			Recharge:            &truePtr,
			ChargeImmediately:   &truePtr,
			SuccessRedirectURL:  fmt.Sprintf("%s/donate/thanks", os.Getenv("WEB_BASE_URL")),
			FailureRedirectURL:  fmt.Sprintf("%s/donate/error", os.Getenv("WEB_BASE_URL")),
		},
	}

	return b.SetPrice(inv).
		SetCustomerData(inv).
		SetItemDetails(inv).
		SetSubscriptionTime(inv)
}

type RecurringChargeRequestBuilder struct {
	request *xrp.CreateParams
}

func (b *RecurringChargeRequestBuilder) SetSubscriptionTime(inv *invoice.Invoice) *RecurringChargeRequestBuilder {
	// TODO change this based on value from invoice
	b.request.Interval = xgo.RecurringPaymentIntervalMonth
	b.request.IntervalCount = 1
	b.request.InvoiceDuration = int((24 * time.Hour).Seconds())
	b.request.TotalRecurrence = 0

	return b
}

func (b *RecurringChargeRequestBuilder) SetPrice(inv *invoice.Invoice) *RecurringChargeRequestBuilder {
	b.request.Amount = inv.GetTotal()
	return b
}

func (b *RecurringChargeRequestBuilder) SetCustomerData(inv *invoice.Invoice) *RecurringChargeRequestBuilder {
	b.request.PayerEmail = inv.BillingAddress.Email
	return b
}

func (b *RecurringChargeRequestBuilder) SetItemDetails(inv *invoice.Invoice) *RecurringChargeRequestBuilder {

	if inv.LineItem == nil {
		return b
	}

	b.request.Description = fmt.Sprintf("%s (%dx): %s",
		inv.LineItem.Name, inv.LineItem.Qty, inv.LineItem.Description)

	return b
}

func (b *RecurringChargeRequestBuilder) Build() (*xrp.CreateParams, error) {
	// TODO validate the request
	return b.request, nil
}
