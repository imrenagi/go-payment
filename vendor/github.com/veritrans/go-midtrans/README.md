# Midtrans Library for Go(lang)

[![Go Report Card](https://goreportcard.com/badge/github.com/veritrans/go-midtrans)](https://goreportcard.com/report/github.com/veritrans/go-midtrans)
[![Apache 2.0 license](https://img.shields.io/badge/license-Apache%202.0-brightgreen.svg)](LICENSE)
[![Build Status](https://travis-ci.org/veritrans/go-midtrans.svg?branch=master)](https://travis-ci.org/veritrans/go-midtrans)

Midtrans :heart: Go !

Go is a very modern, terse, and combine aspect of dynamic and static typing that in a way very
well suited for web development, among other things. Its small memory footprint is also
an advantage of itself. Now, Midtrans is available to be used in Go, too.

## Usage blueprint

1. There is a type named `Client` (`midtrans.Client`) that should be instantiated through `NewClient` which hold any possible setting to the library.
2. There is a gateway classes which you will be using depending on whether you used Core, SNAP, SNAP-Redirect. The gateway type need a Client instance.
3. Any activity (charge, approve, etc) is done in the gateway level.

## Example

We have attached usage examples in this repository in folder `example/simplepay`.
Please proceed there for more detail on how to run the example.

### Core Gateway

```go
    midclient := midtrans.NewClient()
    midclient.ServerKey = "YOUR-VT-SERVER-KEY"
    midclient.ClientKey = "YOUR-VT-CLIENT-KEY"
    midclient.APIEnvType = midtrans.Sandbox

    coreGateway := midtrans.CoreGateway{
        Client: midclient,
    }

    chargeReq := &midtrans.ChargeReq{
        PaymentType: midtrans.SourceCreditCard,
        TransactionDetails: midtrans.TransactionDetails{
            OrderID: "12345",
            GrossAmt: 200000,
        },
        CreditCard: &midtrans.CreditCardDetail{
            TokenID: "YOUR-CC-TOKEN",
        },
        Items: &[]midtrans.ItemDetail{
            midtrans.ItemDetail{
                ID: "ITEM1",
                Price: 200000,
                Qty: 1,
                Name: "Someitem",
            },
        },
    }

    resp, _ := coreGateway.Charge(chargeReq)
```

### How Core API does charge with map type
please refer to file `main.go` in folder `example/simplepay`
```go
func ChargeWithMap(w http.ResponseWriter, r *http.Request) {
	var reqPayload = &midtrans.ChargeReqWithMap{}
	err := json.NewDecoder(r.Body).Decode(reqPayload)
	if err != nil {
		// do something
		return
	}
	
	chargeResp, _ := coreGateway.ChargeWithMap(reqPayload)
	result, err := json.Marshal(chargeResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}
```

### Snap Gateway

Snap is Midtrans existing tool to help merchant charge customers using a
mobile-friendly, in-page, no-redirect checkout facilities. Using snap is
completely simple.

Quick transaction with minimum Snap parameters:

```go
var snapGateway midtrans.SnapGateway
snapGateway = midtrans.SnapGateway{
  Client: midclient,
}

snapResp, err := snapGateway.GetTokenQuick(generateOrderId(), 200000)
var snapToken string
if err != nil {
  snapToken = snapResp.Token
}
```

On the client side:

```javascript
var token = $("#snap-token").val();
snap.pay(token, {
    onSuccess: function(res) { alert("Payment accepted!"); },
    onPending: function(res) { alert("Payment pending", res); },
    onError: function(res) { alert("Error", res); }
});
```

You may want to override those `onSuccess`, `onPending` and `onError`
functions to reflect the behaviour that you wished when the charging
result in their respective state.

Alternativelly, more complete Snap parameter:

```go

    midclient := midtrans.NewClient()
    midclient.ServerKey = "YOUR-VT-SERVER-KEY"
    midclient.ClientKey = "YOUR-VT-CLIENT-KEY"
    midclient.APIEnvType = midtrans.Sandbox

    var snapGateway midtrans.SnapGateway
    snapGateway = midtrans.SnapGateway{
        Client: midclient,
    }

    custAddress := &midtrans.CustAddress{
        FName: "John",
        LName: "Doe",
        Phone: "081234567890",
        Address: "Baker Street 97th",
        City: "Jakarta",
        Postcode: "16000",
        CountryCode: "IDN",
    }

    snapReq := &midtrans.SnapReq{
        TransactionDetails: midtrans.TransactionDetails{
            OrderID: "order-id-go-"+timestamp,
            GrossAmt: 200000,
        },
        CustomerDetail: &midtrans.CustDetail{
            FName: "John",
            LName: "Doe",
            Email: "john@doe.com",
            Phone: "081234567890",
            BillAddr: custAddress,
            ShipAddr: custAddress,
        },
        Items: &[]midtrans.ItemDetail{
            midtrans.ItemDetail{
                ID: "ITEM1",
                Price: 200000,
                Qty: 1,
                Name: "Someitem",
            },
        },
    }

    log.Println("GetToken:")
    snapTokenResp, err := snapGateway.GetToken(snapReq)
```

### Handle HTTP Notification
Create separated web endpoint (notification url) to receive HTTP POST notification callback/webhook. 
HTTP notification will be sent whenever transaction status is changed.
Example also available in `main.go` in folder `example/simplepay` 

```go
func notification(w http.ResponseWriter, r *http.Request) {
	var reqPayload = &midtrans.ChargeReqWithMap{}
	err := json.NewDecoder(r.Body).Decode(reqPayload)
	if err != nil {
		// do something
		return
	}

	encode, _ := json.Marshal(reqPayload)
	resArray := make(map[string]string)
	err = json.Unmarshal(encode, &resArray)

	chargeResp, _ := coreGateway.StatusWithMap(resArray["order_id"])
	result, err := json.Marshal(chargeResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}
```

### Iris Gateway
Iris is Midtrans cash management solution that allows you to disburse payments to any bank accounts in Indonesia securely and easily. Iris connects to the banks’ hosts to enable seamless transfer using integrated APIs.

>Note: `ServerKey` used for `irisGateway`'s `Client` is API Key found in Iris Dashboard. The API Key is different with Midtrans' payment gateway account's key.

```go
var irisGateway midtrans.IrisGateway
irisGateway = midtrans.IrisGateway{
  Client: midclient,
}

res, err := irisGateway.GetListBeneficiaryBank()
if err != nil {
    fmt.Println("err ", err)
    return
}
fmt.Printf("result %v \n ", res.BeneficiaryBanks)
```

## License

See [LICENSE](LICENSE).
