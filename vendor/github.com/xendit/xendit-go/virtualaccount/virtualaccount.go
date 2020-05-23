package virtualaccount

import (
	"context"

	"github.com/xendit/xendit-go"
)

// CreateFixedVA creates new fixed virtual account
func CreateFixedVA(data *CreateFixedVAParams) (*xendit.VirtualAccount, *xendit.Error) {
	return CreateFixedVAWithContext(context.Background(), data)
}

// CreateFixedVAWithContext creates new fixed virtual account with context
func CreateFixedVAWithContext(ctx context.Context, data *CreateFixedVAParams) (*xendit.VirtualAccount, *xendit.Error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.CreateFixedVAWithContext(ctx, data)
}

// GetFixedVA gets a fixed virtual account
func GetFixedVA(data *GetFixedVAParams) (*xendit.VirtualAccount, *xendit.Error) {
	return GetFixedVAWithContext(context.Background(), data)
}

// GetFixedVAWithContext gets a fixed virtual account with context
func GetFixedVAWithContext(ctx context.Context, data *GetFixedVAParams) (*xendit.VirtualAccount, *xendit.Error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.GetFixedVAWithContext(ctx, data)
}

// UpdateFixedVA updates a fixed virtual account
func UpdateFixedVA(data *UpdateFixedVAParams) (*xendit.VirtualAccount, *xendit.Error) {
	return UpdateFixedVAWithContext(context.Background(), data)
}

// UpdateFixedVAWithContext updates a fixed virtual account with context
func UpdateFixedVAWithContext(ctx context.Context, data *UpdateFixedVAParams) (*xendit.VirtualAccount, *xendit.Error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.UpdateFixedWithContext(ctx, data)
}

// GetAvailableBanks gets available virtual account banks
func GetAvailableBanks() ([]xendit.VirtualAccountBank, *xendit.Error) {
	return GetAvailableBanksWithContext(context.Background())
}

// GetAvailableBanksWithContext gets available virtual account banks with context
func GetAvailableBanksWithContext(ctx context.Context) ([]xendit.VirtualAccountBank, *xendit.Error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.GetAvailableBanksWithContext(ctx)
}

// GetPayment gets one fixed virtual account payment
func GetPayment(data *GetPaymentParams) (*xendit.VirtualAccountPayment, *xendit.Error) {
	return GetPaymentWithContext(context.Background(), data)
}

// GetPaymentWithContext gets one fixed virtual account payment with context
func GetPaymentWithContext(ctx context.Context, data *GetPaymentParams) (*xendit.VirtualAccountPayment, *xendit.Error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.GetPaymentWithContext(ctx, data)
}

func getClient() (*Client, *xendit.Error) {
	return &Client{
		Opt:          &xendit.Opt,
		APIRequester: xendit.GetAPIRequester(),
	}, nil
}
