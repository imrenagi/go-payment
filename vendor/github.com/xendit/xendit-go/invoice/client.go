package invoice

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

// Create creates new invoice
func (c *Client) Create(data *CreateParams) (*xendit.Invoice, *xendit.Error) {
	return c.CreateWithContext(context.Background(), data)
}

// CreateWithContext creates new invoice with context
func (c *Client) CreateWithContext(ctx context.Context, data *CreateParams) (*xendit.Invoice, *xendit.Error) {
	if err := validator.ValidateRequired(ctx, data); err != nil {
		return nil, validator.APIValidatorErr(err)
	}

	response := &xendit.Invoice{}
	header := &http.Header{}

	if data.ForUserID != "" {
		header.Add("for-user-id", data.ForUserID)
	}

	err := c.APIRequester.Call(
		ctx,
		"POST",
		fmt.Sprintf("%s/v2/invoices", c.Opt.XenditURL),
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

// Get gets one invoice
func (c *Client) Get(data *GetParams) (*xendit.Invoice, *xendit.Error) {
	return c.GetWithContext(context.Background(), data)
}

// GetWithContext gets one invoice with context
func (c *Client) GetWithContext(ctx context.Context, data *GetParams) (*xendit.Invoice, *xendit.Error) {
	if err := validator.ValidateRequired(ctx, data); err != nil {
		return nil, validator.APIValidatorErr(err)
	}

	response := &xendit.Invoice{}

	err := c.APIRequester.Call(
		ctx,
		"GET",
		fmt.Sprintf("%s/v2/invoices/%s", c.Opt.XenditURL, data.ID),
		c.Opt.SecretKey,
		nil,
		nil,
		response,
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Expire expires the created invoice
func (c *Client) Expire(data *ExpireParams) (*xendit.Invoice, *xendit.Error) {
	return c.ExpireWithContext(context.Background(), data)
}

// ExpireWithContext expires the created invoice with context
func (c *Client) ExpireWithContext(ctx context.Context, data *ExpireParams) (*xendit.Invoice, *xendit.Error) {
	if err := validator.ValidateRequired(ctx, data); err != nil {
		return nil, validator.APIValidatorErr(err)
	}

	response := &xendit.Invoice{}

	err := c.APIRequester.Call(
		ctx,
		"POST",
		fmt.Sprintf("%s/invoices/%s/expire!", c.Opt.XenditURL, data.ID),
		c.Opt.SecretKey,
		nil,
		nil,
		response,
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetAll gets all invoices with conditions
func (c *Client) GetAll(data *GetAllParams) ([]xendit.Invoice, *xendit.Error) {
	return c.GetAllWithContext(context.Background(), data)
}

// GetAllWithContext gets all invoices with conditions
func (c *Client) GetAllWithContext(ctx context.Context, data *GetAllParams) ([]xendit.Invoice, *xendit.Error) {
	response := []xendit.Invoice{}
	var queryString string

	if data != nil {
		queryString = data.QueryString()
	}

	err := c.APIRequester.Call(
		ctx,
		"GET",
		fmt.Sprintf("%s/v2/invoices?%s", c.Opt.XenditURL, queryString),
		c.Opt.SecretKey,
		nil,
		nil,
		&response,
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}
