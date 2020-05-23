package disbursement

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

// Create creates new disbursement
func (c *Client) Create(data *CreateParams) (*xendit.Disbursement, *xendit.Error) {
	return c.CreateWithContext(context.Background(), data)
}

// CreateWithContext creates new disbursement with context
func (c *Client) CreateWithContext(ctx context.Context, data *CreateParams) (*xendit.Disbursement, *xendit.Error) {
	if err := validator.ValidateRequired(ctx, data); err != nil {
		return nil, validator.APIValidatorErr(err)
	}

	response := &xendit.Disbursement{}
	header := &http.Header{}

	if data.IdempotencyKey != "" {
		header.Add("X-IDEMPOTENCY-KEY", data.IdempotencyKey)
	}
	if data.ForUserID != "" {
		header.Add("for-user-id", data.ForUserID)
	}

	err := c.APIRequester.Call(
		ctx,
		"POST",
		fmt.Sprintf("%s/disbursements", c.Opt.XenditURL),
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

// GetByID gets a disbursement by id
func (c *Client) GetByID(data *GetByIDParams) (*xendit.Disbursement, *xendit.Error) {
	return c.GetByIDWithContext(context.Background(), data)
}

// GetByIDWithContext gets a disbursement by id with context
func (c *Client) GetByIDWithContext(ctx context.Context, data *GetByIDParams) (*xendit.Disbursement, *xendit.Error) {
	if err := validator.ValidateRequired(ctx, data); err != nil {
		return nil, validator.APIValidatorErr(err)
	}

	response := &xendit.Disbursement{}
	header := &http.Header{}
	if data.ForUserID != "" {
		header.Add("for-user-id", data.ForUserID)
	}

	err := c.APIRequester.Call(
		ctx,
		"GET",
		fmt.Sprintf("%s/disbursements/%s", c.Opt.XenditURL, data.DisbursementID),
		c.Opt.SecretKey,
		header,
		nil,
		response,
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetByExternalID gets a disbursement by id
func (c *Client) GetByExternalID(data *GetByExternalIDParams) ([]xendit.Disbursement, *xendit.Error) {
	return c.GetByExternalIDWithContext(context.Background(), data)
}

// GetByExternalIDWithContext gets a disbursement by id with context
func (c *Client) GetByExternalIDWithContext(ctx context.Context, data *GetByExternalIDParams) ([]xendit.Disbursement, *xendit.Error) {
	if err := validator.ValidateRequired(ctx, data); err != nil {
		return nil, validator.APIValidatorErr(err)
	}

	response := []xendit.Disbursement{}
	header := &http.Header{}
	if data.ForUserID != "" {
		header.Add("for-user-id", data.ForUserID)
	}

	err := c.APIRequester.Call(
		ctx,
		"GET",
		fmt.Sprintf("%s/disbursements?%s", c.Opt.XenditURL, data.QueryString()),
		c.Opt.SecretKey,
		header,
		nil,
		&response,
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetAvailableBanks gets available disbursement banks
func (c *Client) GetAvailableBanks() ([]xendit.DisbursementBank, *xendit.Error) {
	return c.GetAvailableBanksWithContext(context.Background())
}

// GetAvailableBanksWithContext gets available disbursement banks with context
func (c *Client) GetAvailableBanksWithContext(ctx context.Context) ([]xendit.DisbursementBank, *xendit.Error) {
	response := []xendit.DisbursementBank{}

	err := c.APIRequester.Call(
		ctx,
		"GET",
		fmt.Sprintf("%s/available_disbursements_banks", c.Opt.XenditURL),
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

// CreateBatch creates new batch disbursement
func (c *Client) CreateBatch(data *CreateBatchParams) (*xendit.BatchDisbursement, *xendit.Error) {
	return c.CreateBatchWithContext(context.Background(), data)
}

// CreateBatchWithContext creates new batch disbursement with context
func (c *Client) CreateBatchWithContext(ctx context.Context, data *CreateBatchParams) (*xendit.BatchDisbursement, *xendit.Error) {
	if err := validator.ValidateRequired(ctx, data); err != nil {
		return nil, validator.APIValidatorErr(err)
	}

	response := &xendit.BatchDisbursement{}
	header := &http.Header{}

	if data.IdempotencyKey != "" {
		header.Add("X-IDEMPOTENCY-KEY", data.IdempotencyKey)
	}
	if data.ForUserID != "" {
		header.Add("for-user-id", data.ForUserID)
	}

	err := c.APIRequester.Call(
		ctx,
		"POST",
		fmt.Sprintf("%s/batch_disbursements", c.Opt.XenditURL),
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
