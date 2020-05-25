package config_test

import (
	"testing"
	"time"

	. "github.com/imrenagi/go-payment/config"
	"github.com/stretchr/testify/assert"
)

var paymentCfg = []byte(`
card_payment:
  payment_type: "credit_card"
  installments:    
    - type: offline
      display_name: ""
      gateway: midtrans
      bank: bca
      channel: migs
      default: true
      active: true
      terms:
        - term: 0
          admin_fee:
            IDR:
              val_percentage: 2.9
              val_currency: 2000
              currency: "IDR"
        - term: 3
          installment_fee:
            IDR:
              val_percentage: 5.5
              val_currency: 2200
              currency: "IDR"  
        - term: 6
          installment_fee:
            IDR:
              val_percentage: 7.5
              val_currency: 2200
              currency: "IDR"
        - term: 12
          installment_fee:
            IDR:
              val_percentage: 9
              val_currency: 2200
              currency: "IDR"                              
bank_transfers:
  - gateway: midtrans
    payment_type: "bca_va"
    display_name: "BCA"
    admin_fee:
      IDR:
        val_percentage: 0
        val_currency: 4000
        currency: "IDR"
    waiting_time:
      duration: 1
      unit: day
  - gateway: midtrans
    payment_type: "echannel"
    display_name: "Mandiri Bill"
    admin_fee:
      IDR:
        val_percentage: 0
        val_currency: 4000
        currency: "IDR" 
    waiting_time:
      duration: 1
      unit: day        
  - gateway: midtrans
    payment_type: "bni_va"
    display_name: "BNI"
    admin_fee:
      IDR:
        val_percentage: 0
        val_currency: 4000
        currency: "IDR"
    waiting_time:
      duration: 1
      unit: day        
  - gateway: midtrans
    payment_type: "permata_va"
    display_name: "Bank Permata"
    admin_fee:
      IDR:
        val_percentage: 0
        val_currency: 4000
        currency: "IDR"
    waiting_time:
      duration: 1
      unit: day        
  - gateway: midtrans
    payment_type: "other_va"
    display_name: "Bank Lainnya"
    admin_fee:
      IDR:
        val_percentage: 0
        val_currency: 4000
        currency: "IDR"
    waiting_time:
      duration: 1
      unit: day                             
ewallets:
  - gateway: midtrans
    payment_type: "gopay"
    display_name: "Gopay"
    admin_fee:
      IDR:
        val_percentage: 2
        val_currency: 0
        currency: "IDR"
    waiting_time:
      duration: 15
      unit: minute           
cstores:
  - gateway: midtrans
    payment_type: alfamart
    display_name: "Alfamart"
    admin_fee: 
      IDR:
        val_percentage: 0
        val_currency: 5000
        currency: "IDR"
    waiting_time:
      duration: 1
      unit: day        
cardless_credits:
  - gateway: midtrans
    payment_type: akulaku
    display_name: "Akulaku"
    admin_fee: 
      IDR:
        val_percentage: 2.0
        val_currency: 0
        currency: "IDR"
    waiting_time:
      duration: 1
      unit: day        

`)

func TestLoadSecret(t *testing.T) {
	t.Run("Test read config", func(t *testing.T) {
		configs, _ := LoadPaymentConfigs(paymentCfg)

		assert.Equal(t, 24*time.Hour, *configs.BankTransfers[0].GetPaymentWaitingTime())
		assert.Equal(t, 15*time.Minute, *configs.EWallets[0].GetPaymentWaitingTime())

	})
}
