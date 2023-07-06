package subscription

import (
	"context"
	"encoding/json"
	"fmt"
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
	}
}

// Subscription is object recording the recurring payment
type Subscription struct {
	payment.Model
	Number              string              `json:"number" gorm:"unique_index:subs_number_k"`
	Name                string              `json:"name"`
	Description         string              `json:"description" gorm:"type:text"`
	Amount              float64             `json:"amount"`
	UserID              string              `json:"user_id"`
	Currency            string              `json:"currency"`
	Schedule            Schedule            `json:"schedule" gorm:"ForeignKey:SubscriptionID"`
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

func (s *Subscription) TableName() string {
	return "goldfish_subscription"
}

// MarshalJSON ...
func (s *Subscription) MarshalJSON() ([]byte, error) {
	type Alias Subscription

	return json.Marshal(&struct {
		*Alias
		Status             string `json:"status"`
		RecurrenceProgress int    `json:"recurrence_progress"`
	}{
		Alias:              (*Alias)(s),
		Status:             s.Status.String(),
		RecurrenceProgress: s.recurrenceProgress(),
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

	s.Schedule.PreviousExecutionAt = s.Schedule.StartAt
	if s.TotalReccurence == 0 || s.TotalReccurence > 1 {
		next := s.Schedule.NextSince(*s.Schedule.PreviousExecutionAt)
		s.Schedule.NextExecutionAt = &next
	}

	return nil
}

func (s Subscription) recurrenceProgress() int {
	return len(s.Invoices)
}

// Pause change the subscription status to paused and stop the schedule
func (s *Subscription) Pause(ctx context.Context, p pauser) error {

	if s.Status != StatusActive {
		return fmt.Errorf("can't pause subscription if it is not in active state: %w", payment.ErrCantProceed)
	}

	if err := p.Pause(ctx, s); err != nil {
		return err
	}
	s.Status = StatusPaused
	return nil
}

// Resume ...
func (s *Subscription) Resume(ctx context.Context, r resumer) error {

	if s.Status != StatusPaused {
		return fmt.Errorf("can't resume subscription if it is not in paused state: %w", payment.ErrCantProceed)
	}

	if err := r.Resume(ctx, s); err != nil {
		return err
	}

	s.Schedule.NextExecutionAt = s.Schedule.NextAfterPause()
	s.Status = StatusActive
	return nil
}

// Stop should stop subscription
func (s *Subscription) Stop(ctx context.Context, st stopper) error {

	if s.Status == StatusStop {
		return fmt.Errorf("subscriptions has been stopped: %w", payment.ErrCantProceed)
	}

	if err := st.Stop(ctx, s); err != nil {
		return err
	}

	s.Schedule.NextExecutionAt = nil
	s.Status = StatusStop
	return nil
}

// Save stores invoice created for subscription and renew subscription
// schedule
func (s *Subscription) Save(inv *invoice.Invoice) error {

	if s.TotalReccurence != 0 && s.recurrenceProgress() >= s.TotalReccurence {
		return fmt.Errorf("should not accept more invoice since all invoices has been recorded %w", payment.ErrCantProceed)
	}

	inv.SubscriptionID = &s.ID
	s.Invoices = append(s.Invoices, *inv)

	if s.Schedule.NextExecutionAt != nil {
		next := s.Schedule.NextSince(*s.Schedule.NextExecutionAt)
		s.Schedule.PreviousExecutionAt = s.Schedule.NextExecutionAt
		s.Schedule.NextExecutionAt = &next
	}
	return nil
}

// NewSchedule create new payment schedule
func NewSchedule(interval int, unit IntervalUnit, start *time.Time) *Schedule {
	s := &Schedule{
		Interval:     interval,
		IntervalUnit: unit,
		StartAt:      start,
	}
	return s
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

func (s *Schedule) TableName() string {
	return "goldfish_schedules"
}

// NextSince ...
func (s *Schedule) NextSince(t time.Time) time.Time {
	return t.Add(time.Duration(s.Interval) * s.IntervalUnit.Duration())
}

// NextAfterPause calculate when the next payment should be executed after it is paused
func (s *Schedule) NextAfterPause() *time.Time {

	// if this schedule is only one time, thus no next charge
	if s.NextExecutionAt == nil {
		return nil
	}

	now := time.Now()
	if s.NextExecutionAt.After(now) {
		return s.NextExecutionAt
	}

	if s.NextExecutionAt.Before(now) {
		var next time.Time
		prev := s.NextExecutionAt
		for {
			next = s.NextSince(*prev)
			if next.After(now) {
				break
			}
			prev = &next
		}
		return &next
	}
	return nil
}
