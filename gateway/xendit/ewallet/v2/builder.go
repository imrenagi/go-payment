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

func newBuilder(inv *invoice.Invoice) *builer {

  b := &builer{
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

type builer struct {
  request *xenewallet.CreateEWalletChargeParams
}

func (b *builer) setCustomerData(inv *invoice.Invoice) *builer {
  // b.request.CustomerID = inv.BillingAddress.Email
  return b
}

func (b *builer) setPrice(inv *invoice.Invoice) *builer {
  b.request.Amount = inv.GetTotal()
  b.request.Currency = inv.Currency
  return b
}

func (b *builer) setItems(inv *invoice.Invoice) *builer {
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

func (b *builer) SetPaymentMethod(m EWalletTypeEnum) *builer {
  b.request.ChannelCode = string(m)
  return b
}

func (b *builer) SetChannelProperties(props map[string]string) *builer {
  b.request.ChannelProperties = props
  return b
}

func (b *builer) Build() (*xenewallet.CreateEWalletChargeParams, error) {
  return b.request, nil
}
