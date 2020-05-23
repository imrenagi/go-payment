package card

import (
	"context"

	"github.com/xendit/xendit-go"
)

/* Charge */

// CreateCharge creates new card charge
func CreateCharge(data *CreateChargeParams) (*xendit.CardCharge, *xendit.Error) {
	return CreateChargeWithContext(context.Background(), data)
}

// CreateChargeWithContext creates new card charge with context
func CreateChargeWithContext(ctx context.Context, data *CreateChargeParams) (*xendit.CardCharge, *xendit.Error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.CreateChargeWithContext(ctx, data)
}

// CaptureCharge captures a card charge
func CaptureCharge(data *CaptureChargeParams) (*xendit.CardCharge, *xendit.Error) {
	return CaptureChargeWithContext(context.Background(), data)
}

// CaptureChargeWithContext captures a card charge with context
func CaptureChargeWithContext(ctx context.Context, data *CaptureChargeParams) (*xendit.CardCharge, *xendit.Error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.CaptureChargeWithContext(ctx, data)
}

// GetCharge gets a card charge
func GetCharge(data *GetChargeParams) (*xendit.CardCharge, *xendit.Error) {
	return GetChargeWithContext(context.Background(), data)
}

// GetChargeWithContext gets a card charge with context
func GetChargeWithContext(ctx context.Context, data *GetChargeParams) (*xendit.CardCharge, *xendit.Error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.GetChargeWithContext(ctx, data)
}

// CreateRefund gets a card charge
func CreateRefund(data *CreateRefundParams) (*xendit.CardRefund, *xendit.Error) {
	return CreateRefundWithContext(context.Background(), data)
}

// CreateRefundWithContext gets a card charge with context
func CreateRefundWithContext(ctx context.Context, data *CreateRefundParams) (*xendit.CardRefund, *xendit.Error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.CreateRefundWithContext(ctx, data)
}

/* Authorization */

// ReverseAuthorization reverses a card authorization
func ReverseAuthorization(data *ReverseAuthorizationParams) (*xendit.CardReverseAuthorization, *xendit.Error) {
	return ReverseAuthorizationWithContext(context.Background(), data)
}

// ReverseAuthorizationWithContext reverses a card authorization with context
func ReverseAuthorizationWithContext(ctx context.Context, data *ReverseAuthorizationParams) (*xendit.CardReverseAuthorization, *xendit.Error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.ReverseAuthorizationWithContext(ctx, data)
}

func getClient() (*Client, *xendit.Error) {
	return &Client{
		Opt:          &xendit.Opt,
		APIRequester: xendit.GetAPIRequester(),
	}, nil
}
