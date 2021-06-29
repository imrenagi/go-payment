package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/gateway/midtrans"
	"github.com/imrenagi/go-payment/gateway/xendit"
	"github.com/imrenagi/go-payment/invoice"
	"github.com/imrenagi/go-payment/manage"
	"gorm.io/gorm"
	mgo "github.com/veritrans/go-midtrans"
)

func migrate(db *gorm.DB) {
	db.AutoMigrate(
		&midtrans.TransactionStatus{},
		&invoice.Invoice{},
		&invoice.Payment{},
		&invoice.CreditCardDetail{},
		&invoice.LineItem{},
		&invoice.BillingAddress{},
	)
}

func NewServer(m manage.Payment) *Server {
	return &Server{
		Manager: m,
	}
}

// Server payment server struct
type Server struct {
	Manager manage.Payment
}

// GetInvoiceRequestHandler returns handler func that will return invoice with given invoice number
func (s Server) GetInvoiceRequestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		invoiceNumber := vars["invoice_number"]
		inv, err := s.Manager.GetInvoice(r.Context(), invoiceNumber)
		if err != nil {
			WriteFailResponseFromError(w, err)
			return
		}
		WriteSuccessResponse(w, http.StatusOK, inv, nil)
	}
}

// GetPaymentMethodsHandler return handler func that will retrieve all
// payment methods available
func (s Server) GetPaymentMethodsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		options, err := payment.NewPaymentMethodListOptions(r)
		if err != nil {
			WriteFailResponse(w, http.StatusBadRequest, Error{StatusCode: http.StatusBadRequest, Message: err.Error()})
			return
		}
		methods, err := s.Manager.GetPaymentMethods(r.Context(), options...)
		if err != nil {
			WriteFailResponseFromError(w, err)
			return
		}
		WriteSuccessResponse(w, http.StatusOK, methods, nil)
	}
}

func (s Server) xenditCallbackToken(r *http.Request) (string, error) {
	token := r.Header.Get("X-CALLBACK-TOKEN")
	return token, nil
}

// CreateInvoiceHandler handles request for creating the invoice
func (s Server) CreateInvoiceHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req manage.GenerateInvoiceRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			WriteFailResponse(w, http.StatusBadRequest, Error{StatusCode: http.StatusBadRequest, Message: err.Error()})
			return
		}
		inv, err := s.Manager.GenerateInvoice(r.Context(), &req)
		if err != nil {
			WriteFailResponseFromError(w, err)
			return
		}
		WriteSuccessResponse(w, http.StatusOK, inv, nil)
	}
}

// CreateSubscriptionHandler handles request for creating new subscription
func (s Server) CreateSubscriptionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req manage.CreateSubscriptionRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			WriteFailResponse(w, http.StatusBadRequest, Error{StatusCode: http.StatusBadRequest, Message: err.Error()})
			return
		}
		subs, err := s.Manager.CreateSubscription(r.Context(), &req)
		if err != nil {
			WriteFailResponseFromError(w, err)
			return
		}
		WriteSuccessResponse(w, http.StatusOK, subs, nil)
	}
}

// PauseSubscriptionHandler returns handler for pausing subscription
func (s Server) PauseSubscriptionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		subscriptionNumber := vars["subscription_number"]
		subs, err := s.Manager.PauseSubscription(r.Context(), subscriptionNumber)
		if err != nil {
			WriteFailResponseFromError(w, err)
			return
		}
		WriteSuccessResponse(w, http.StatusOK, subs, nil)
	}
}

// StopSubscriptionHandler returns stop subscription handler
func (s Server) StopSubscriptionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		subscriptionNumber := vars["subscription_number"]
		subs, err := s.Manager.StopSubscription(r.Context(), subscriptionNumber)
		if err != nil {
			WriteFailResponseFromError(w, err)
			return
		}
		WriteSuccessResponse(w, http.StatusOK, subs, nil)
	}
}

// ResumeSubscriptionHandler returns resume susbcription handler
func (s Server) ResumeSubscriptionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		subscriptionNumber := vars["subscription_number"]
		subs, err := s.Manager.ResumeSubscription(r.Context(), subscriptionNumber)
		if err != nil {
			WriteFailResponseFromError(w, err)
			return
		}
		WriteSuccessResponse(w, http.StatusOK, subs, nil)
	}
}

// MidtransTransactionCallbackHandler handles incoming notification about payment status from midtrans.
func (s *Server) MidtransTransactionCallbackHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var notification mgo.Response
		err := decoder.Decode(&notification)
		if err != nil {
			WriteFailResponse(w, http.StatusBadRequest, Error{
				StatusCode: http.StatusBadRequest,
				Message:    "Request can't be parsed",
			})
			return
		}
		err = s.Manager.ProcessMidtransCallback(r.Context(), notification)
		if err != nil {
			WriteFailResponseFromError(w, err)
			return
		}
		WriteSuccessResponse(w, http.StatusOK, Empty{}, nil)
		return
	}
}

// XenditOVOCallbackHandler handles notification updates for ovo from xendit
func (s *Server) XenditOVOCallbackHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var status xendit.OVOPaymentStatus
		err := decoder.Decode(&status)
		if err != nil {
			WriteFailResponse(w, http.StatusBadRequest, Error{
				StatusCode: http.StatusBadRequest,
				Message:    "Request can't be parsed",
			})
			return
		}
		// TODO add ovo callback token to status
		err = s.Manager.ProcessOVOCallback(r.Context(), &status)
		if err != nil {
			WriteFailResponseFromError(w, err)
			return
		}
		WriteSuccessResponse(w, http.StatusOK, Empty{}, nil)
		return
	}
}

// XenditLinkAjaCallbackHandler handles incoming xendit notification about link aja
func (s *Server) XenditLinkAjaCallbackHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var status xendit.LinkAjaPaymentStatus
		err := decoder.Decode(&status)
		if err != nil {
			WriteFailResponse(w, http.StatusBadRequest, Error{
				StatusCode: http.StatusBadRequest,
				Message:    "Request can't be parsed",
			})
			return
		}
		err = s.Manager.ProcessLinkAjaCallback(r.Context(), &status)
		if err != nil {
			WriteFailResponseFromError(w, err)
			return
		}
		WriteSuccessResponse(w, http.StatusOK, Empty{}, nil)
		return
	}
}

// XenditInvoiceCallbackHandler handles incoming xendit notification about xenInvoice
func (s *Server) XenditInvoiceCallbackHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var status xendit.InvoicePaymentStatus
		err := decoder.Decode(&status)
		if err != nil {
			WriteFailResponse(w, http.StatusBadRequest, Error{
				StatusCode: http.StatusBadRequest,
				Message:    "Request can't be parsed",
			})
			return
		}

		cbToken, err := s.xenditCallbackToken(r)
		if err != nil {
			WriteFailResponse(w, http.StatusBadRequest, Error{
				StatusCode: http.StatusBadRequest,
				Message:    "cant read xendit callback token",
			})
			return
		}

		status.CallbackAuthToken = cbToken

		err = s.Manager.ProcessXenditInvoicesCallback(r.Context(), &status)
		if err != nil {
			WriteFailResponseFromError(w, err)
			return
		}
		WriteSuccessResponse(w, http.StatusOK, Empty{}, nil)
		return
	}
}

// XenditDanaCallbackHandler handles incoming xendit notification about dana
func (s *Server) XenditDanaCallbackHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var status xendit.DANAPaymentStatus
		err := decoder.Decode(&status)
		if err != nil {
			WriteFailResponse(w, http.StatusBadRequest, Error{
				StatusCode: http.StatusBadRequest,
				Message:    "Request can't be parsed",
			})
			return
		}
		err = s.Manager.ProcessDANACallback(r.Context(), &status)
		if err != nil {
			WriteFailResponseFromError(w, err)
			return
		}
		WriteSuccessResponse(w, http.StatusOK, Empty{}, nil)
		return
	}
}
