package ewallet

import (
  goxendit "github.com/xendit/xendit-go"
  "github.com/xendit/xendit-go/ewallet"

  "github.com/imrenagi/go-payment/invoice"
)

// Deprecated: newBuilder generate legacy ewallet body request for xendit. This API is
// deprecated. Consider to use the newEWalletChargeRequestBuilder
func newBuilder(inv *invoice.Invoice) *builder {

  b := &builder{
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

type builder struct {
  request *ewallet.CreatePaymentParams
}

func (b *builder) SetItemDetails(inv *invoice.Invoice) *builder {

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

func (b *builder) SetExpiration(inv *invoice.Invoice) *builder {
  b.request.ExpirationDate = &inv.DueDate
  return b
}

func (b *builder) SetCustomerData(inv *invoice.Invoice) *builder {
  b.request.Phone = inv.BillingAddress.PhoneNumber
  return b
}

func (b *builder) SetPrice(inv *invoice.Invoice) *builder {
  b.request.Amount = inv.GetTotal()
  return b
}

func (b *builder) SetPaymentMethod(m goxendit.EWalletTypeEnum) *builder {
  b.request.EWalletType = m
  return b
}

func (b *builder) SetCallback(url string) *builder {
  b.request.CallbackURL = url
  return b
}

func (b *builder) SetRedirect(url string) *builder {
  b.request.RedirectURL = url
  return b
}

func (b *builder) Build() (*ewallet.CreatePaymentParams, error) {
  // TODO validate the request
  // phone number for ovo must be 08xxxxx format only for ovo
  return b.request, nil
}
