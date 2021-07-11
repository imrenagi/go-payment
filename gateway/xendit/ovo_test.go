package xendit_test

import (
  "testing"

  "github.com/stretchr/testify/assert"
  "github.com/xendit/xendit-go"
  "github.com/xendit/xendit-go/ewallet"

  xendit2 "github.com/imrenagi/go-payment/gateway/xendit"
  "github.com/imrenagi/go-payment/invoice"
)

func TestOvoCharge(t *testing.T) {
  tests := []struct{
    name string
    invoice *invoice.Invoice
    req *ewallet.CreateEWalletChargeParams
  } {
    {
      name: "successfully build the ewallet charge request builder",
      invoice: dummyInv,
      req: &ewallet.CreateEWalletChargeParams{
        ReferenceID:       "a-random-invoice-number",
        Currency:          "IDR",
        Amount:            15000,
        CheckoutMethod:    "ONE_TIME_PAYMENT",
        ChannelCode:       "ID_OVO",
        ChannelProperties: map[string]string{
          "mobile_number": "+628111231234",
        },
        CustomerID:        "",
        Basket:            []xendit.EWalletBasketItem{
          {
            ReferenceID: "1",
            Name:        "random-item",
            Category:    "HOME",
            Currency:    "IDR",
            Price:       15000,
            Quantity:    1,
            Type:        "PRODUCT",
            Description: "just description",
          },
        },
        Metadata:          nil,
      },
    },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      req, err := xendit2.NewOVOCharge(tt.invoice)
      assert.NoError(t, err)
      assert.EqualValues(t, tt.req, req)
    })
  }
}
