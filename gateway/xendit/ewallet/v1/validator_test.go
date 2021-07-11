package ewallet_test

import (
  "testing"

  "github.com/stretchr/testify/assert"

  "github.com/imrenagi/go-payment/gateway/xendit/ewallet/v1"
)

func TestOvoLegacyPhoneValidator_IsValid(t *testing.T) {

  tests := []struct {
    name    string
    phone   string
    isValid bool
  }{
    {
      name:    "valid phone number with 08xxx",
      phone:   "08111231234",
      isValid: true,
    },
    {
      name:    "invalid phone number even if it is using 08xxx",
      phone:   "0-811-123-1234",
      isValid: false,
    },
    {
      name:    "invalid phone number",
      phone:   "+628111231234",
      isValid: false,
    },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      o := ewallet.OvoPhoneValidator
      assert.Equal(t, tt.isValid, o.IsValid(tt.phone))
    })
  }
}
