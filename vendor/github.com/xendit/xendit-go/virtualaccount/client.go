package virtualaccount

import (
	"context"
	"fmt"
	"net/http"

	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/utils/validator"
)

// Client is the client used to invoke virtual account API.
type Client struct {
	Opt          *xendit.Option
	APIRequester xendit.APIRequester
}

// CreateFixedVA creates new fixed virtual account
func (c *Client) CreateFixedVA(data *CreateFixedVAParams) (*xendit.VirtualAccount, *xendit.Error) {
	return c.CreateFixedVAWithContext(context.Background(), data)
}

// CreateFixedVAWithContext creates new fixed virtual account with context
func (c *Client) CreateFixedVAWithContext(ctx context.Context, data *CreateFixedVAParams) (*xendit.VirtualAccount, *xendit.Error) {
	if err := validator.ValidateRequired(ctx, data); err != nil {
		return nil, validator.APIValidatorErr(err)
	}

	response := &xendit.VirtualAccount{}
	header := &http.Header{}

	if data.ForUserID != "" {
		header.Add("for-user-id", data.ForUserID)
	}

	err := c.APIRequester.Call(
		ctx,
		"POST",
		fmt.Sprintf("%s/callback_virtual_accounts", c.Opt.XenditURL),
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

// GetFixedVA gets one fixed virtual account
func (c *Client) GetFixedVA(data *GetFixedVAParams) (*xendit.VirtualAccount, *xendit.Error) {
	return c.GetFixedVAWithContext(context.Background(), data)
}

// GetFixedVAWithContext gets one fixed virtual account with context
func (c *Client) GetFixedVAWithContext(ctx context.Context, data *GetFixedVAParams) (*xendit.VirtualAccount, *xendit.Error) {
	if err := validator.ValidateRequired(ctx, data); err != nil {
		return nil, validator.APIValidatorErr(err)
	}

	response := &xendit.VirtualAccount{}

	err := c.APIRequester.Call(
		ctx,
		"GET",
		fmt.Sprintf("%s/callback_virtual_accounts/%s", c.Opt.XenditURL, data.ID),
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

// UpdateFixedVA updates one fixed virtual account
func (c *Client) UpdateFixedVA(data *UpdateFixedVAParams) (*xendit.VirtualAccount, *xendit.Error) {
	return c.UpdateFixedWithContext(context.Background(), data)
}

// UpdateFixedWithContext updates one fixed virtual account with context
func (c *Client) UpdateFixedWithContext(ctx context.Context, data *UpdateFixedVAParams) (*xendit.VirtualAccount, *xendit.Error) {
	if err := validator.ValidateRequired(ctx, data); err != nil {
		return nil, validator.APIValidatorErr(err)
	}

	response := &xendit.VirtualAccount{}

	err := c.APIRequester.Call(
		ctx,
		"PATCH",
		fmt.Sprintf("%s/callback_virtual_accounts/%s", c.Opt.XenditURL, data.ID),
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

// GetAvailableBanks gets available virtual account banks
func (c *Client) GetAvailableBanks() ([]xendit.VirtualAccountBank, *xendit.Error) {
	return c.GetAvailableBanksWithContext(context.Background())
}

// GetAvailableBanksWithContext gets available virtual account banks with context
func (c *Client) GetAvailableBanksWithContext(ctx context.Context) ([]xendit.VirtualAccountBank, *xendit.Error) {
	response := []xendit.VirtualAccountBank{}

	err := c.APIRequester.Call(
		ctx,
		"GET",
		fmt.Sprintf("%s/available_virtual_account_banks", c.Opt.XenditURL),
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

// GetPayment gets one fixed virtual account payment
func (c *Client) GetPayment(data *GetPaymentParams) (*xendit.VirtualAccountPayment, *xendit.Error) {
	return c.GetPaymentWithContext(context.Background(), data)
}

// GetPaymentWithContext gets one fixed virtual account payment with context
func (c *Client) GetPaymentWithContext(ctx context.Context, data *GetPaymentParams) (*xendit.VirtualAccountPayment, *xendit.Error) {
	if err := validator.ValidateRequired(ctx, data); err != nil {
		return nil, validator.APIValidatorErr(err)
	}

	response := &xendit.VirtualAccountPayment{}

	err := c.APIRequester.Call(
		ctx,
		"GET",
		fmt.Sprintf("%s/callback_virtual_account_payments/payment_id=%s", c.Opt.XenditURL, data.PaymentID),
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
