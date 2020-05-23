package service

import (
	"context"

	midfactory "github.com/imrenagi/go-payment/gateway/midtrans"
	midgateway "github.com/imrenagi/go-payment/gateway/midtrans"
	"github.com/imrenagi/go-payment/invoice"
)

type midtransCharger struct {
	MidtransGateway *midgateway.Gateway
}

func (c midtransCharger) Create(ctx context.Context, inv *invoice.Invoice) (*invoice.ChargeResponse, error) {

	snapRequest, err := midfactory.NewSnapRequestFromInvoice(inv)
	if err != nil {
		return nil, err
	}

	// bytes, err := json.MarshalIndent(snapRequest, "", "\t")
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Println(string(bytes))

	resp, err := c.MidtransGateway.SnapGateway.GetToken(snapRequest)
	if err != nil {
		return nil, err
	}

	return &invoice.ChargeResponse{
		PaymentToken: resp.Token,
		PaymentURL:   resp.RedirectURL,
		// TransactionID: resp.

	}, nil
}

func (c midtransCharger) Gateway() string {
	// TODO change this later
	return "midtrans"
}
