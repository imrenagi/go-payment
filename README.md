Go-Payment
===

Payment module used as proxy for multiple payment gateways. Currently it only supports [Midtrans SNAP](https://snap-docs.midtrans.com/) and Xendit [Ewallet](https://xendit.github.io/apireference/#ewallets) and [XenInvoice](https://xendit.github.io/apireference/#invoices). Support for other channels will be added incrementally.

> This payment proxy is a payment service I used for my [personal site](https://imrenagi.com/donate). Thinking that this might be useful to help other people so that they can start accept money ASAP.

## Features

### Payment Channels Supported
In general, this payment proxy can support payment through this following channels:
* Credit card payment with/without installment 
* Ewallet (GoPay, OVO, Dana, LinkAja)
* Retail Outlet (Alfamart, Alfamidi, Dan+Dan)
* Cardless Credit (Akulaku)
* Bank Transfer via Virtual Account (BCA, BNI, Mandiri, Permata, Other Bank). BRI channel is coming.

> :heavy_exclamation_mark: Support for recurring payment will be added soon!

### Why you should use this payment proxy?
* If you are planning to use Midtrans SNAP and Xendit Invoice as the UI for the payment, you are strongly encouraged to use this proxy because it supports both UIs.
* This proxy helps you managing the payment gateway used for each channel. It internally connects to both payment gateway as you need, in no time. What your API user knows is only one single API to generate `Invoice`
* This proxy helps you seemlesly switch the gateway for a payment channel whenever one of them is not functioning properly/down for maintenance. For instance, Bank Transfer by VA, are supported by Midtrans and Xendit. If Midtrans VA is going south, you can easily switch the gateway to Xendit simply by updating the configuration file.
* You can choose whether to absorb the admin/installment fees by yourself or to off load it to your user by changing the payment configuration written in yaml.
* This proxy can generate `Invoice` storing informations about the customer info, item, payment method selected, and its state. `Invoice` state will change over the time depends on the transaction status callback sent by payment gateway.
* You can opt-in to store payment notification callback to your database. Currently it only stores midtrans transaction status. Support for xendit will be added soon.

### Current Limitations

1. This proxy only support 1 item per invoice. Support for multiple items per invoice is not priority as of now.
1. For simplify the query creation for database join, I use [gorm.io](https://gorm.io/) as the ORM library. The repository interfaces are provided indeed. However, default implementations with [gorm.io](https://gorm.io/) for several entities are provided in `datastore/mysql` package.
1. This proxy is not made for supporting all use cases available out there. It's hard requirement is just so that people can accept payment with as low effort as possible without need to worry about custom UI flow.
1. No callback trigger at least of now once the payment manager is done procesing this request. This will be the next priority of the next release.

### Implemented Channels

This tables shows which payment channels that has been implemented by this proxy.

:white_check_mark: : ready

:heavy_exclamation_mark: : in progress

:x: : not yet supported natively by payment gateway

|  Channels | Midtrans | Xendit |
|---|---|---|
| Credit Card without installment  | :white_check_mark: | :heavy_exclamation_mark: |
| Credit Card with installment | :white_check_mark: | :x:  |
| BCA VA  | :white_check_mark:  | :heavy_exclamation_mark:  |
| Mandiri VA  | :white_check_mark:  | :heavy_exclamation_mark:  |
| BNI VA  | :white_check_mark:  | :heavy_exclamation_mark:  |
| Permata VA  | :white_check_mark:  | :heavy_exclamation_mark:  |
| Other VA  | :white_check_mark:  | :heavy_exclamation_mark:  |
| BRI VA  | :x:  | :heavy_exclamation_mark:  |
| Alfamart, Alfamidi, Dan+Dan  | :white_check_mark:  | :heavy_exclamation_mark: |
| QRIS | :white_check_mark: via Gopay Option | :x: |
| Gopay | :white_check_mark: | :x: |
| OVO | :x: | :white_check_mark: |
| DANA | :x: | :white_check_mark: |
| LinkAja | :x: | :white_check_mark: |
| Akulaku |  :white_check_mark: |  :x: |
| Kredivo | :x: | :heavy_exclamation_mark: |


## Getting Started

### Payment Gateway Registration

#### Midtrans

#### Xendit

### Payment Gateway Callback

#### Midtrans

#### Xendit

### Application Secret

```yaml
db:
  host: "127.0.0.1"
  port: 3306
  username: "imrenagi"
  password: "imrenagi"
  dbname: "imrenagi"
payment:
  midtrans:
    secretKey: "midtrans-server-secret"
    clientKey: "midtrans-client-key"
    clientId: "midtrans-client-id"
  xendit:
    secretKey: "xendit-api-key"
    callbackToken: "xendit-callback-token"
```

### Configuration File

You can take a look sample configuration file named [payment-methods.yml](/example/server/payment-methods.yml).

```yaml
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
```
With above configuration, for installment `offline` with `bca`, you can apply this following fees to the invoice after user generates new invoice:
1. 2.9% + IDR 2000 admin fee for credit card transaction without any installment, or
1. 5.5% + IDR 2200 installment fee for credit card transaction with installment for 3 month tenure.

If you want to absorb the fee, you can simply set `val_percentage` and `val_currency` as `0`

If you only want to apply fee just either by using `val_pecentage` or `val_currency`, simply set the value to one of them and give `0` to the other. For instance:
```yaml
bank_transfers:
  - gateway: midtrans
    payment_type: "bca_va"
    display_name: "BCA"
    admin_fee:
      IDR:
        val_percentage: 0
        val_currency: 4000
        currency: "IDR"
```

### Mandatory Environment Variables

```bash
# ENVIRONMENT can be either staging or production
export ENVIRONMENT=staging
export LOG_LEVEL=DEBUG
export SERVER_BASE_URL="http://localhost:8080"
export WEB_BASE_URL="https://imrenagi.com"
export SUCCESS_REDIRECT_PATH="/donate/thanks"
export PENDING_REDIRECT_PATH="/donate/pending"
export FAILED_REDIRECT_PATH="/donate/error"
```

## Example Code

To start using this module, you can try the example [server.go](/example/server/server.go)

```go
package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/imrenagi/go-payment/datastore/inmemory"
	dsmysql "github.com/imrenagi/go-payment/datastore/mysql"
	"github.com/imrenagi/go-payment/gateway/midtrans"
	"github.com/imrenagi/go-payment/invoice"
	"github.com/imrenagi/go-payment/manage"
	"github.com/imrenagi/go-payment/server"
	"github.com/imrenagi/go-payment/util/db/mysql"
	"github.com/imrenagi/go-payment/util/localconfig"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
)

func main() {

	secret, err := localconfig.LoadSecret("example/server/secret.yaml")
	if err != nil {
		panic(err)
	}

	db := mysql.NewGorm(secret.DB)
	db.AutoMigrate(
		&midtrans.TransactionStatus{},
		&invoice.Invoice{},
		&invoice.Payment{},
		&invoice.CreditCardDetail{},
		&invoice.LineItem{},
		&invoice.BillingAddress{},
	)

	m := manage.NewManager(secret.Payment)
	m.MustMidtransTransactionStatusRepository(dsmysql.NewMidtransTransactionRepository(db))
	m.MustInvoiceRepository(dsmysql.NewInvoiceRepository(db))
	m.MustPaymentConfigReader(inmemory.NewPaymentConfigRepository("example/server/payment-methods.yml"))

	srv := srv{
		Router:     mux.NewRouter(),
		paymentSrv: server.NewServer(m),
	}
	srv.routes()

	if err := http.ListenAndServe(":8080", srv.GetHandler()); err != nil {
		log.Fatal().Msgf("Server can't run. Got: `%v`", err)
	}

}
```

To run the application, simply use:
```bash
$ go run example/server/server.go
```

> :heavy_exclamation_mark: If you want to accept payment callback from the payment gateway, consider to use [ngrok.io](https://ngrok.io) to expose your localhost to the internet and update the callback base URL in payment gateway dashboard accordingly.

## Contributing



