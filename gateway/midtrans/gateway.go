package midtrans

import (
	"os"

	"github.com/imrenagi/go-payment/util/localconfig"
	gomidtrans "github.com/veritrans/go-midtrans"
)

// NewGateway creates new midtrans payment gateway
func NewGateway(creds localconfig.APICredential) *Gateway {

	midclient := gomidtrans.NewClient()
	midclient.ServerKey = creds.SecretKey
	midclient.ClientKey = creds.ClientKey

	switch os.Getenv("ENVIRONMENT") {
	case "prod":
		midclient.APIEnvType = gomidtrans.Production
	default:
		midclient.APIEnvType = gomidtrans.Sandbox
	}

	gateway := Gateway{
		midClient: midclient,
		CoreGateway: &gomidtrans.CoreGateway{
			Client: midclient,
		},
		SnapGateway: &gomidtrans.SnapGateway{
			Client: midclient,
		},
	}

	return &gateway
}

// Gateway stores go-midtrans gateway and client
type Gateway struct {
	midClient   gomidtrans.Client
	CoreGateway CoreGateway
	SnapGateway SnapGateway
}

// GetServerKey returns midtrans server key
func (g Gateway) GetServerKey() string {
	return g.midClient.ServerKey
}

type SnapGateway interface {
	GetToken(r *gomidtrans.SnapReq) (gomidtrans.SnapResponse, error)
}

type CoreGateway interface {
	Status(orderID string) (gomidtrans.Response, error)
}
