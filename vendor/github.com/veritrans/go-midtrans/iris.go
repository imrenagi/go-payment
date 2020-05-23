package midtrans

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

// IrisGateway struct
type IrisGateway struct {
	Client Client
}

// Call : base method to call IRIS API
func (gateway *IrisGateway) Call(method, path string, body io.Reader, v interface{}) error {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	path = gateway.Client.APIEnvType.IrisURL() + path
	return gateway.Client.Call(method, path, body, v)
}

// GetListBeneficiaryBank : Show list of supported banks in IRIS. (https://iris-docs.midtrans.com/#list-banks)
func (gateway *IrisGateway) GetListBeneficiaryBank() (IrisBeneficiaryBanksResponse, error) {
	resp := IrisBeneficiaryBanksResponse{}

	err := gateway.Call("GET", "api/v1/beneficiary_banks", nil, &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error getting beneficiary banks: ", err)
		return resp, err
	}

	return resp, nil
}

// CreateBeneficiaries : Create Beneficiaries (https://iris-docs.midtrans.com/#create-beneficiaries)
func (gateway *IrisGateway) CreateBeneficiaries(req *IrisBeneficiaries) (bool, error) {
	resp := IrisBeneficiariesResponse{}
	jsonReq, _ := json.Marshal(req)

	err := gateway.Call("POST", "api/v1/beneficiaries", bytes.NewBuffer(jsonReq), &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error creating beneficiaries: ", err)
		return false, err
	}

	if resp.Status != "created" {
		gateway.Client.Logger.Println("Error creating beneficiaries: ", resp.Errors)
		return false, errors.New(strings.Join(resp.Errors, ","))
	}

	return true, nil
}

// UpdateBeneficiaries : Update Beneficiaries (https://iris-docs.midtrans.com/#update-beneficiaries)
func (gateway *IrisGateway) UpdateBeneficiaries(aliasName string, req *IrisBeneficiaries) (bool, error) {
	resp := IrisBeneficiariesResponse{}
	jsonReq, _ := json.Marshal(req)

	err := gateway.Call("PATCH", fmt.Sprintf("api/v1/beneficiaries/%s", aliasName), bytes.NewBuffer(jsonReq), &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error updating beneficiaries: ", err)
		return false, err
	}

	if resp.Status != "updated" {
		gateway.Client.Logger.Println("Error updating beneficiaries: ", resp.Errors)
		return false, errors.New(strings.Join(resp.Errors, ","))
	}

	return true, nil
}

// GetListBeneficiaries : Get List Beneficiaries (https://iris-docs.midtrans.com/#list-beneficiaries)
func (gateway *IrisGateway) GetListBeneficiaries() ([]IrisBeneficiaries, error) {
	var resp []IrisBeneficiaries

	err := gateway.Call("GET", "api/v1/beneficiaries", nil, &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error get list beneficiaries: ", err)
		return resp, err
	}

	return resp, nil
}

// CreatePayouts : This API is for Creator to create a payout. It can be used for single payout and also multiple payouts. (https://iris-docs.midtrans.com/#create-payouts)
func (gateway *IrisGateway) CreatePayouts(req IrisCreatePayoutReq) (IrisCreatePayoutResponse, error) {
	resp := IrisCreatePayoutResponse{}
	jsonReq, _ := json.Marshal(req)

	err := gateway.Call("POST", "api/v1/payouts", bytes.NewBuffer(jsonReq), &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error creating payouts: ", err)
		return resp, err
	}

	if resp.ErrorMessage != "" {
		return resp, errors.New(resp.ErrorMessage)
	}

	return resp, nil
}

// ApprovePayouts : Use this API for Apporver to approve multiple payout request. (https://iris-docs.midtrans.com/#approve-payouts)
func (gateway *IrisGateway) ApprovePayouts(req IrisApprovePayoutReq) (IrisApprovePayoutResponse, error) {
	resp := IrisApprovePayoutResponse{}
	jsonReq, _ := json.Marshal(req)

	err := gateway.Call("POST", "api/v1/payouts/approve", bytes.NewBuffer(jsonReq), &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error approving payouts: ", err)
		return resp, err
	}

	if len(resp.Errors) > 0 {
		return resp, errors.New(strings.Join(resp.Errors, ", "))
	}

	if resp.Status != "ok" {
		return resp, errors.New("Error approving payouts, status from API not OK")
	}

	return resp, nil
}

// RejectPayouts : Use this API for Apporver to reject multiple payout request. (https://iris-docs.midtrans.com/#reject-payouts)
func (gateway *IrisGateway) RejectPayouts(req IrisRejectPayoutReq) (IrisRejectPayoutResponse, error) {
	resp := IrisRejectPayoutResponse{}
	jsonReq, _ := json.Marshal(req)

	err := gateway.Call("POST", "api/v1/payouts/reject", bytes.NewBuffer(jsonReq), &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error rejecting payouts: ", err)
		return resp, err
	}

	if len(resp.Errors) > 0 {
		return resp, errors.New(strings.Join(resp.Errors, ", "))
	}

	if resp.Status != "ok" {
		return resp, errors.New("Error rejecting payouts, status from API not OK")
	}

	return resp, nil
}

// GetPayoutDetails : Get details of a single payout (https://iris-docs.midtrans.com/#get-payout-details)
func (gateway *IrisGateway) GetPayoutDetails(referenceNo string) (IrisPayoutDetailResponse, error) {
	resp := IrisPayoutDetailResponse{}

	// handle conflict call with payout history (https://iris-docs.midtrans.com/#payout-history)
	if referenceNo == "" {
		return resp, errors.New("you must specified referenceNo")
	}

	err := gateway.Call("GET", "api/v1/payouts/"+referenceNo, nil, &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error getting payout details: ", err)
		return resp, err
	}

	if resp.ErrorMessage != "" {
		return resp, errors.New(resp.ErrorMessage)
	}

	return resp, nil
}

// ValidateBankAccount : Check if an account is valid, if valid return account information. (https://iris-docs.midtrans.com/#validate-bank-account)
func (gateway *IrisGateway) ValidateBankAccount(bankName string, accountNo string) (IrisBankAccountDetailResponse, error) {
	resp := IrisBankAccountDetailResponse{}

	err := gateway.Call("GET", fmt.Sprintf("api/v1/account_validation?bank=%s&account=%s", bankName, accountNo), nil, &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error validating bank account: ", err)
		return resp, err
	}

	if resp.ErrorMessage != "" {
		errMsg := []string{}
		if len(resp.Errors.Account) > 0 {
			accountErr := "account: " + strings.Join(resp.Errors.Account, ", ")
			errMsg = append(errMsg, accountErr)
		}

		if len(resp.Errors.Bank) > 0 {
			bankErr := "bank: " + strings.Join(resp.Errors.Bank, ", ")
			errMsg = append(errMsg, bankErr)
		}

		return resp, errors.New(strings.Join(errMsg, " & "))
	}

	return resp, nil
}

// CheckBalance : Check Balance (Aggregator) (https://iris-docs.midtrans.com/#check-balance-aggregator)
func (gateway *IrisGateway) CheckBalance() (IrisBalanceResponse, error) {
	resp := IrisBalanceResponse{}

	err := gateway.Call("GET", "api/v1/balance", nil, &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error check balance: ", err)
		return resp, err
	}

	return resp, nil
}

// GetPayoutHistory : Returns all the payout details for specific dates (https://iris-docs.midtrans.com/#payout-history)
func (gateway *IrisGateway) GetPayoutHistory(fromDate string, toDate string) ([]IrisPayoutDetailResponse, error) {
	resp := []IrisPayoutDetailResponse{}

	err := gateway.Call("GET", fmt.Sprintf("api/v1/payouts?from_date=%s&to_date=%s", fromDate, toDate), nil, &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error get payout history: ", err)
		return resp, err
	}

	return resp, nil
}
