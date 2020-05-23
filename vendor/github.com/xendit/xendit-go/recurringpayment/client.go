package recurringpayment

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

// Create creates new recurring payment
func (c *Client) Create(data *CreateParams) (*xendit.RecurringPayment, *xendit.Error) {
	return c.CreateWithContext(context.Background(), data)
}

// CreateWithContext creates new recurring payment with context
func (c *Client) CreateWithContext(ctx context.Context, data *CreateParams) (*xendit.RecurringPayment, *xendit.Error) {
	if err := validator.ValidateRequired(ctx, data); err != nil {
		return nil, validator.APIValidatorErr(err)
	}

	response := &xendit.RecurringPayment{}
	header := &http.Header{}

	if data.ForUserID != "" {
		header.Add("for-user-id", data.ForUserID)
	}

	err := c.APIRequester.Call(
		ctx,
		"POST",
		fmt.Sprintf("%s/recurring_payments", c.Opt.XenditURL),
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

// Get gets one recurring payment
func (c *Client) Get(data *GetParams) (*xendit.RecurringPayment, *xendit.Error) {
	return c.GetWithContext(context.Background(), data)
}

// GetWithContext gets one recurring payment with context
func (c *Client) GetWithContext(ctx context.Context, data *GetParams) (*xendit.RecurringPayment, *xendit.Error) {
	if err := validator.ValidateRequired(ctx, data); err != nil {
		return nil, validator.APIValidatorErr(err)
	}

	response := &xendit.RecurringPayment{}

	err := c.APIRequester.Call(
		ctx,
		"GET",
		fmt.Sprintf("%s/recurring_payments/%s", c.Opt.XenditURL, data.ID),
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

// Edit edits a recurring payment
func (c *Client) Edit(data *EditParams) (*xendit.RecurringPayment, *xendit.Error) {
	return c.EditWithContext(context.Background(), data)
}

// EditWithContext edits a recurring payment with context
func (c *Client) EditWithContext(ctx context.Context, data *EditParams) (*xendit.RecurringPayment, *xendit.Error) {
	if err := validator.ValidateRequired(ctx, data); err != nil {
		return nil, validator.APIValidatorErr(err)
	}

	response := &xendit.RecurringPayment{}

	err := c.APIRequester.Call(
		ctx,
		"PATCH",
		fmt.Sprintf("%s/recurring_payments/%s", c.Opt.XenditURL, data.ID),
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

// Stop stops one recurring payment
func (c *Client) Stop(data *StopParams) (*xendit.RecurringPayment, *xendit.Error) {
	return c.StopWithContext(context.Background(), data)
}

// StopWithContext stops one recurring payment with context
func (c *Client) StopWithContext(ctx context.Context, data *StopParams) (*xendit.RecurringPayment, *xendit.Error) {
	if err := validator.ValidateRequired(ctx, data); err != nil {
		return nil, validator.APIValidatorErr(err)
	}

	response := &xendit.RecurringPayment{}

	err := c.APIRequester.Call(
		ctx,
		"POST",
		fmt.Sprintf("%s/recurring_payments/%s/stop!", c.Opt.XenditURL, data.ID),
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

// Pause pauses one recurring payment
func (c *Client) Pause(data *PauseParams) (*xendit.RecurringPayment, *xendit.Error) {
	return c.PauseWithContext(context.Background(), data)
}

// PauseWithContext pauses one recurring payment with context
func (c *Client) PauseWithContext(ctx context.Context, data *PauseParams) (*xendit.RecurringPayment, *xendit.Error) {
	if err := validator.ValidateRequired(ctx, data); err != nil {
		return nil, validator.APIValidatorErr(err)
	}

	response := &xendit.RecurringPayment{}

	err := c.APIRequester.Call(
		ctx,
		"POST",
		fmt.Sprintf("%s/recurring_payments/%s/pause!", c.Opt.XenditURL, data.ID),
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

// Resume resumes one recurring payment
func (c *Client) Resume(data *ResumeParams) (*xendit.RecurringPayment, *xendit.Error) {
	return c.ResumeWithContext(context.Background(), data)
}

// ResumeWithContext resumes one recurring payment with context
func (c *Client) ResumeWithContext(ctx context.Context, data *ResumeParams) (*xendit.RecurringPayment, *xendit.Error) {
	if err := validator.ValidateRequired(ctx, data); err != nil {
		return nil, validator.APIValidatorErr(err)
	}

	response := &xendit.RecurringPayment{}

	err := c.APIRequester.Call(
		ctx,
		"POST",
		fmt.Sprintf("%s/recurring_payments/%s/resume!", c.Opt.XenditURL, data.ID),
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
