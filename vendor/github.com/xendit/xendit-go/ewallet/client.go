package ewallet

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/utils/validator"
)

// Client is the client used to invoke e-wallet API.
type Client struct {
	Opt          *xendit.Option
	APIRequester xendit.APIRequester
}

// getPaymentStatusResponse is e-wallet data that is contained in API response of Get Payment Status.
// It exists because the type of `Amount` in Get Payment Status json response is string,
// different from the CreatePayment
type getPaymentStatusResponse struct {
	EWalletType     xendit.EWalletTypeEnum `json:"ewallet_type"`
	ExternalID      string                 `json:"external_id"`
	Amount          float64                `json:"amount,string"`
	TransactionDate *time.Time             `json:"transaction_date,omitempty"`
	CheckoutURL     string                 `json:"checkout_url,omitempty"`
	BusinessID      string                 `json:"business_id,omitempty"`
}

func (r *getPaymentStatusResponse) toEwalletResponse() *xendit.EWallet {
	return &xendit.EWallet{
		EWalletType:     r.EWalletType,
		ExternalID:      r.ExternalID,
		Amount:          r.Amount,
		TransactionDate: r.TransactionDate,
		CheckoutURL:     r.CheckoutURL,
		BusinessID:      r.BusinessID,
	}
}

// CreatePayment creates new payment
func (c *Client) CreatePayment(data *CreatePaymentParams) (*xendit.EWallet, *xendit.Error) {
	return c.CreatePaymentWithContext(context.Background(), data)
}

// CreatePaymentWithContext creates new payment
func (c *Client) CreatePaymentWithContext(ctx context.Context, data *CreatePaymentParams) (*xendit.EWallet, *xendit.Error) {
	if err := validator.ValidateRequired(ctx, data); err != nil {
		return nil, validator.APIValidatorErr(err)
	}

	response := &xendit.EWallet{}
	header := &http.Header{}

	if data.ForUserID != "" {
		header.Add("for-user-id", data.ForUserID)
	}
	if data.XApiVersion != "" {
		header.Add("X-API-VERSION", data.XApiVersion)
	}

	err := c.APIRequester.Call(
		ctx,
		"POST",
		fmt.Sprintf("%s/ewallets", c.Opt.XenditURL),
		c.Opt.SecretKey,
		header,
		data,
		response,
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetPaymentStatus gets one payment with its status
func (c *Client) GetPaymentStatus(data *GetPaymentStatusParams) (*xendit.EWallet, *xendit.Error) {
	return c.GetPaymentStatusWithContext(context.Background(), data)
}

// GetPaymentStatusWithContext gets one payment with its status
func (c *Client) GetPaymentStatusWithContext(ctx context.Context, data *GetPaymentStatusParams) (*xendit.EWallet, *xendit.Error) {
	if err := validator.ValidateRequired(ctx, data); err != nil {
		return nil, validator.APIValidatorErr(err)
	}

	tempResponse := &getPaymentStatusResponse{}
	var queryString string

	if data != nil {
		queryString = data.QueryString()
	}

	err := c.APIRequester.Call(
		ctx,
		"GET",
		fmt.Sprintf("%s/ewallets?%s", c.Opt.XenditURL, queryString),
		c.Opt.SecretKey,
		nil,
		nil,
		tempResponse,
	)
	if err != nil {
		return nil, err
	}

	response := tempResponse.toEwalletResponse()

	return response, nil
}
