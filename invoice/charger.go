package invoice

import "context"

// PaymentCharger will call the API
type PaymentCharger interface {
	Create(ctx context.Context, inv *Invoice) (*ChargeResponse, error)
	Gateway() string
}

// ChargeResponse stores the important data after a payment charge request
// has been created
type ChargeResponse struct {
	TransactionID string
	PaymentToken  string
	PaymentURL    string
}
