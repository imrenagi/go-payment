package invoice

import (
	"context"

	"github.com/imrenagi/go-payment"
)

// PaymentCharger will call the API
type PaymentCharger interface {
	Create(ctx context.Context, inv *Invoice) (*ChargeResponse, error)
	Gateway() payment.Gateway
}

// ChargeResponse stores the important data after a payment charge request
// has been created
type ChargeResponse struct {
	TransactionID string
	PaymentToken  string
	PaymentURL    string
}
