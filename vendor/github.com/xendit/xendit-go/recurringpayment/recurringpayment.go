package recurringpayment

import (
	"context"

	"github.com/xendit/xendit-go"
)

// Create creates new recurring payment
func Create(data *CreateParams) (*xendit.RecurringPayment, *xendit.Error) {
	return CreateWithContext(context.Background(), data)
}

// CreateWithContext creates new recurring payment with context
func CreateWithContext(ctx context.Context, data *CreateParams) (*xendit.RecurringPayment, *xendit.Error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.CreateWithContext(ctx, data)
}

// Get gets a recurring payment
func Get(data *GetParams) (*xendit.RecurringPayment, *xendit.Error) {
	return GetWithContext(context.Background(), data)
}

// GetWithContext gets a recurring payment with context
func GetWithContext(ctx context.Context, data *GetParams) (*xendit.RecurringPayment, *xendit.Error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.GetWithContext(ctx, data)
}

// Edit gets a recurring payment
func Edit(data *EditParams) (*xendit.RecurringPayment, *xendit.Error) {
	return EditWithContext(context.Background(), data)
}

// EditWithContext gets a recurring payment with context
func EditWithContext(ctx context.Context, data *EditParams) (*xendit.RecurringPayment, *xendit.Error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.EditWithContext(ctx, data)
}

// Stop gets a recurring payment
func Stop(data *StopParams) (*xendit.RecurringPayment, *xendit.Error) {
	return StopWithContext(context.Background(), data)
}

// StopWithContext gets a recurring payment with context
func StopWithContext(ctx context.Context, data *StopParams) (*xendit.RecurringPayment, *xendit.Error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.StopWithContext(ctx, data)
}

// Pause gets a recurring payment
func Pause(data *PauseParams) (*xendit.RecurringPayment, *xendit.Error) {
	return PauseWithContext(context.Background(), data)
}

// PauseWithContext gets a recurring payment with context
func PauseWithContext(ctx context.Context, data *PauseParams) (*xendit.RecurringPayment, *xendit.Error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.PauseWithContext(ctx, data)
}

// Resume gets a recurring payment
func Resume(data *ResumeParams) (*xendit.RecurringPayment, *xendit.Error) {
	return ResumeWithContext(context.Background(), data)
}

// ResumeWithContext gets a recurring payment with context
func ResumeWithContext(ctx context.Context, data *ResumeParams) (*xendit.RecurringPayment, *xendit.Error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.ResumeWithContext(ctx, data)
}

func getClient() (*Client, *xendit.Error) {
	return &Client{
		Opt:          &xendit.Opt,
		APIRequester: xendit.GetAPIRequester(),
	}, nil
}
