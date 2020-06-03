package subscription

import (
	"context"

	"github.com/imrenagi/go-payment"
)

type gateway interface {
	Gateway() payment.Gateway
}

type creator interface {
	gateway
	Create(ctx context.Context, sub *Subscription) (*CreateResponse, error)
}

// Controller is payment gateway interface for subscription handling
type Controller interface {
	creator
}
