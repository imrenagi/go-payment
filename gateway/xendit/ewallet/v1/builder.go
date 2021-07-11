package ewallet

import (
  goxendit "github.com/xendit/xendit-go"
  "github.com/xendit/xendit-go/ewallet"

  "github.com/imrenagi/go-payment/invoice"
)

// Deprecated: mewEWalletRequest generate legacy ewallet body request for xendit. This API is
// deprecated. Consider to use the newEWalletChargeRequestBuilder
func mewEWalletRequest(inv *invoice.Invoice) *eWalletRequestBuilder {

  b := &eWalletRequestBuilder{
    request: &ewallet.CreatePaymentParams{
      XApiVersion: "2020-02-01",
      ExternalID:  inv.Number,
    },
  }

  return b.SetCustomerData(inv).
    SetPrice(inv).
    SetItemDetails(inv).
    SetExpiration(inv)
}

type eWalletRequestBuilder struct {
  request *ewallet.CreatePaymentParams
}

func (b *eWalletRequestBuilder) SetItemDetails(inv *invoice.Invoice) *eWalletRequestBuilder {

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

func (b *eWalletRequestBuilder) SetExpiration(inv *invoice.Invoice) *eWalletRequestBuilder {
  b.request.ExpirationDate = &inv.DueDate
  return b
}

func (b *eWalletRequestBuilder) SetCustomerData(inv *invoice.Invoice) *eWalletRequestBuilder {
  b.request.Phone = inv.BillingAddress.PhoneNumber
  return b
}

func (b *eWalletRequestBuilder) SetPrice(inv *invoice.Invoice) *eWalletRequestBuilder {
  b.request.Amount = inv.GetTotal()
  return b
}

func (b *eWalletRequestBuilder) SetPaymentMethod(m goxendit.EWalletTypeEnum) *eWalletRequestBuilder {
  b.request.EWalletType = m
  return b
}

func (b *eWalletRequestBuilder) SetCallback(url string) *eWalletRequestBuilder {
  b.request.CallbackURL = url
  return b
}

func (b *eWalletRequestBuilder) SetRedirect(url string) *eWalletRequestBuilder {
  b.request.RedirectURL = url
  return b
}

func (b *eWalletRequestBuilder) Build() (*ewallet.CreatePaymentParams, error) {
  // TODO validate the request
  // phone number for ovo must be 08xxxxx format only for ovo
  return b.request, nil
}
