//go:generate mockery -dir . -name Controller -output ./mocks -filename controller.go

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

type pauser interface {
	Pause(ctx context.Context, sub *Subscription) error
}

type stopper interface {
	Stop(ctx context.Context, stop *Subscription) error
}

type resumer interface {
	Resume(ctx context.Context, sub *Subscription) error
}

// Controller is payment gateway interface for subscription handling
type Controller interface {
	creator
	pauser
	stopper
	resumer
}
