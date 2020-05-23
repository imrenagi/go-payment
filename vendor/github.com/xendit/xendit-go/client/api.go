// Package client provides a Xendit client for invoking APIs across all products
package client

import (
	"net/http"

	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/balance"
	"github.com/xendit/xendit-go/card"
	"github.com/xendit/xendit-go/cardlesscredit"
	"github.com/xendit/xendit-go/disbursement"
	"github.com/xendit/xendit-go/ewallet"
	"github.com/xendit/xendit-go/invoice"
	"github.com/xendit/xendit-go/payout"
	"github.com/xendit/xendit-go/recurringpayment"
	"github.com/xendit/xendit-go/retailoutlet"
	"github.com/xendit/xendit-go/virtualaccount"
)

// API is the Xendit client which contains all products
type API struct {
	opt              xendit.Option
	apiRequester     xendit.APIRequester
	Invoice          *invoice.Client
	EWallet          *ewallet.Client
	Balance          *balance.Client
	RetailOutlet     *retailoutlet.Client
	VirtualAccount   *virtualaccount.Client
	Card             *card.Client
	Payout           *payout.Client
	RecurringPayment *recurringpayment.Client
	CardlessCredit   *cardlesscredit.Client
	Disbursement     *disbursement.Client
}

func (a *API) init() {
	a.Invoice = &invoice.Client{Opt: &a.opt, APIRequester: a.apiRequester}
	a.EWallet = &ewallet.Client{Opt: &a.opt, APIRequester: a.apiRequester}
	a.Balance = &balance.Client{Opt: &a.opt, APIRequester: a.apiRequester}
	a.RetailOutlet = &retailoutlet.Client{Opt: &a.opt, APIRequester: a.apiRequester}
	a.VirtualAccount = &virtualaccount.Client{Opt: &a.opt, APIRequester: a.apiRequester}
	a.Card = &card.Client{Opt: &a.opt, APIRequester: a.apiRequester}
	a.Payout = &payout.Client{Opt: &a.opt, APIRequester: a.apiRequester}
	a.RecurringPayment = &recurringpayment.Client{Opt: &a.opt, APIRequester: a.apiRequester}
	a.CardlessCredit = &cardlesscredit.Client{Opt: &a.opt, APIRequester: a.apiRequester}
	a.Disbursement = &disbursement.Client{Opt: &a.opt, APIRequester: a.apiRequester}
}

// New creates a new Xendit API client
func New(secretKey string) *API {
	api := API{
		opt: xendit.Option{
			SecretKey: secretKey,
			XenditURL: "https://api.xendit.co",
		},
		apiRequester: xendit.GetAPIRequester(),
	}
	api.init()

	return &api
}

// WithAPIRequester set custom APIRequester for Xendit Client
// Can be chained with constructor like below:
// 		client.New(yourSecretKey).WithAPIRequester(yourCustomRequester)
func (a *API) WithAPIRequester(apiRequester xendit.APIRequester) *API {
	a.apiRequester = apiRequester
	a.init()
	return a
}

// WithCustomURL set custom xendit URL for Xendit Client
// Can be chained with constructor like below:
// 		client.New(yourSecretKey).WithCustomURL(yourCustomURL)
func (a *API) WithCustomURL(xenditURL string) *API {
	a.opt.XenditURL = xenditURL
	a.init()
	return a
}

// WithCustomHTTPClient set custom HTTP Client for default API Requester
// Can be chained with constructor like below:
// 		client.New(yourSecretKey).WithCustomHTTPClient(yourCustomHTTPClient)
func (a *API) WithCustomHTTPClient(client *http.Client) *API {
	a.apiRequester = &xendit.APIRequesterImplementation{
		HTTPClient: client,
	}
	a.init()
	return a
}
