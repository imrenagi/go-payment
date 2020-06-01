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

func (srb *SnapRequestBuilder) SetItemDetails(inv *invoice.Invoice) *SnapRequestBuilder {
	var out []gomidtrans.ItemDetail

	for _, item := range inv.LineItems {

		name := item.Name
		if len(item.Name) > 50 {
			runes := []rune(name)
			name = fmt.Sprintf("%s", string(runes[0:50]))
		}

		out = append(out, gomidtrans.ItemDetail{
			ID:           item.Category,
			Name:         name,
			Price:        int64(item.UnitPrice),
			Qty:          int32(item.Qty),
			Category:     item.Category,
			MerchantName: item.MerchantName,
		})
	}

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
