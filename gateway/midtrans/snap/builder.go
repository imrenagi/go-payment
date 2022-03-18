package snap

import (
	"fmt"
	"math"
	"time"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"

	"github.com/imrenagi/go-payment/invoice"
)

func newBuilder(inv *invoice.Invoice) *builder {
	srb := &builder{
		req: &snap.Request{
			Items: &[]midtrans.ItemDetails{},
		},
	}

	return srb.
		setTransactionDetails(inv).
		setCustomerDetail(inv).
		setExpiration(inv).
		setItemDetails(inv)
}

type builder struct {
	req *snap.Request
}

func (b *builder) setItemDetails(inv *invoice.Invoice) *builder {
	var out []midtrans.ItemDetails

	for _, item := range inv.LineItems {

		name := item.Name
		if len(item.Name) > 50 {
			runes := []rune(name)
			name = fmt.Sprintf("%s", string(runes[0:50]))
		}

		out = append(out, midtrans.ItemDetails{
			ID:           fmt.Sprintf("%d", item.ID),
			Name:         name,
			Price:        int64(item.UnitPrice),
			Qty:          int32(item.Qty),
			Category:     item.Category,
			MerchantName: item.MerchantName,
		})
	}

	if inv.ServiceFee > 0 {
		out = append(out, midtrans.ItemDetails{
			ID:       "adminfee",
			Name:     "Biaya Admin",
			Price:    int64(inv.ServiceFee),
			Qty:      1,
			Category: "FEE",
		})
	}
	if inv.InstallmentFee > 0 {
		out = append(out, midtrans.ItemDetails{
			ID:       "installmentfee",
			Name:     "Installment Fee",
			Price:    int64(inv.InstallmentFee),
			Qty:      1,
			Category: "FEE",
		})
	}
	if inv.Discount > 0 {
		out = append(out, midtrans.ItemDetails{
			ID:       "discount",
			Name:     "Discount",
			Price:    int64(-1 * inv.Discount),
			Qty:      1,
			Category: "DISCOUNT",
		})
	}
	if inv.Tax > 0 {
		out = append(out, midtrans.ItemDetails{
			ID:       "tax",
			Name:     "Tax",
			Price:    int64(inv.Tax),
			Qty:      1,
			Category: "TAX",
		})
	}

	b.req.Items = &out

	return b
}

func (b *builder) setCustomerDetail(inv *invoice.Invoice) *builder {
	b.req.CustomerDetail = &midtrans.CustomerDetails{
		FName: inv.BillingAddress.FullName,
		Email: inv.BillingAddress.Email,
		Phone: inv.BillingAddress.PhoneNumber,
		BillAddr: &midtrans.CustomerAddress{
			FName: inv.BillingAddress.FullName,
			Phone: inv.BillingAddress.PhoneNumber,
		},
	}
	return b
}

func (b *builder) setExpiration(inv *invoice.Invoice) *builder {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	invDate := inv.InvoiceDate.In(loc)
	duration := inv.DueDate.Sub(inv.InvoiceDate)
	b.req.Expiry = &snap.ExpiryDetails{
		StartTime: invDate.Format("2006-01-02 15:04:05 -0700"),
		Unit:      "minute",
		Duration:  int64(math.Round(duration.Minutes())),
	}
	return b
}

func (b *builder) setTransactionDetails(inv *invoice.Invoice) *builder {
	b.req.TransactionDetails = midtrans.TransactionDetails{
		OrderID:  inv.Number,
		GrossAmt: int64(inv.GetTotal()),
	}
	return b
}

func (b *builder) AddPaymentMethods(m snap.SnapPaymentType) *builder {
	b.req.EnabledPayments = append(b.req.EnabledPayments, m)
	return b
}

func (b *builder) SetCreditCardDetail(d *snap.CreditCardDetails) *builder {
	b.req.CreditCard = d
	return b
}

func (b *builder) Build() (*snap.Request, error) {
	return b.req, nil
}
