package main

import (
	"net/http"

	"github.com/gorilla/mux"
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

	midTrxRepo := dsmysql.NewMidtransTransactionRepository(db)
	invoiceRepo := dsmysql.NewInvoiceRepository(db)

	m := manage.NewDefaultManager(secret.Payment)
	m.MustMidtransTransactionStatusRepository(midTrxRepo)
	m.MustInvoiceRepository(invoiceRepo)

	srv := srv{
		Router:     mux.NewRouter(),
		paymentSrv: server.NewServer(m),
	}
	srv.routes()

	if err := http.ListenAndServe(":8080", srv.GetHandler()); err != nil {
		log.Fatal().Msgf("Server can't run. Got: `%v`", err)
	}

}

type srv struct {
	Router     *mux.Router
	paymentSrv *server.Server
}

// GetHandler returns http.Handler which intercepted by the cors checker.
func (s *srv) GetHandler() http.Handler {

	c := cors.New(cors.Options{
		AllowedOrigins:     []string{"http://localhost:3000", "https://localhost:3000"},
		AllowedMethods:     []string{"POST", "GET", "PUT", "DELETE", "HEAD", "OPTIONS"},
		AllowedHeaders:     []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Mode"},
		MaxAge:             60, //1 minutes
		AllowCredentials:   true,
		OptionsPassthrough: false,
		Debug:              false,
	})

	return c.Handler(s.Router)
}

func (s *srv) Healthcheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}
}

func (s srv) routes() {
	s.Router.HandleFunc("/payment/methods", s.paymentSrv.GetPaymentMethodsHandler()).Methods("GET")
	s.Router.HandleFunc("/payment/invoices", s.paymentSrv.CreateInvoiceHandler()).Methods("POST")
	s.Router.HandleFunc("/payment/midtrans/callback", s.paymentSrv.MidtransTransactionCallbackHandler()).Methods("POST")
	s.Router.HandleFunc("/payment/xendit/invoice/callback", s.paymentSrv.XenditInvoiceCallbackHandler()).Methods("POST")
	s.Router.HandleFunc("/payment/xendit/ovo/callback", s.paymentSrv.XenditOVOCallbackHandler()).Methods("POST")
	s.Router.HandleFunc("/payment/xendit/dana/callback", s.paymentSrv.XenditDanaCallbackHandler()).Methods("POST")
	s.Router.HandleFunc("/payment/xendit/linkaja/callback", s.paymentSrv.XenditLinkAjaCallbackHandler()).Methods("POST")
}
