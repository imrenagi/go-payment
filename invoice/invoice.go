package invoice

import (
	"context"
	"encoding/json"
	"time"

	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/config"

	"github.com/google/uuid"
)

func newInvoiceNumber() string {
	return uuid.New().String()
}

func NewDefault(message string) *Invoice {
	now := time.Now()
	inv := New(now, now.AddDate(0, 0, 1))
	inv.Message = message

	return inv
}

func New(invoiceDate, dueDate time.Time) *Invoice {
	return &Invoice{
		Number:          newInvoiceNumber(),
		InvoiceDate:     invoiceDate,
		DueDate:         dueDate,
		Currency:        "IDR",
		State:           Draft,
		StateController: &DraftState{},
	}
}

// Invoice ...
type Invoice struct {
	payment.Model
	Title           string          `json:"-"`
	Number          string          `json:"number" gorm:"unique_index:inv_number_k"`
	InvoiceDate     time.Time       `json:"invoice_date"`
	DueDate         time.Time       `json:"due_date"`
	PaidAt          *time.Time      `json:"paid_at"`
	Currency        string          `json:"-"`
	SubTotal        float64         `json:"-"`
	Discount        float64         `json:"-"`
	Tax             float64         `json:"-"`
	ServiceFee      float64         `json:"-"`
	InstallmentFee  float64         `json:"-"`
	State           State           `json:"-"`
	StateController StateController `json:"-" gorm:"-"`
	LineItem        *LineItem       `json:"item"`
	Payment         *Payment        `json:"payment" gorm:"ForeignKey:InvoiceID"`
	BillingAddress  *BillingAddress `json:"billing_address" gorm:"ForeignKey:InvoiceID"`
	Message         string          `json:"message" gorm:"not null;type:text"`
}

func (i *Invoice) GetTitle() string {
	return i.Title
}

// MarshalJSON ...
func (i *Invoice) MarshalJSON() ([]byte, error) {
	type Alias Invoice

	type value struct {
		Currency       string  `json:"currency"`
		Total          float64 `json:"total_amount"`
		SubTotal       float64 `json:"sub_total_amount"`
		Discount       float64 `json:"discount_amount"`
		Tax            float64 `json:"tax_amount"`
		ServiceFee     float64 `json:"admin_fee_amount"`
		InstallmentFee float64 `json:"installment_fee_amount"`
	}

	return json.Marshal(&struct {
		*Alias
		InvTitle string `json:"title"`
		State    string `json:"state"`
		Value    value  `json:"transaction_values"`
	}{
		Alias:    (*Alias)(i),
		InvTitle: i.GetTitle(),
		State:    i.GetState().String(),
		Value: value{
			Currency:       i.Currency,
			Total:          i.GetTotal(),
			SubTotal:       i.SubTotal,
			Discount:       i.Discount,
			Tax:            i.Tax,
			ServiceFee:     i.ServiceFee,
			InstallmentFee: i.InstallmentFee,
		},
	})
}

func (i *Invoice) GetTotal() float64 {
	return i.SubTotal + i.Tax + i.ServiceFee + i.InstallmentFee - i.Discount
}

type paymentMethodFinder interface {
	FindByPaymentType(ctx context.Context, paymentType payment.PaymentType, opts ...payment.PaymentOption) (config.FeeConfigReader, error)
}

// UpdateFee will update service and installment fee when the payment method is set
func (i *Invoice) UpdateFee(ctx context.Context, finder paymentMethodFinder, opts ...payment.PaymentOption) error {
	if i.Payment == nil {
		return InvoiceError{InvoiceErrorPaymentMethodNotSet}
	}

	feeCalculator, err := finder.FindByPaymentType(ctx, i.Payment.PaymentType, opts...)
	if err != nil {
		return err
	}

	i.ServiceFee = 0
	i.InstallmentFee = 0

	if f := feeCalculator.GetAdminFeeConfig(i.Currency); f != nil {
		i.ServiceFee = f.Estimate(i.SubTotal)
	}

	if f := feeCalculator.GetInstallmentFeeConfig(i.Currency); f != nil {
		i.InstallmentFee = f.Estimate(i.SubTotal)
	}
	return nil
}

// Clear ...
func (i *Invoice) Clear() {
	i.SubTotal = 0
	i.Discount = 0
	i.Tax = 0
	i.ServiceFee = 0
	i.InstallmentFee = 0
	i.State = Draft
	i.LineItem = nil
	i.Payment = nil
}

func (i *Invoice) AfterFind() error {
	i.StateController = NewState(i.State.String())
	return nil
}

// SetBillingAddress ...
func (i *Invoice) SetBillingAddress(name, email, phoneNumber string) error {
	if i.BillingAddress == nil {
		addr, err := NewBillingAddress(name, email, phoneNumber)
		if err != nil {
			return err
		}
		i.BillingAddress = addr
	} else {
		err := i.BillingAddress.Update(name, email, phoneNumber)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetPaymentMethod ...
func (i *Invoice) SetPaymentMethod(payment *Payment) error {
	if payment != nil {
		i.Payment = payment
		i.Payment.InvoiceID = i.ID
	}
	return nil
}

// CreateChargeRequest will use `charger` to create a charge request to
// the payment gateway
func (i *Invoice) CreateChargeRequest(ctx context.Context, charger PaymentCharger) error {
	res, err := charger.Create(ctx, i)
	if err != nil {
		return err
	}

	if i.Payment != nil {
		i.Payment.Gateway = charger.Gateway()
		i.Payment.Token = res.PaymentToken
		i.Payment.RedirectURL = res.PaymentURL
		i.Payment.TransactionID = res.TransactionID
	}

	return nil
}

// TableName ...
func (Invoice) TableName() string {
	return "invoices"
}

func (i *Invoice) SetItem(ctx context.Context, item LineItem) error {
	i.LineItem = &LineItem{
		InvoiceID:    i.ID,
		Name:         item.Name,
		Category:     item.Category,
		MerchantName: item.MerchantName,
		Currency:     item.Currency,
		UnitPrice:    item.UnitPrice,
		Qty:          1,
	}

	i.SubTotal = i.LineItem.SubTotal()
	return nil
}

func (i *Invoice) AddDiscount(value float64) error {
	if value < 0 {
		return InvoiceError{InvoiceErrorInvalidDiscountValue}
	}
	i.Discount = value
	return nil
}

func (i *Invoice) RemoveDiscount() error {
	i.Discount = 0
	return nil
}

func (i *Invoice) GetSubTotal() float64 {
	return i.LineItem.SubTotal()
}

func (i *Invoice) SetState(state StateController) error {
	i.StateController = state
	i.State = state.State(i)
	return nil
}

func (i *Invoice) Publish(ctx context.Context) error {

	if i.Payment == nil {
		return InvoiceError{InvoiceErrorPaymentMethodNotSet}
	}

	if i.BillingAddress == nil {
		return InvoiceError{InvoiceErrorBillingAddressNotSet}
	}

	now := time.Now()
	due := now.AddDate(0, 0, 1)
	i.InvoiceDate = now
	i.DueDate = due

	return i.StateController.Publish(i)
}

func (i *Invoice) Pay(ctx context.Context, transactionID string) error {
	now := time.Now()
	i.PaidAt = &now
	return i.StateController.Pay(i, transactionID)
}

func (i *Invoice) Process(ctx context.Context) error {
	return i.StateController.Process(i)
}

func (i *Invoice) Fail(ctx context.Context) error {
	return i.StateController.Fail(i)
}

func (i *Invoice) Reset(ctx context.Context) error {

	i.Number = newInvoiceNumber()
	i.ServiceFee = 0
	i.InstallmentFee = 0

	if i.Payment != nil {
		i.Payment.Reset()
	}

	now := time.Now()
	due := now.Add(i.DueDate.Sub(i.InvoiceDate))
	i.InvoiceDate = now
	i.DueDate = due

	return i.StateController.Reset(i)
}

func (i *Invoice) GetState() State {
	return i.StateController.State(i)
}
