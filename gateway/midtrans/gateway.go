package midtrans

import (
	"os"

	"github.com/imrenagi/go-payment/util/localconfig"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	mgo "github.com/veritrans/go-midtrans"
)

// NewGateway creates new midtrans payment gateway
func NewGateway(creds localconfig.APICredential) *Gateway {

	midclient := mgo.NewClient()
	midclient.ServerKey = creds.SecretKey
	midclient.ClientKey = creds.ClientKey

	snapClient := &snap.Client{
		ServerKey: creds.SecretKey,
	}

	switch os.Getenv("ENVIRONMENT") {
	case "prod":
		midclient.APIEnvType = mgo.Production
		snapClient.Env = midtrans.Production
		snapClient.HttpClient = midtrans.GetHttpClient(midtrans.Production)
	default:
		midclient.APIEnvType = mgo.Sandbox
		snapClient.Env = midtrans.Sandbox
		snapClient.HttpClient = midtrans.GetHttpClient(midtrans.Sandbox)
	}

	gateway := Gateway{
		midClient: midclient,
		SnapGateway: &mgo.SnapGateway{
			Client: midclient,
		},
		SnapV2Gateway: snapClient,
	}

	return &gateway
}

// Gateway stores go-midtrans gateway and client
type Gateway struct {
	midClient     mgo.Client
	SnapGateway   SnapGateway
	SnapV2Gateway SnapV2Gateway
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

type SnapV2Gateway interface {
	CreateTransaction(req *snap.Request) (*snap.Response, *midtrans.Error)
}
