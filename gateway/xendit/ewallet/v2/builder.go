package ewallet

import (
  "fmt"

  goxendit "github.com/xendit/xendit-go"
  xenewallet "github.com/xendit/xendit-go/ewallet"

  "github.com/imrenagi/go-payment/invoice"
)

type EWalletTypeEnum string

const (
  EWalletIDOVO       EWalletTypeEnum = "ID_OVO"
  EWalletIDDana      EWalletTypeEnum = "ID_DANA"
  EWalletIDLinkAja   EWalletTypeEnum = "ID_LINKAJA"
  EwalletIDShopeePay EWalletTypeEnum = "ID_SHOPEEPAY"
)

func newEWalletChargeRequestBuilder(inv *invoice.Invoice) *eWalletChargeRequestBuilder {

  b := &eWalletChargeRequestBuilder{
    request: &xenewallet.CreateEWalletChargeParams{
      ReferenceID:    inv.Number,
      CheckoutMethod: "ONE_TIME_PAYMENT",
    },
  }

  b.setCustomerData(inv).
    setPrice(inv).
    setItems(inv)

  return b
}

type eWalletChargeRequestBuilder struct {
  request *xenewallet.CreateEWalletChargeParams
}

func (b *eWalletChargeRequestBuilder) setCustomerData(inv *invoice.Invoice) *eWalletChargeRequestBuilder {
  // b.request.CustomerID = inv.BillingAddress.Email
  return b
}

func (b *eWalletChargeRequestBuilder) setPrice(inv *invoice.Invoice) *eWalletChargeRequestBuilder {
  b.request.Amount = inv.GetTotal()
  b.request.Currency = inv.Currency
  return b
}

func (b *eWalletChargeRequestBuilder) setItems(inv *invoice.Invoice) *eWalletChargeRequestBuilder {
  if inv.LineItems == nil {
    return b
  }

  var items []goxendit.EWalletBasketItem
  for _, item := range inv.LineItems {
    items = append(items, goxendit.EWalletBasketItem{
      ReferenceID: fmt.Sprintf("%d", item.ID),
      Name:        item.Name,
      Category:    item.Category,
      Currency:    item.Currency,
      Price:       item.UnitPrice,
      Quantity:    item.Qty,
      Type:        "PRODUCT", // TODO do not hardcode this
      Description: item.Description,
    })
  }

  b.request.Basket = items
  return b
}

func (b *eWalletChargeRequestBuilder) SetPaymentMethod(m EWalletTypeEnum) *eWalletChargeRequestBuilder {
  b.request.ChannelCode = string(m)
  return b
}

func (b *eWalletChargeRequestBuilder) SetChannelProperties(props map[string]string) *eWalletChargeRequestBuilder {
  b.request.ChannelProperties = props
  return b
}

func (b *eWalletChargeRequestBuilder) Build() (*xenewallet.CreateEWalletChargeParams, error) {
  return b.request, nil
}

type ewalletRequestBuilder interface {
  Build() (*xenewallet.CreatePaymentParams, error)
}

// Deprecated: NewEWalletRequest generate legacy ewallet body request for xendit. This API is
// deprecated. Consider to use the newEWalletChargeRequestBuilder
func NewEWalletRequest(inv *invoice.Invoice) *EWalletRequestBuilder {

  b := &EWalletRequestBuilder{
    request: &xenewallet.CreatePaymentParams{
      XApiVersion: "2020-02-01",
      ExternalID:  inv.Number,
    },
  }

  return b.SetCustomerData(inv).
    SetPrice(inv).
    SetItemDetails(inv).
    SetExpiration(inv)
}

type EWalletRequestBuilder struct {
  request *xenewallet.CreatePaymentParams
}

func (b *EWalletRequestBuilder) SetItemDetails(inv *invoice.Invoice) *EWalletRequestBuilder {

  if inv.LineItems == nil {
    return b
  }

  var out []xenewallet.Item
  for _, item := range inv.LineItems {
    out = append(out, xenewallet.Item{
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

func (b *EWalletRequestBuilder) Build() (*xenewallet.CreatePaymentParams, error) {
  // TODO validate the request
  // phone number for ovo must be 08xxxxx format only for ovo
  return b.request, nil
}
