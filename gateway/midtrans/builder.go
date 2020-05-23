package midtrans

import (
	"fmt"
	"math"
	"time"

	"github.com/imrenagi/go-payment/invoice"

	gomidtrans "github.com/veritrans/go-midtrans"
)

type requestBuilder interface {
	Build() (*gomidtrans.SnapReq, error)
}

func NewSnapRequestBuilder(inv *invoice.Invoice) *SnapRequestBuilder {
	srb := &SnapRequestBuilder{
		snapRequest: &gomidtrans.SnapReq{
			Items: &[]gomidtrans.ItemDetail{},
		},
	}

	return srb.
		SetTransactionDetails(inv).
		SetCustomerDetail(inv).
		SetExpiration(inv).
		SetItemDetails(inv)
}

type SnapRequestBuilder struct {
	snapRequest *gomidtrans.SnapReq
}

func (srb *SnapRequestBuilder) concatProductName(provider, product string) string {

	sProviderName := provider
	if len(provider) > 10 {
		runes := []rune(provider)
		sProviderName = fmt.Sprintf("%s...", string(runes[0:10]))
	}

	sProductName := product
	if len(product) > 30 {
		runes := []rune(product)
		sProductName = fmt.Sprintf("%s...", string(runes[0:30]))
	}

	return fmt.Sprintf("%s - %s", sProviderName, sProductName)
}

func (srb *SnapRequestBuilder) SetItemDetails(inv *invoice.Invoice) *SnapRequestBuilder {
	var out []gomidtrans.ItemDetail
	i := inv.LineItem
	out = append(out, gomidtrans.ItemDetail{
		ID:           i.Category,
		Name:         i.Name,
		Price:        int64(i.UnitPrice),
		Qty:          int32(i.Qty),
		Category:     i.Category,
		MerchantName: i.MerchantName,
	})

	if inv.ServiceFee > 0 {
		out = append(out, gomidtrans.ItemDetail{
			ID:       "adminfee",
			Name:     "Biaya Admin",
			Price:    int64(inv.ServiceFee),
			Qty:      1,
			Category: "FEE",
		})
	}
	if inv.InstallmentFee > 0 {
		out = append(out, gomidtrans.ItemDetail{
			ID:       "installmentfee",
			Name:     "Installment Fee",
			Price:    int64(inv.InstallmentFee),
			Qty:      1,
			Category: "FEE",
		})
	}
	if inv.Discount > 0 {
		out = append(out, gomidtrans.ItemDetail{
			ID:       "discount",
			Name:     "Discount",
			Price:    int64(-1 * inv.Discount),
			Qty:      1,
			Category: "DISCOUNT",
		})
	}

	srb.snapRequest.Items = &out

	return srb
}

func (srb *SnapRequestBuilder) SetCustomerDetail(inv *invoice.Invoice) *SnapRequestBuilder {
	srb.snapRequest.CustomerDetail = &gomidtrans.CustDetail{
		FName: inv.BillingAddress.FullName,
		Email: inv.BillingAddress.Email,
		Phone: inv.BillingAddress.PhoneNumber,
		BillAddr: &gomidtrans.CustAddress{
			FName: inv.BillingAddress.FullName,
			Phone: inv.BillingAddress.PhoneNumber,
		},
	}
	return srb
}

func (srb *SnapRequestBuilder) SetExpiration(inv *invoice.Invoice) *SnapRequestBuilder {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	invDate := inv.InvoiceDate.In(loc)
	duration := inv.DueDate.Sub(inv.InvoiceDate)
	srb.snapRequest.Expiry = &gomidtrans.ExpiryDetail{
		StartTime: invDate.Format("2006-01-02 15:04:05 -0700"),
		Unit:      "hour",
		Duration:  int64(math.Round(duration.Hours())),
	}
	return srb
}

func (srb *SnapRequestBuilder) SetTransactionDetails(inv *invoice.Invoice) *SnapRequestBuilder {
	srb.snapRequest.TransactionDetails = gomidtrans.TransactionDetails{
		OrderID:  inv.Number,
		GrossAmt: int64(inv.GetTotal()),
	}
	return srb
}

func (srb *SnapRequestBuilder) Build() (*gomidtrans.SnapReq, error) {
	return srb.snapRequest, nil
}
