package manage

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/imrenagi/go-payment"

	"github.com/imrenagi/go-payment/gateway/xendit"
	factory "github.com/imrenagi/go-payment/gateway/xendit"
	"github.com/imrenagi/go-payment/subscription"

	goxendit "github.com/xendit/xendit-go"
	xrecurring "github.com/xendit/xendit-go/recurringpayment"
)

type xenditSubscriptionController struct {
	XenditGateway *xendit.Gateway
}

func (sc xenditSubscriptionController) Create(ctx context.Context, sub *subscription.Subscription) (*subscription.CreateResponse, error) {

	l := log.Ctx(ctx).With().
		Str("function", "xenditSubscriptionController.Create").
		Logger()

	recurringRequest, err := factory.NewRecurringChargeRequestBuilder(sub).Build()
	if err != nil {
		return nil, err
	}

	bytes, err := json.MarshalIndent(recurringRequest, "", "\t")
	if err != nil {
		return nil, err
	}
	l.Debug().
		RawJSON("payload", bytes).
		Msg("recurring request payload is created")

	xres, err := sc.XenditGateway.Recurring.CreateWithContext(ctx, recurringRequest)
	if err != nil {
		var xError *goxendit.Error
		if ok := errors.As(err, &xError); ok && xError != nil {
			l.Error().Err(xError).Msg("unable to create recurring request")
			return nil, xError
		}
	}

	l.Info().Msg("recurring request is created successfully")

	return &subscription.CreateResponse{
		ID:                    xres.ID,
		Status:                xendit.NewStatus(xres.Status),
		LastCreatedInvoiceURL: xres.LastCreatedInvoiceURL,
	}, nil
}

func (sc xenditSubscriptionController) Resume(ctx context.Context, sub *subscription.Subscription) error {

	l := log.Ctx(ctx).With().
		Str("function", "xenditSubscriptionController.Resume").
		Logger()

	_, err := sc.XenditGateway.Recurring.ResumeWithContext(ctx, &xrecurring.ResumeParams{
		ID: sub.GatewayRecurringID,
	})
	if err != nil {
		var xError *goxendit.Error
		if ok := errors.As(err, &xError); ok && xError != nil {
			l.Error().Err(xError).Msg("unable to resume subscription")
			return xError
		}
	}

	l.Info().Msg("resume subscription request is completed")
	return nil
}

func (sc xenditSubscriptionController) Stop(ctx context.Context, sub *subscription.Subscription) error {

	l := log.Ctx(ctx).With().
		Str("function", "xenditSubscriptionController.Stop").
		Logger()

	_, err := sc.XenditGateway.Recurring.StopWithContext(ctx, &xrecurring.StopParams{
		ID: sub.GatewayRecurringID,
	})
	if err != nil {
		var xError *goxendit.Error
		if ok := errors.As(err, &xError); ok && xError != nil {
			l.Error().Err(xError).Msg("unable to stop subscription")
			return xError
		}
	}
	l.Info().Msg("stop subscription request is completed")
	return nil
}

func (sc xenditSubscriptionController) Pause(ctx context.Context, sub *subscription.Subscription) error {

	l := log.Ctx(ctx).With().
		Str("function", "xenditSubscriptionController.Stop").
		Logger()

	_, err := sc.XenditGateway.Recurring.PauseWithContext(ctx, &xrecurring.PauseParams{
		ID: sub.GatewayRecurringID,
	})
	if err != nil {
		var xError *goxendit.Error
		if ok := errors.As(err, &xError); ok && xError != nil {
			l.Error().Err(xError).Msg("unable to pause subscription")
			return xError
		}
	}
	l.Info().Msg("pause subscription request is completed")
	return nil
}

func (sc xenditSubscriptionController) Gateway() payment.Gateway {
	return payment.GatewayXendit
}
