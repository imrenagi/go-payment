package disbursement

import (
	"context"

	"github.com/xendit/xendit-go"
)

// Create creates new disbursement
func Create(data *CreateParams) (*xendit.Disbursement, *xendit.Error) {
	return CreateWithContext(context.Background(), data)
}

// CreateWithContext creates new disbursement with context
func CreateWithContext(ctx context.Context, data *CreateParams) (*xendit.Disbursement, *xendit.Error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.CreateWithContext(ctx, data)
}

// GetByID gets a disbursement
func GetByID(data *GetByIDParams) (*xendit.Disbursement, *xendit.Error) {
	return GetByIDWithContext(context.Background(), data)
}

// GetByIDWithContext gets a disbursement with context
func GetByIDWithContext(ctx context.Context, data *GetByIDParams) (*xendit.Disbursement, *xendit.Error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.GetByIDWithContext(ctx, data)
}

// GetByExternalID gets a disbursement
func GetByExternalID(data *GetByExternalIDParams) ([]xendit.Disbursement, *xendit.Error) {
	return GetByExternalIDWithContext(context.Background(), data)
}

// GetByExternalIDWithContext gets a disbursement with context
func GetByExternalIDWithContext(ctx context.Context, data *GetByExternalIDParams) ([]xendit.Disbursement, *xendit.Error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.GetByExternalIDWithContext(ctx, data)
}

// GetAvailableBanks gets available disbursement banks
func GetAvailableBanks() ([]xendit.DisbursementBank, *xendit.Error) {
	return GetAvailableBanksWithContext(context.Background())
}

// GetAvailableBanksWithContext gets available disbursement banks with context
func GetAvailableBanksWithContext(ctx context.Context) ([]xendit.DisbursementBank, *xendit.Error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.GetAvailableBanksWithContext(ctx)
}

// CreateBatch creates new batch disbursement
func CreateBatch(data *CreateBatchParams) (*xendit.BatchDisbursement, *xendit.Error) {
	return CreateBatchWithContext(context.Background(), data)
}

// CreateBatchWithContext creates new batch disbursement with context
func CreateBatchWithContext(ctx context.Context, data *CreateBatchParams) (*xendit.BatchDisbursement, *xendit.Error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.CreateBatchWithContext(ctx, data)
}

func getClient() (*Client, *xendit.Error) {
	return &Client{
		Opt:          &xendit.Opt,
		APIRequester: xendit.GetAPIRequester(),
	}, nil
}
