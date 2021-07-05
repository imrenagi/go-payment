package xendit

import (
  "fmt"

  "github.com/xendit/xendit-go"
  goxendit "github.com/xendit/xendit-go"
  "github.com/xendit/xendit-go/ewallet"

  "github.com/imrenagi/go-payment/invoice"
)

type EWalletTypeEnum string

const (
  EWalletIDOVO       EWalletTypeEnum = "ID_OVO"
  EWalletIDDana      EWalletTypeEnum = "ID_DANA"
  EWalletIDLinkAja   EWalletTypeEnum = "ID_LINKAJA"
  EwalletIDShopeePay EWalletTypeEnum = "ID_SHOPEEPAY"
)

type ewalletRequestBuilderV2 interface {
  Build() (*ewallet.CreateEWalletChargeParams, error)
}

func NewEWalletChargeRequestBuilder(inv *invoice.Invoice) *EWalletChargeRequestBuilder {

  b := &EWalletChargeRequestBuilder{
    request: &ewallet.CreateEWalletChargeParams{
      ReferenceID: inv.Number,
      CheckoutMethod: "ONE_TIME_PAYMENT", // TODO dont hardcode this
    },
  }

  b.setCustomerData(inv).
    setPrice(inv).
    setItems(inv)

  return b
}

type EWalletChargeRequestBuilder struct {
  request *ewallet.CreateEWalletChargeParams
}

func (b *EWalletChargeRequestBuilder) setCustomerData(inv *invoice.Invoice) *EWalletChargeRequestBuilder {
  // b.request.CustomerID = inv.BillingAddress.Email

  return b
}

func (b *EWalletChargeRequestBuilder) setPrice(inv *invoice.Invoice) *EWalletChargeRequestBuilder {
  b.request.Amount = inv.GetTotal()
  b.request.Currency = "IDR"
  return b
}

func (b *EWalletChargeRequestBuilder) setItems(inv *invoice.Invoice) *EWalletChargeRequestBuilder {
  if inv.LineItems == nil {
    return b
  }

  var items []xendit.EWalletBasketItem
  for _, item := range inv.LineItems {
    items = append(items, xendit.EWalletBasketItem{
      ReferenceID: fmt.Sprintf("%d", item.ID),
      Name:        item.Name,
      Category:    item.Category,
      Currency:    item.Currency,
      Price:       item.UnitPrice,
      Quantity:    item.Qty,
      Type:        "SERVICE", // TODO do not hardcode this
      Description: item.Description,
    })
  }

  b.request.Basket = items
  return b
}

func (b *EWalletChargeRequestBuilder) SetPaymentMethod(m EWalletTypeEnum) *EWalletChargeRequestBuilder {
  b.request.ChannelCode = string(m)
  return b
}

func (b *EWalletChargeRequestBuilder) SetChannelProperties(props map[string]string) *EWalletChargeRequestBuilder {
  b.request.ChannelProperties = props
  return b
}

func (b *EWalletChargeRequestBuilder) Build() (*ewallet.CreateEWalletChargeParams, error) {
  // TODO validate the request
  return b.request, nil
}


type ewalletRequestBuilder interface {
  Build() (*ewallet.CreatePaymentParams, error)
}

func NewEWalletRequest(inv *invoice.Invoice) *EWalletRequestBuilder {

  b := &EWalletRequestBuilder{
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
