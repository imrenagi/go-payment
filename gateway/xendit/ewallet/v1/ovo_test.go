package ewallet_test

import (
  "fmt"
  "testing"

  "github.com/stretchr/testify/assert"
  "github.com/xendit/xendit-go/ewallet"

  v1 "github.com/imrenagi/go-payment/gateway/xendit/ewallet/v1"
  "github.com/imrenagi/go-payment/invoice"
)

func TestNewOvo(t *testing.T) {

  tests := []struct {
    name string
    inv  *invoice.Invoice
    req  *ewallet.CreatePaymentParams
    // callbackURL string
    // redirectURL string
    wantErr error
  }{
    {
      name: "should create correct params",
      inv:  dummyInv,
      // callbackURL: "http://example.com/callback",
      // redirectURL: "http://example.com/success",
      wantErr: nil,
      req: &ewallet.CreatePaymentParams{
        XApiVersion:    "2020-02-01",
        EWalletType:    "OVO",
        ExternalID:     "a-random-invoice-number",
        Amount:         15000,
        Phone:          "08111231234",
        ExpirationDate: &fakeDueDate,
        // CallbackURL:    "http://example.com/callback",
        // RedirectURL:    "http://example.com/success",
        Items: []ewallet.Item{
          {
            ID:       "HOME",
            Name:     "random-item",
            Price:    15000,
            Quantity: 1,
          },
        },
      },
    },
    {
      name:    "should return error if phone number is not using correct format",
      inv:     incorrectPhoneDummyInv,
      wantErr: fmt.Errorf("invalid phone number. must be in 08xxxx format"),
      req:     nil,
    },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      params, err := v1.NewOVO(tt.inv)
      assert.Equal(t, tt.wantErr, err)
      assert.EqualValues(t, tt.req, params)
    })
  }

}
