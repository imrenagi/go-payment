package invoice

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"

	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/config"

	"github.com/google/uuid"
)

func newInvoiceNumber() string {
	return uuid.New().String()
}

// NewDefault creates new invoice with 1 day expiration time
func NewDefault() *Invoice {
	now := time.Now()
	inv := New(now, now.AddDate(0, 0, 1))
	return inv
}

// NewWithDurationLimit creates new invoice with any duration
func NewWithDurationLimit(duration time.Duration) *Invoice {
	now := time.Now()
	inv := New(now, now.Add(duration))
	return inv
}

// New accept invoice start and due date and return a new invoice in a draft state.
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
	Title              string          `json:"-"`
	Number             string          `json:"number" gorm:"unique_index:inv_number_k"`
	InvoiceDate        time.Time       `json:"invoice_date"`
	DueDate            time.Time       `json:"due_date"`
	PaidAt             *time.Time      `json:"paid_at"`
	Currency           string          `json:"-"`
	SubTotal           float64         `json:"-"`
	Discount           float64         `json:"-"`
	Tax                float64         `json:"-"`
	ServiceFee         float64         `json:"-"`
	InstallmentFee     float64         `json:"-"`
	State              State           `json:"-"`
	StateController    StateController `json:"-" gorm:"-"`
	LineItems          []LineItem      `json:"items"`
	Payment            *Payment        `json:"payment" gorm:"ForeignKey:InvoiceID"`
	BillingAddress     *BillingAddress `json:"billing_address" gorm:"ForeignKey:InvoiceID"`
	SubscriptionID     *uint64         `json:"-" gorm:"sql:index;"`
	SuccessRedirectURL string          `json:"success_redirect_url"`
	FailureRedirectURL string          `json:"failure_redirect_url"`
}

// GetTitle ...
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

// GetTotal adds up the subtotal, tax, service fee and installment fee and
// deduct the discout value
func (i *Invoice) GetTotal() float64 {
	return i.SubTotal + i.Tax + i.ServiceFee + i.InstallmentFee - i.Discount
}

type paymentMethodFinder interface {
	FindByPaymentType(ctx context.Context, paymentType payment.PaymentType, opts ...payment.Option) (config.FeeConfigReader, error)
}

// Clear remove all invoice values (subtotal, discount, tax, fee, line item and payment method)
// and reset the state to draft
func (i *Invoice) Clear() {
	i.SubTotal = 0
	i.Discount = 0
	i.Tax = 0
	i.ServiceFee = 0
	i.InstallmentFee = 0
	i.State = Draft
	i.LineItems = []LineItem{}
	i.Payment = nil
}

// AfterFind assign a state controller after the entity is fetched from
// database
func (i *Invoice) AfterFind(tx *gorm.DB) error {
	i.StateController = NewState(i.State.String())
	return nil
}

func (i *Invoice) GetStateController() StateController {
	if i.StateController == nil {
		i.StateController = NewState(i.State.String())
	}
	return i.StateController
}

// UpsertBillingAddress create new billing address or update the existing one if it exist
func (i *Invoice) UpsertBillingAddress(name, email, phoneNumber string) error {
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

// UpdatePaymentMethod set the payment method and recalculate the service and installment fee
// of the invoice
func (i *Invoice) UpdatePaymentMethod(ctx context.Context, payment *Payment, finder paymentMethodFinder, opts ...payment.Option) error {

	if payment == nil {
		return InvoiceError{InvoiceErrorPaymentMethodNotSet}
	}

	i.Payment = payment
	i.Payment.InvoiceID = i.ID

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

// CreateChargeRequest will use `charger` to create a charge request to
// the payment gateway and update invoice payment attributes
func (i *Invoice) CreateChargeRequest(ctx context.Context, charger PaymentCharger) error {
	res, err := charger.Create(ctx, i)
	if err != nil {
		return err
	}

	if i.Payment != nil {
		i.Payment.Gateway = charger.Gateway().String()
		i.Payment.Token = res.PaymentToken
		i.Payment.RedirectURL = res.PaymentURL
		i.Payment.TransactionID = res.TransactionID
	}

	return nil
}

// TableName returns the table name used for gorm
func (Invoice) TableName() string {
	return "invoices"
}

// SetItems set the informations of the invoice item
func (i *Invoice) SetItems(ctx context.Context, items []LineItem) error {
	i.LineItems = items
	i.SubTotal = i.GetSubTotal()
	return nil
}

// AddDiscount adds discount if any to the invoice. If discount value is less than 0, error is returned.
func (i *Invoice) AddDiscount(value float64) error {
	if value < 0 {
		return InvoiceError{InvoiceErrorInvalidDiscountValue}
	}
	i.Discount = value
	return nil
}

// RemoveDiscount sets the discount to 0
func (i *Invoice) RemoveDiscount() error {
	i.Discount = 0
	return nil
}

// GetSubTotal returns to total price of all items within the invoice
func (i *Invoice) GetSubTotal() float64 {
	var sum float64
	for _, item := range i.LineItems {
		sum += item.SubTotal()
	}
	return sum
}

// SetState set the invoice state to the given state
func (i *Invoice) SetState(state StateController) error {
	i.StateController = state
	i.State = state.State(i)
	return nil
}

// Publish checks whether payment and billing address of invoice are set.
func (i *Invoice) Publish(ctx context.Context) error {

	if i.Payment == nil {
		return InvoiceError{InvoiceErrorPaymentMethodNotSet}
	}

	if i.BillingAddress == nil {
		return InvoiceError{InvoiceErrorBillingAddressNotSet}
	}

	return i.GetStateController().Publish(i)
}

// Pay set the PaidAt time. It later delegate the action to its state controller.
func (i *Invoice) Pay(ctx context.Context, transactionID string) error {
	now := time.Now()
	i.PaidAt = &now
	return i.GetStateController().Pay(i, transactionID)
}

// Process delegates the action to its state controller
func (i *Invoice) Process(ctx context.Context) error {
	return i.GetStateController().Process(i)
}

// Fail delegates the action to its state controller
func (i *Invoice) Fail(ctx context.Context) error {
	return i.GetStateController().Fail(i)
}

// Reset changes the invoice invoice so that it can be used again
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

	return i.GetStateController().Reset(i)
}

// GetState returns the current state of invoice
func (i *Invoice) GetState() State {
	return i.GetStateController().State(i)
}
