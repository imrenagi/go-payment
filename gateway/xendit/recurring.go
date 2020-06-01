package xendit

import (
	"fmt"
	"os"
	"strings"
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
			SuccessRedirectURL:  fmt.Sprintf("%s%s", os.Getenv("WEB_BASE_URL"), os.Getenv("SUCCESS_REDIRECT_PATH")),
			FailureRedirectURL:  fmt.Sprintf("%s%s", os.Getenv("WEB_BASE_URL"), os.Getenv("FAILED_REDIRECT_PATH")),
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

	if inv.LineItems == nil || len(inv.LineItems) == 0 {
		return b
	}

	var sb strings.Builder
	for _, item := range inv.LineItems {
		fmt.Fprintf(&sb, "- ")
		fmt.Fprintf(&sb, "%dx %s: %s.", item.Qty, item.Name, item.Description)
	}

	b.request.Description = sb.String()
	return b
}

func (b *RecurringChargeRequestBuilder) Build() (*xrp.CreateParams, error) {
	// TODO validate the request
	return b.request, nil
}
