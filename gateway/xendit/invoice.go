package xendit

import (
	"fmt"
	"os"
	"strings"

	"github.com/imrenagi/go-payment/invoice"
	xinvoice "github.com/xendit/xendit-go/invoice"
)

type invoiceRequestBuilder interface {
	Build() (*xinvoice.CreateParams, error)
}

func NewInvoiceRequestBuilder(inv *invoice.Invoice) *InvoiceRequestBuilder {

	var shouldSendEmail bool = true

	b := &InvoiceRequestBuilder{
		request: &xinvoice.CreateParams{
			ExternalID:         inv.Number,
			ShouldSendEmail:    &shouldSendEmail,
			SuccessRedirectURL: fmt.Sprintf("%s%s", os.Getenv("WEB_BASE_URL"), os.Getenv("SUCCESS_REDIRECT_PATH")),
			FailureRedirectURL: fmt.Sprintf("%s%s", os.Getenv("WEB_BASE_URL"), os.Getenv("FAILED_REDIRECT_PATH")),
			Currency:           "IDR",
			PaymentMethods:     make([]string, 0),
		},
	}

	return b.SetPrice(inv).
		SetCustomerData(inv).
		SetItemDetails(inv).
		SetExpiration(inv)
}

type InvoiceRequestBuilder struct {
	request *xinvoice.CreateParams
}

func (b *InvoiceRequestBuilder) SetPrice(inv *invoice.Invoice) *InvoiceRequestBuilder {
	b.request.Amount = inv.GetTotal()
	return b
}

func (b *InvoiceRequestBuilder) SetCustomerData(inv *invoice.Invoice) *InvoiceRequestBuilder {
	b.request.PayerEmail = inv.BillingAddress.Email
	return b
}

func (b *InvoiceRequestBuilder) SetItemDetails(inv *invoice.Invoice) *InvoiceRequestBuilder {

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

func (b *InvoiceRequestBuilder) SetExpiration(inv *invoice.Invoice) *InvoiceRequestBuilder {
	b.request.InvoiceDuration = int(inv.DueDate.Sub(inv.InvoiceDate).Seconds())
	return b
}

func (b *InvoiceRequestBuilder) AddPaymentMethod(m string) *InvoiceRequestBuilder {
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

func (b *InvoiceRequestBuilder) Build() (*xinvoice.CreateParams, error) {
	// TODO validate the request
	return b.request, nil
}
