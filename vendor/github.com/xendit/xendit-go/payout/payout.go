package payout

import (
	"context"

	"github.com/xendit/xendit-go"
)

// Create creates new payout
func Create(data *CreateParams) (*xendit.Payout, *xendit.Error) {
	return CreateWithContext(context.Background(), data)
}

// CreateWithContext creates new payout with context
func CreateWithContext(ctx context.Context, data *CreateParams) (*xendit.Payout, *xendit.Error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.CreateWithContext(ctx, data)
}

// Get gets one payout
func Get(data *GetParams) (*xendit.Payout, *xendit.Error) {
	return GetWithContext(context.Background(), data)
}

// GetWithContext gets one payout with context
func GetWithContext(ctx context.Context, data *GetParams) (*xendit.Payout, *xendit.Error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.GetWithContext(ctx, data)
}

// Void voids the created payout
func Void(data *VoidParams) (*xendit.Payout, *xendit.Error) {
	return VoidWithContext(context.Background(), data)
}

// VoidWithContext voids the created payout with context
func VoidWithContext(ctx context.Context, data *VoidParams) (*xendit.Payout, *xendit.Error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.VoidWithContext(ctx, data)
}

func getClient() (*Client, *xendit.Error) {
	return &Client{
		Opt:          &xendit.Opt,
		APIRequester: xendit.GetAPIRequester(),
	}, nil
}
