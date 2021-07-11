package ewallet_test

import (
  "os"
  "testing"

  "github.com/stretchr/testify/assert"
  "github.com/xendit/xendit-go"
  "github.com/xendit/xendit-go/ewallet"

  . "github.com/imrenagi/go-payment/gateway/xendit/ewallet/v2"
  "github.com/imrenagi/go-payment/invoice"
)

func TestLinkAjaCharge(t *testing.T) {
  tests := []struct{
    name string
    invoice *invoice.Invoice
    req *ewallet.CreateEWalletChargeParams
    successRedirectURL string
  } {
    {
      name:               "successfully build the ewallet charge request for linkaja",
      invoice:            dummyInv,
      successRedirectURL: "http://example.com/success",
      req: &ewallet.CreateEWalletChargeParams{
        ReferenceID:       "a-random-invoice-number",
        Currency:          "IDR",
        Amount:            15000,
        CheckoutMethod:    "ONE_TIME_PAYMENT",
        ChannelCode:       "ID_LINKAJA",
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
      err := os.Setenv("LINKAJA_SUCCESS_REDIRECT_URL", tt.successRedirectURL)
      assert.NoError(t, err)

      req, err := NewLinkAja(tt.invoice)
      assert.NoError(t, err)
      assert.EqualValues(t, tt.req, req)
    })
  }
}
