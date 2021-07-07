Example
===

To start using this module, you can try the example [server.go](/server.go)

This example is using sqlite as the datastore. If you want to change the database,
please read gorm docs for connecting to database [here](https://gorm.io/docs/connecting_to_the_database.html)
```go
package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/imrenagi/go-payment/datastore/inmemory"
	dssql "github.com/imrenagi/go-payment/datastore/sql"
	"github.com/imrenagi/go-payment/gateway/midtrans"
	"github.com/imrenagi/go-payment/invoice"
	"github.com/imrenagi/go-payment/manage"
	"github.com/imrenagi/go-payment/server"
	"github.com/imrenagi/go-payment/util/localconfig"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
)

func main() {

	secret, err := localconfig.LoadSecret("example/server/secret.yaml")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open("example/server/gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	db.AutoMigrate(
		&midtrans.TransactionStatus{},
		&invoice.Invoice{},
		&invoice.Payment{},
		&invoice.CreditCardDetail{},
		&invoice.LineItem{},
		&invoice.BillingAddress{},
	)

	m := manage.NewManager(secret.Payment)
	m.MustMidtransTransactionStatusRepository(dssql.NewMidtransTransactionRepository(db))
	m.MustInvoiceRepository(dssql.NewInvoiceRepository(db))
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

```console
$ go run example/server/server.go
```

> :heavy_exclamation_mark: If you want to accept payment callback from the payment gateway on your local computer for development purpose, consider to use [ngrok.io](https://ngrok.io) to expose your localhost to the internet and update the callback base URL in all url environment variables.