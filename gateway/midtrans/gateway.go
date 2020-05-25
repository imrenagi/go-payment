package midtrans

import (
	"os"

	"github.com/imrenagi/go-payment/util/localconfig"
	mgo "github.com/veritrans/go-midtrans"
)

// NewGateway creates new midtrans payment gateway
func NewGateway(creds localconfig.APICredential) *Gateway {

	midclient := mgo.NewClient()
	midclient.ServerKey = creds.SecretKey
	midclient.ClientKey = creds.ClientKey

	switch os.Getenv("ENVIRONMENT") {
	case "prod":
		midclient.APIEnvType = mgo.Production
	default:
		midclient.APIEnvType = mgo.Sandbox
	}

	gateway := Gateway{
		midClient: midclient,
		SnapGateway: &mgo.SnapGateway{
			Client: midclient,
		},
	}

	return &gateway
}

// Gateway stores go-midtrans gateway and client
type Gateway struct {
	midClient   mgo.Client
	SnapGateway SnapGateway
}

// NotificationValidationKey returns midtrans server key used for validating
// midtransa transaction status
func (g Gateway) NotificationValidationKey() string {
	return g.midClient.ServerKey
}

// SnapGateway interaction interface with midtrans snap server
type SnapGateway interface {
	GetToken(r *mgo.SnapReq) (mgo.SnapResponse, error)
}
