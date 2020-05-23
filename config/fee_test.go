package config_test

import (
	"testing"

	. "github.com/imrenagi/go-payment/config"
	"github.com/stretchr/testify/assert"
)

func TestFee(t *testing.T) {

	t.Run("Fee with percentage", func(t *testing.T) {

		cases := []struct {
			Fee           Fee
			Price         float64
			ExpectedValue float64
		}{
			{
				Fee:           Fee{PercentageVal: 5},
				Price:         10000,
				ExpectedValue: 500.00,
			},
			{
				Fee:           Fee{PercentageVal: 3.9},
				Price:         110,
				ExpectedValue: 5.00,
			},
			{
				Fee:           Fee{PercentageVal: 2.9},
				Price:         11,
				ExpectedValue: 1.00,
			},
		}

		for _, c := range cases {
			assert.Equal(t, c.ExpectedValue, c.Fee.Estimate(c.Price))
		}

	})

	t.Run("Fee with value", func(t *testing.T) {

		cases := []struct {
			Fee           Fee
			Price         float64
			ExpectedValue float64
		}{
			{
				Fee:           Fee{CurrencyVal: 100},
				Price:         10000,
				ExpectedValue: 100.00,
			},
		}

		for _, c := range cases {
			assert.Equal(t, c.ExpectedValue, c.Fee.Estimate(c.Price))
		}

	})

	t.Run("Fee with value and percentage", func(t *testing.T) {

		cases := []struct {
			Fee           Fee
			Price         float64
			ExpectedValue float64
		}{
			{
				Fee:           Fee{PercentageVal: 3.9, CurrencyVal: 100.0},
				Price:         110,
				ExpectedValue: 105.00,
			},
		}

		for _, c := range cases {
			assert.Equal(t, c.ExpectedValue, c.Fee.Estimate(c.Price))
		}

	})
}
