package balance

import (
	"context"

	"github.com/xendit/xendit-go"
)

// Get gets balance
func Get(data *GetParams) (*xendit.Balance, *xendit.Error) {
	return GetWithContext(context.Background(), data)
}

// GetWithContext gets balance with context
func GetWithContext(ctx context.Context, data *GetParams) (*xendit.Balance, *xendit.Error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.GetWithContext(ctx, data)
}

func getClient() (*Client, *xendit.Error) {
	return &Client{
		Opt:          &xendit.Opt,
		APIRequester: xendit.GetAPIRequester(),
	}, nil
}
