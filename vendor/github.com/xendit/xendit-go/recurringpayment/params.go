package recurringpayment

import (
	"time"

	"github.com/xendit/xendit-go"
)

// CreateParams contains parameters for Create
type CreateParams struct {
	ForUserID           string                              `json:"for_user_id,omitempty"`
	ExternalID          string                              `json:"external_id" validate:"required"`
	PayerEmail          string                              `json:"payer_email" validate:"required"`
	Description         string                              `json:"description" validate:"required"`
	Amount              float64                             `json:"amount" validate:"required"`
	Interval            xendit.RecurringPaymentIntervalEnum `json:"interval" validate:"required"`
	IntervalCount       int                                 `json:"interval_count" validate:"required"`
	TotalRecurrence     int                                 `json:"total_recurrence,omitempty"`
	InvoiceDuration     int                                 `json:"invoice_duration,omitempty"`
	ShouldSendEmail     *bool                               `json:"should_send_email,omitempty"`
	MissedPaymentAction xendit.MissedPaymentActionEnum      `json:"missed_payment_action,omitempty"`
	CreditCardToken     string                              `json:"credit_card_token,omitempty"`
	StartDate           *time.Time                          `json:"start_date,omitempty"`
	SuccessRedirectURL  string                              `json:"success_redirect_url,omitempty"`
	FailureRedirectURL  string                              `json:"failure_redirect_url,omitempty"`
	Recharge            *bool                               `json:"recharge,omitempty"`
	ChargeImmediately   *bool                               `json:"charge_immediately,omitempty"`
}

// GetParams contains parameters for Get
type GetParams struct {
	ID string `json:"id" validate:"required"`
}

// EditParams contains parameters for Edit
type EditParams struct {
	ID                  string                              `json:"-" validate:"required"`
	Amount              float64                             `json:"amount,omitempty"`
	Interval            xendit.RecurringPaymentIntervalEnum `json:"interval,omitempty"`
	IntervalCount       int                                 `json:"interval_count,omitempty"`
	InvoiceDuration     int                                 `json:"invoice_duration,omitempty"`
	ShouldSendEmail     *bool                               `json:"should_send_email,omitempty"`
	MissedPaymentAction xendit.MissedPaymentActionEnum      `json:"missed_payment_action,omitempty"`
	CreditCardToken     string                              `json:"credit_card_token,omitempty"`
}

// StopParams contains parameters for Stop
type StopParams struct {
	ID string `json:"id" validate:"required"`
}

// PauseParams contains parameters for Pause
type PauseParams struct {
	ID string `json:"id" validate:"required"`
}

// ResumeParams contains parameters for Resume
type ResumeParams struct {
	ID string `json:"id" validate:"required"`
}
