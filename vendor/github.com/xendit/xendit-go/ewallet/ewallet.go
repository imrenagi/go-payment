package ewallet

import (
	"context"

	"github.com/xendit/xendit-go"
)

// CreatePayment creates new payment
func CreatePayment(data *CreatePaymentParams) (*xendit.EWallet, *xendit.Error) {
	return CreatePaymentWithContext(context.Background(), data)
}

// CreatePaymentWithContext creates new payment
func CreatePaymentWithContext(ctx context.Context, data *CreatePaymentParams) (*xendit.EWallet, *xendit.Error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.CreatePaymentWithContext(ctx, data)
}

// GetPaymentStatus gets one payment with its status
func GetPaymentStatus(data *GetPaymentStatusParams) (*xendit.EWallet, *xendit.Error) {
	return GetPaymentStatusWithContext(context.Background(), data)
}

// GetPaymentStatusWithContext gets one payment with its status
func GetPaymentStatusWithContext(ctx context.Context, data *GetPaymentStatusParams) (*xendit.EWallet, *xendit.Error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.GetPaymentStatusWithContext(ctx, data)
}

func getClient() (*Client, *xendit.Error) {
	return &Client{
		Opt:          &xendit.Opt,
		APIRequester: xendit.GetAPIRequester(),
	}, nil
}
