package cardlesscredit

import (
	"context"

	"github.com/xendit/xendit-go"
)

// CreatePayment creates new payment
func CreatePayment(data *CreatePaymentParams) (*xendit.CardlessCredit, *xendit.Error) {
	return CreatePaymentWithContext(context.Background(), data)
}

// CreatePaymentWithContext creates new payment
func CreatePaymentWithContext(ctx context.Context, data *CreatePaymentParams) (*xendit.CardlessCredit, *xendit.Error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.CreatePaymentWithContext(ctx, data)
}

func getClient() (*Client, *xendit.Error) {
	return &Client{
		Opt:          &xendit.Opt,
		APIRequester: xendit.GetAPIRequester(),
	}, nil
}
