package payout

import (
	"context"
	"fmt"

	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/utils/validator"
)

// Client is the client used to invoke invoice API.
type Client struct {
	Opt          *xendit.Option
	APIRequester xendit.APIRequester
}

// Create creates new payout
func (c *Client) Create(data *CreateParams) (*xendit.Payout, *xendit.Error) {
	return c.CreateWithContext(context.Background(), data)
}

// CreateWithContext creates new payout with context
func (c *Client) CreateWithContext(ctx context.Context, data *CreateParams) (*xendit.Payout, *xendit.Error) {
	if err := validator.ValidateRequired(ctx, data); err != nil {
		return nil, validator.APIValidatorErr(err)
	}

	response := &xendit.Payout{}

	err := c.APIRequester.Call(
		ctx,
		"POST",
		fmt.Sprintf("%s/payouts", c.Opt.XenditURL),
		c.Opt.SecretKey,
		nil,
		data,
		response,
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Get gets one payout
func (c *Client) Get(data *GetParams) (*xendit.Payout, *xendit.Error) {
	return c.GetWithContext(context.Background(), data)
}

// GetWithContext gets one payout with context
func (c *Client) GetWithContext(ctx context.Context, data *GetParams) (*xendit.Payout, *xendit.Error) {
	if err := validator.ValidateRequired(ctx, data); err != nil {
		return nil, validator.APIValidatorErr(err)
	}

	response := &xendit.Payout{}

	err := c.APIRequester.Call(
		ctx,
		"GET",
		fmt.Sprintf("%s/payouts/%s", c.Opt.XenditURL, data.ID),
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

// Void voids the created payout
func (c *Client) Void(data *VoidParams) (*xendit.Payout, *xendit.Error) {
	return c.VoidWithContext(context.Background(), data)
}

// VoidWithContext voids the created payout with context
func (c *Client) VoidWithContext(ctx context.Context, data *VoidParams) (*xendit.Payout, *xendit.Error) {
	if err := validator.ValidateRequired(ctx, data); err != nil {
		return nil, validator.APIValidatorErr(err)
	}

	response := &xendit.Payout{}

	err := c.APIRequester.Call(
		ctx,
		"POST",
		fmt.Sprintf("%s/payouts/%s/void", c.Opt.XenditURL, data.ID),
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
