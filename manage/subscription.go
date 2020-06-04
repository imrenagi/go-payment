package manage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/imrenagi/go-payment"

	"github.com/imrenagi/go-payment/gateway/xendit"
	factory "github.com/imrenagi/go-payment/gateway/xendit"
	"github.com/imrenagi/go-payment/subscription"

	goxendit "github.com/xendit/xendit-go"
)

type xenditSubscriptionController struct {
	XenditGateway *xendit.Gateway
}

func (sc xenditSubscriptionController) Create(ctx context.Context, sub *subscription.Subscription) (*subscription.CreateResponse, error) {
	recurringRequest, err := factory.NewRecurringChargeRequestBuilder(sub).Build()
	if err != nil {
		return nil, err
	}

	bytes, err := json.MarshalIndent(recurringRequest, "", "\t")
	if err != nil {
		return nil, err
	}
	fmt.Println(string(bytes))

	xres, err := sc.XenditGateway.Recurring.CreateWithContext(ctx, recurringRequest)
	var xError *goxendit.Error
	if ok := errors.As(err, &xError); ok && xError != nil {
		return nil, xError
	}

	return &subscription.CreateResponse{
		ID:                    xres.ID,
		Status:                xendit.NewStatus(xres.Status),
		LastCreatedInvoiceURL: xres.LastCreatedInvoiceURL,
	}, nil
}

func (sc xenditSubscriptionController) Gateway() payment.Gateway {
	return payment.GatewayXendit
}
