package ewallet_test

import (
  "testing"

  "github.com/stretchr/testify/assert"

  . "github.com/imrenagi/go-payment/gateway/xendit/ewallet/v2"
)

func TestOvoChargePhoneValidator_IsValid(t *testing.T) {

  tests := []struct {
    name    string
    phone   string
    isValid bool
  }{
    {
      name:    "valid phone number with +62",
      phone:   "+628111231234",
      isValid: true,
    },
    {
      name:    "invalid phone number even if it is using +62",
      phone:   "+62-811-123-1234",
      isValid: false,
    },
    {
      name:    "invalid phone number",
      phone:   "08111231234",
      isValid: false,
    },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      o := OvoChargePhoneValidator
      assert.Equal(t, tt.isValid, o.IsValid(tt.phone))
    })
  }
}
