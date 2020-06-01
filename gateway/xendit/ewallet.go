package xendit

import (
	"github.com/imrenagi/go-payment/invoice"
	goxendit "github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/ewallet"
)

type ewalletRequestBuilder interface {
	Build() (*ewallet.CreatePaymentParams, error)
}

func NewEWalletRequest(inv *invoice.Invoice) *EWalletRequestBuilder {

	b := &EWalletRequestBuilder{
		request: &ewallet.CreatePaymentParams{
			ExternalID: inv.Number,
		},
	}

	return b.SetCustomerData(inv).
		SetPrice(inv).
		SetItemDetails(inv).
		SetExpiration(inv)
}

type EWalletRequestBuilder struct {
	request *ewallet.CreatePaymentParams
}

func (b *EWalletRequestBuilder) SetItemDetails(inv *invoice.Invoice) *EWalletRequestBuilder {
	if inv.LineItems == nil {
		return b
	}

	var out []ewallet.Item
	for _, item := range inv.LineItems {
		out = append(out, ewallet.Item{
			ID:       item.Category,
			Name:     item.Name,
			Price:    item.UnitPrice,
			Quantity: item.Qty,
		})
	}

	b.request.Items = out
	return b
}

func (b *EWalletRequestBuilder) SetExpiration(inv *invoice.Invoice) *EWalletRequestBuilder {
	b.request.ExpirationDate = &inv.DueDate
	return b
}

func (b *EWalletRequestBuilder) SetCustomerData(inv *invoice.Invoice) *EWalletRequestBuilder {
	b.request.Phone = inv.BillingAddress.PhoneNumber
	return b
}

func (b *EWalletRequestBuilder) SetPrice(inv *invoice.Invoice) *EWalletRequestBuilder {
	b.request.Amount = inv.GetTotal()
	return b
}

func (b *EWalletRequestBuilder) SetPaymentMethod(m goxendit.EWalletTypeEnum) *EWalletRequestBuilder {
	b.request.EWalletType = m
	return b
}

func (b *EWalletRequestBuilder) SetCallback(url string) *EWalletRequestBuilder {
	b.request.CallbackURL = url
	return b
}

func (b *EWalletRequestBuilder) SetRedirect(url string) *EWalletRequestBuilder {
	b.request.RedirectURL = url
	return b
}

func (b *EWalletRequestBuilder) Build() (*ewallet.CreatePaymentParams, error) {
	// TODO validate the request
	return b.request, nil
}
