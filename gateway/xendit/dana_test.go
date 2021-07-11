package xendit_test

import (
  "os"
  "testing"

  "github.com/stretchr/testify/assert"
  "github.com/xendit/xendit-go"
  "github.com/xendit/xendit-go/ewallet"

  xendit2 "github.com/imrenagi/go-payment/gateway/xendit"
  "github.com/imrenagi/go-payment/invoice"
)

func TestDanaCharge(t *testing.T) {
  tests := []struct{
    name string
    invoice *invoice.Invoice
    req *ewallet.CreateEWalletChargeParams
    successRedirectURL string
  } {
    {
      name: "successfully build the ewallet charge request for dana",
      invoice: dummyInv,
      successRedirectURL: "http://example.com/success",
      req: &ewallet.CreateEWalletChargeParams{
        ReferenceID:       "a-random-invoice-number",
        Currency:          "IDR",
        Amount:            15000,
        CheckoutMethod:    "ONE_TIME_PAYMENT",
        ChannelCode:       "ID_DANA",
        ChannelProperties: map[string]string{
          "success_redirect_url": "http://example.com/success",
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
      err := os.Setenv("DANA_SUCCESS_REDIRECT_URL", tt.successRedirectURL)
      assert.NoError(t, err)

      req, err := xendit2.NewDanaCharge(tt.invoice)
      assert.NoError(t, err)
      assert.EqualValues(t, tt.req, req)
    })
  }
}
