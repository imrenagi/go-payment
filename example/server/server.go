package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gorm.io/gorm"

	"github.com/rs/zerolog/log"

	"github.com/imrenagi/go-payment/datastore/inmemory"
	dssql "github.com/imrenagi/go-payment/datastore/sql"
	"github.com/imrenagi/go-payment/gateway/midtrans"
	"github.com/imrenagi/go-payment/invoice"
	"github.com/imrenagi/go-payment/manage"
	"github.com/imrenagi/go-payment/server"
	"github.com/imrenagi/go-payment/subscription"
	"github.com/imrenagi/go-payment/util/localconfig"
)

func main() {
	config, err := localconfig.LoadConfig("./config.yaml")
	if err != nil {
		panic(err)
	}

	secret, err := localconfig.LoadSecret("./secret.yaml")
	if err != nil {
		panic(err)
	}

	var dbDriver gorm.Dialector
	if secret.DB.Driver == "mysql" {
		access := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", secret.DB.UserName, secret.DB.Password, secret.DB.Host, secret.DB.Port, secret.DB.DBName)
		dbDriver = mysql.Open(access)
	} else if secret.DB.Driver == "postgres" {
		dbDriver = postgres.Open(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s", secret.DB.Host, secret.DB.UserName, secret.DB.Password, secret.DB.DBName, secret.DB.Port, secret.DB.Timezone)) // DSN postgres
	} else if secret.DB.Driver == "sqlite" {
		dbDriver = sqlite.Open("./gorm.db")
	} else { // no database driver provided
		log.Fatal().Msg("Can't connect to database Server")
	}

	db, err := gorm.Open(dbDriver, &gorm.Config{})
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	db.AutoMigrate(
		&midtrans.TransactionStatus{},
		&subscription.Subscription{},
		&subscription.Schedule{},
		&invoice.Invoice{},
		&invoice.Payment{},
		&invoice.CreditCardDetail{},
		&invoice.LineItem{},
		&invoice.BillingAddress{},
	)

	m := manage.NewManager(*config, secret.Payment)
	m.MustMidtransTransactionStatusRepository(dssql.NewMidtransTransactionRepository(db))
	m.MustInvoiceRepository(dssql.NewInvoiceRepository(db))
	m.MustSubscriptionRepository(dssql.NewSubscriptionRepository(db))
	m.MustPaymentConfigReader(inmemory.NewPaymentConfigRepository("./payment-methods.yaml"))

	srv := srv{
		Router:     mux.NewRouter(),
		paymentSrv: server.NewServer(m),
	}
	srv.routes()

	if err := http.ListenAndServe(fmt.Sprintf(":%s", secret.Setting.RunningPort), srv.GetHandler()); err != nil {
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
		MaxAge:             60, // 1 minutes
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
	s.Router.HandleFunc("/payment/xendit/ewallet/callback", s.paymentSrv.XenditEWalletCallbackHandler()).Methods("POST")
	s.Router.HandleFunc("/payment/subscriptions", s.paymentSrv.CreateSubscriptionHandler()).Methods("POST")
	s.Router.HandleFunc("/payment/subscriptions/{subscription_number}/pause", s.paymentSrv.PauseSubscriptionHandler()).Methods("POST", "PUT")
	s.Router.HandleFunc("/payment/subscriptions/{subscription_number}/stop", s.paymentSrv.StopSubscriptionHandler()).Methods("POST", "PUT")
	s.Router.HandleFunc("/payment/subscriptions/{subscription_number}/resume", s.paymentSrv.ResumeSubscriptionHandler()).Methods("POST", "PUT")
}
