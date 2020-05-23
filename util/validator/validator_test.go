package validator_test

import (
	"testing"

	. "github.com/imrenagi/go-payment/util/validator"
)

func TestEmailValidator_IsValid(t *testing.T) {

	var validator Validator = EmailValidator{}

	tests := []struct {
		email   string
		isValid bool
	}{
		{email: "ç$€§/az@gmail.com", isValid: false},
		{email: "abcd@gmail_yahoo.com", isValid: false},
		{email: "abcd@gmail-yahoo.com", isValid: true},
		{email: "abcd@gmailyahoo", isValid: true},
		{email: "abcd@gmail.yahoo", isValid: true},
		{email: "example.example@gmail.com", isValid: true},
	}

	for _, v := range tests {
		t.Run("Validate email", func(t *testing.T) {
			if validator.IsValid(v.email) != v.isValid {
				t.Fail()
				t.Logf("Expect %v to be %v got %v", v.email, v.isValid, !v.isValid)
			}
		})
	}
}

func TestPhoneValidator_IsValid(t *testing.T) {

	var validator Validator = PhoneNumberValidator{}

	tests := []struct {
		phone   string
		isValid bool
	}{
		{phone: "1(234)5678901x1234", isValid: true},
		{phone: "(+351) 282 43 50 50", isValid: true},
		{phone: "90191919908", isValid: true},
		{phone: "555-8909", isValid: true},
		{phone: "001 6867684", isValid: true},
		{phone: "001 6867684x1", isValid: true},
		{phone: "1 (234) 567-8901", isValid: true},
		{phone: "1-234-567-8901 ext1234", isValid: true},
		{phone: "+62811132431", isValid: true},
		{phone: "+62811132431aaa", isValid: false},
		{phone: "+62-811-132-431", isValid: true},
		{phone: "62(751)142345", isValid: true},
		{phone: "6285274507699", isValid: true},
		{phone: "089899992834", isValid: true},
		{phone: "+6285274507699", isValid: true},
		{phone: "+085274507699", isValid: false},
		{phone: "85274507699", isValid: true},
		{phone: "8527450769999", isValid: true},
	}

	for _, v := range tests {
		t.Run("Validate phone number", func(t *testing.T) {
			if validator.IsValid(v.phone) != v.isValid {
				t.Fail()
				t.Logf("Expect %v to be %v got %v", v.phone, v.isValid, !v.isValid)
			}
		})
	}
}
