package xeninvoice

import (
	"fmt"
	"os"
	"strings"

	xinvoice "github.com/xendit/xendit-go/invoice"

	"github.com/imrenagi/go-payment/invoice"
)

func newBuilder(inv *invoice.Invoice) *builer {
	var shouldSendEmail bool = true

	successRedirectURL := os.Getenv("INVOICE_SUCCESS_REDIRECT_URL")
	if inv.SuccessRedirectURL != "" {
		successRedirectURL = inv.SuccessRedirectURL
	}

	failureRedirectURL := os.Getenv("INVOICE_FAILED_REDIRECT_URL")
	if inv.FailureRedirectURL != "" {
		failureRedirectURL = inv.FailureRedirectURL
	}

	b := &builer{
		request: &xinvoice.CreateParams{
			ExternalID:         inv.Number,
			ShouldSendEmail:    &shouldSendEmail,
			SuccessRedirectURL: successRedirectURL,
			FailureRedirectURL: failureRedirectURL,
			Currency:           "IDR",
			PaymentMethods:     make([]string, 0),
		},
	}

	return b.SetPrice(inv).
		SetCustomerData(inv).
		SetItemDetails(inv).
		SetExpiration(inv)
}

type builer struct {
	request *xinvoice.CreateParams
}

func (b *builer) SetPrice(inv *invoice.Invoice) *builer {
	b.request.Amount = inv.GetTotal()
	return b
}

func (b *builer) SetCustomerData(inv *invoice.Invoice) *builer {
	b.request.PayerEmail = inv.BillingAddress.Email
	return b
}

func (b *builer) SetItemDetails(inv *invoice.Invoice) *builer {

	if inv.LineItems == nil || len(inv.LineItems) == 0 {
		return b
	}

	var sb strings.Builder
	for _, item := range inv.LineItems {
		fmt.Fprintf(&sb, "")
		fmt.Fprintf(&sb, "%s %s (%d package). \n", item.Name, item.Description, item.Qty)
	}

	b.request.Description = sb.String()

	return b
}

func (b *builer) SetExpiration(inv *invoice.Invoice) *builer {
	b.request.InvoiceDuration = int(inv.DueDate.Sub(inv.InvoiceDate).Seconds())
	return b
}

func (b *builer) AddPaymentMethod(m string) *builer {
	switch strings.ToUpper(m) {
	case "BCA",
		"BRI",
		"MANDIRI",
		"BNI",
		"PERMATA",
		"ALFAMART",
		"CREDIT_CARD",
		"DANA",
		"LINKAJA",
		"OVO":
		b.request.PaymentMethods = append(b.request.PaymentMethods, m)
	}
	return b
}

func (b *builer) Build() (*xinvoice.CreateParams, error) {
	// TODO validate the request
	return b.request, nil
}
