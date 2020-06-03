package subscription

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/imrenagi/go-payment"
	"github.com/imrenagi/go-payment/invoice"
)

// New creates empty subscription with valid UUID
func New() *Subscription {
	return &Subscription{
		Number: uuid.New().String(),
		// TODO change this with configuration
		MissedPaymentAction: MissedPaymentActionIgnore,
		Recharge:            true,
		ShouldSendEmail:     true,
		InvoiceDuration:     7 * 24 * time.Hour,
		Invoices:            make([]invoice.Invoice, 0),
		ChargeImmediately:   false,
	}
}

// Subscription is object recording the recurring payment
type Subscription struct {
	payment.Model
	Number              string              `json:"number"`
	Name                string              `json:"name"`
	Description         string              `json:"description"`
	Amount              float64             `json:"amount"`
	UserID              string              `json:"user_id"`
	Currency            string              `json:"currency"`
	Schedule            Schedule            `json:"schedule"`
	TotalReccurence     int                 `json:"total_recurrence"`
	InvoiceDuration     time.Duration       `json:"invoice_duration"`
	ShouldSendEmail     bool                `json:"should_send_email"`
	MissedPaymentAction MissedPaymentAction `json:"missed_payment_action"`
	Recharge            bool                `json:"recharge"`
	CardToken           string              `json:"card_token"`
	GatewayRecurringID  string              `json:"gateway_recurring_id"`
	Gateway             string              `json:"gateway"`
	Invoices            []invoice.Invoice   `json:"invoices"`
	// ChargeImmediately will create first invoice no matter
	// what the startat value is
	ChargeImmediately  bool   `json:"charge_immediately"`
	LastCreatedInvoice string `json:"last_created_invoice"`
	Status             Status `json:"-"`
}

// MarshalJSON ...
func (s *Subscription) MarshalJSON() ([]byte, error) {
	type Alias Subscription

	return json.Marshal(&struct {
		*Alias
		Status string `json:"status"`
	}{
		Alias:  (*Alias)(s),
		Status: s.Status.String(),
	})
}

// Start will create subscription to the payment gateway and update its properties
func (s *Subscription) Start(ctx context.Context, c creator) error {
	res, err := c.Create(ctx, s)
	if err != nil {
		return err
	}

	s.Gateway = c.Gateway().String()
	s.GatewayRecurringID = res.ID
	s.Status = res.Status
	s.LastCreatedInvoice = res.LastCreatedInvoiceURL

	return nil
}

func (s *Subscription) Stop() error {
	return nil
}

func (s *Subscription) Pause() error {
	return nil
}

func (s *Subscription) Resume() error {
	return nil
}

// Charge should create new invoice belong to the subscription
func (s *Subscription) Charge() (*invoice.Invoice, error) {
	return nil, nil
}

// Schedule tells when subscription starts and charges
type Schedule struct {
	payment.Model
	SubscriptionID      uint64       `json:"-" gorm:"index:schedule_subs_id"`
	Interval            int          `json:"interval"`
	IntervalUnit        IntervalUnit `json:"interval_unit"`
	StartAt             *time.Time   `json:"start_at"`
	PreviousExecutionAt *time.Time   `json:"previous_execution_at"`
	NextExecutionAt     *time.Time   `json:"next_execution_at"`
}
