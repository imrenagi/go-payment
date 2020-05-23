package midtrans

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
)

// CoreGateway struct
type CoreGateway struct {
	Client Client
}

// Call : base method to call Core API
func (gateway *CoreGateway) Call(method, path string, body io.Reader, v interface{}) error {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	path = gateway.Client.APIEnvType.String() + path
	return gateway.Client.Call(method, path, body, v)
}

// Charge : Perform transaction using ChargeReq
func (gateway *CoreGateway) Charge(req *ChargeReq) (Response, error) {
	resp := Response{}
	jsonReq, _ := json.Marshal(req)

	err := gateway.Call("POST", "v2/charge", bytes.NewBuffer(jsonReq), &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error charging: ", err)
		return resp, err
	}

	if resp.StatusMessage != "" {
		gateway.Client.Logger.Println(resp.StatusMessage)
	}

	return resp, nil
}

// ChargeWithMap : Perform transaction using ChargeReqWithMap
func (gateway *CoreGateway) ChargeWithMap(req *ChargeReqWithMap) (ResponseWithMap, error) {
	resp := ResponseWithMap{}
	jsonReq, _ := json.Marshal(req)

	err := gateway.Call("POST", "v2/charge", bytes.NewBuffer(jsonReq), &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error charging: ", err)
		return resp, err
	}

	if resp["status_code"] != nil {
		gateway.Client.Logger.Println("==== Status Code === : ", resp["status_code"])
		gateway.Client.Logger.Println("==== Message === : ", resp["status_message"])
	} else {
		gateway.Client.Logger.Println("==== Error Message === : ", resp["message"])
	}

	return resp, nil
}

// PreauthCard : Perform authorized transactions using ChargeReq
func (gateway *CoreGateway) PreauthCard(req *ChargeReq) (Response, error) {
	req.CreditCard.Type = "authorize"
	return gateway.Charge(req)
}

// CaptureCard : Capture an authorized transaction for card payment
func (gateway *CoreGateway) CaptureCard(req *CaptureReq) (Response, error) {
	resp := Response{}
	jsonReq, _ := json.Marshal(req)

	err := gateway.Call("POST", "v2/capture", bytes.NewBuffer(jsonReq), &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error capturing: ", err)
		return resp, err
	}

	if resp.StatusMessage != "" {
		gateway.Client.Logger.Println(resp.StatusMessage)
	}

	return resp, nil
}

// Approve : Approve order using order ID
func (gateway *CoreGateway) Approve(orderID string) (Response, error) {
	resp := Response{}

	err := gateway.Call("POST", "v2/"+orderID+"/approve", nil, &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error approving: ", err)
		return resp, err
	}

	if resp.StatusMessage != "" {
		gateway.Client.Logger.Println(resp.StatusMessage)
	}

	return resp, nil
}

// Cancel : Cancel order using order ID
func (gateway *CoreGateway) Cancel(orderID string) (Response, error) {
	resp := Response{}

	err := gateway.Call("POST", "v2/"+orderID+"/cancel", nil, &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error while cancel: ", err)
		return resp, err
	}

	if resp.StatusMessage != "" {
		gateway.Client.Logger.Println(resp.StatusMessage)
	}

	return resp, nil
}

// Expire : change order status to expired using order ID
func (gateway *CoreGateway) Expire(orderID string) (Response, error) {
	resp := Response{}

	err := gateway.Call("POST", "v2/"+orderID+"/expire", nil, &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error while expire: ", err)
		return resp, err
	}

	if resp.StatusMessage != "" {
		gateway.Client.Logger.Println(resp.StatusMessage)
	}

	return resp, nil
}

// Status : get order status using order ID
func (gateway *CoreGateway) Status(orderID string) (Response, error) {
	resp := Response{}

	err := gateway.Call("GET", "v2/"+orderID+"/status", nil, &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error while get status: ", err)
		return resp, err
	}

	if resp.StatusMessage != "" {
		gateway.Client.Logger.Println(resp.StatusMessage)
	}

	return resp, nil
}

// Refund : refund order using order ID
func (gateway *CoreGateway) Refund(orderID string, req *RefundReq) (Response, error) {
	resp := Response{}
	jsonReq, _ := json.Marshal(req)

	err := gateway.Call("POST", "v2/"+orderID+"/refund", bytes.NewBuffer(jsonReq), &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error while refund: ", err)
		return resp, err
	}

	if resp.StatusMessage != "" {
		gateway.Client.Logger.Println(resp.StatusMessage)
	}

	return resp, nil
}

// DirectRefund : refund order using order ID
func (gateway *CoreGateway) DirectRefund(orderID string, req *RefundReq) (Response, error) {
	resp := Response{}
	jsonReq, _ := json.Marshal(req)

	err := gateway.Call("POST", "v2/"+orderID+"/refund/online/direct", bytes.NewBuffer(jsonReq), &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error while direct refund: ", err)
		return resp, err
	}

	if resp.StatusMessage != "" {
		gateway.Client.Logger.Println(resp.StatusMessage)
	}

	return resp, nil
}

// StatusWithMap : get order status using order ID
func (gateway *CoreGateway) StatusWithMap(orderID string) (ResponseWithMap, error) {
	resp := ResponseWithMap{}

	err := gateway.Call("GET", "v2/"+orderID+"/status", nil, &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error approving: ", err)
		return resp, err
	}

	if resp["status_code"] != nil {
		gateway.Client.Logger.Println("==== Status Code === : ", resp["status_code"])
		gateway.Client.Logger.Println("==== Message === : ", resp["status_message"])
	} else {
		gateway.Client.Logger.Println("==== Error Message === : ", resp["message"])
	}

	return resp, nil
}

// Subscribe : Perform transaction using subscriptions
func (gateway *CoreGateway) Subscribe(req *SubscribeReq) (SubscribeResponse, error) {
	resp := SubscribeResponse{}
	jsonReq, _ := json.Marshal(req)

	err := gateway.Call("POST", "v1/subscriptions", bytes.NewBuffer(jsonReq), &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error while subscribing: ", err)
		return resp, err
	}

	if resp.Status != "" {
		gateway.Client.Logger.Println(resp.Status)
	}

	if resp.StatusMessage != "" {
		gateway.Client.Logger.Println(resp.StatusMessage)
	}

	return resp, nil
}

// SubscribeDetail : Perform get subscription details
func (gateway *CoreGateway) SubscribeDetail(subscriptionID string) (SubscribeResponse, error) {
	resp := SubscribeResponse{}

	err := gateway.Call("GET", "v1/subscriptions/"+subscriptionID, nil, &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error while get subscription details: ", err)
		return resp, err
	}

	if resp.Status != "" {
		gateway.Client.Logger.Println(resp.Status)
	}

	if resp.StatusMessage != "" {
		gateway.Client.Logger.Println(resp.StatusMessage)
	}

	return resp, nil
}

// SubscribeUpdate : Perform update a subscription
func (gateway *CoreGateway) SubscribeUpdate(subscriptionID string, req *SubscribeReq) (SubscribeResponse, error) {
	resp := SubscribeResponse{}
	jsonReq, _ := json.Marshal(req)

	err := gateway.Call("PATCH", "v1/subscriptions/"+subscriptionID, bytes.NewBuffer(jsonReq), &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error while update subscription: ", err)
		return resp, err
	}

	if resp.Status != "" {
		gateway.Client.Logger.Println(resp.Status)
	}

	if resp.StatusMessage != "" {
		gateway.Client.Logger.Println(resp.StatusMessage)
	}

	return resp, nil
}

// SubscribeDisable : Perform disable a subscription
func (gateway *CoreGateway) SubscribeDisable(subscriptionID string) (SubscribeResponse, error) {
	resp := SubscribeResponse{}

	err := gateway.Call("POST", "v1/subscriptions/"+subscriptionID+"/disable", nil, &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error while disable subscription: ", err)
		return resp, err
	}

	if resp.Status != "" {
		gateway.Client.Logger.Println(resp.Status)
	}

	if resp.StatusMessage != "" {
		gateway.Client.Logger.Println(resp.StatusMessage)
	}

	return resp, nil
}

// SubscribeEnable : Perform enable a subscription
func (gateway *CoreGateway) SubscribeEnable(subscriptionID string) (SubscribeResponse, error) {
	resp := SubscribeResponse{}

	err := gateway.Call("POST", "v1/subscriptions/"+subscriptionID+"/enable", nil, &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error while enable subscription: ", err)
		return resp, err
	}

	if resp.Status != "" {
		gateway.Client.Logger.Println(resp.Status)
	}

	if resp.StatusMessage != "" {
		gateway.Client.Logger.Println(resp.StatusMessage)
	}

	return resp, nil
}
