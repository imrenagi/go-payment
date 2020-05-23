package invoice

import (
	"context"

	"github.com/xendit/xendit-go"
)

// Create creates new invoice
func Create(data *CreateParams) (*xendit.Invoice, *xendit.Error) {
	return CreateWithContext(context.Background(), data)
}

// CreateWithContext creates new invoice with context
func CreateWithContext(ctx context.Context, data *CreateParams) (*xendit.Invoice, *xendit.Error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.CreateWithContext(ctx, data)
}

// Get gets one invoice
func Get(data *GetParams) (*xendit.Invoice, *xendit.Error) {
	return GetWithContext(context.Background(), data)
}

// GetWithContext gets one invoice with context
func GetWithContext(ctx context.Context, data *GetParams) (*xendit.Invoice, *xendit.Error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.GetWithContext(ctx, data)
}

// Expire expires the created invoice
func Expire(data *ExpireParams) (*xendit.Invoice, *xendit.Error) {
	return ExpireWithContext(context.Background(), data)
}

// ExpireWithContext expires the created invoice with context
func ExpireWithContext(ctx context.Context, data *ExpireParams) (*xendit.Invoice, *xendit.Error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.ExpireWithContext(ctx, data)
}

// GetAll gets all invoices with conditions
func GetAll(data *GetAllParams) ([]xendit.Invoice, *xendit.Error) {
	return GetAllWithContext(context.Background(), data)
}

// GetAllWithContext gets all invoices with conditions
func GetAllWithContext(ctx context.Context, data *GetAllParams) ([]xendit.Invoice, *xendit.Error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.GetAllWithContext(ctx, data)
}

func getClient() (*Client, *xendit.Error) {
	return &Client{
		Opt:          &xendit.Opt,
		APIRequester: xendit.GetAPIRequester(),
	}, nil
}
