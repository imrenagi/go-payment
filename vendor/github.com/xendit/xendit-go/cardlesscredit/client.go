package cardlesscredit

import (
	"context"
	"fmt"
	"net/http"

	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/utils/validator"
)

// Client is the client used to invoke invoice API.
type Client struct {
	Opt          *xendit.Option
	APIRequester xendit.APIRequester
}

// CreatePayment creates new cardless credit payment
func (c *Client) CreatePayment(data *CreatePaymentParams) (*xendit.CardlessCredit, *xendit.Error) {
	return c.CreatePaymentWithContext(context.Background(), data)
}

// CreatePaymentWithContext creates new cardless credit payment with context
func (c *Client) CreatePaymentWithContext(ctx context.Context, data *CreatePaymentParams) (*xendit.CardlessCredit, *xendit.Error) {
	if err := validator.ValidateRequired(ctx, data); err != nil {
		return nil, validator.APIValidatorErr(err)
	}

	response := &xendit.CardlessCredit{}
	header := &http.Header{}

	err := c.APIRequester.Call(
		ctx,
		"POST",
		fmt.Sprintf("%s/cardless-credit", c.Opt.XenditURL),
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
