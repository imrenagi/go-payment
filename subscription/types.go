package subscription

import (
	"strings"
)

// MissedPaymentAction defines action should be taken
// if user miss the payment
type MissedPaymentAction string

const (
	// MissedPaymentActionIgnore subscription will be continue on the
	// next term
	MissedPaymentActionIgnore MissedPaymentAction = "ignore"
	// MissedPaymentActionStop subscription will be stopped immediately
	MissedPaymentActionStop MissedPaymentAction = "stop"
)

// IntervalUnit defines the unit of recurrence period
type IntervalUnit string

const (
	// IntervalUnitDay recurrence happend every day
	IntervalUnitDay IntervalUnit = "day"
	// IntervalUnitWeek ...
	IntervalUnitWeek IntervalUnit = "week"
	// IntervalUnitMonth ...
	IntervalUnitMonth IntervalUnit = "month"
)

// NewIntervalUnit return an IntervalUnit based on the string given
func NewIntervalUnit(s string) IntervalUnit {
	switch strings.ToLower(s) {
	case "day":
		return IntervalUnitDay
	case "week":
		return IntervalUnitWeek
	default:
		return IntervalUnitMonth
	}
}

// Status shows subscription status
type Status int

const (
	// StatusActive is active
	StatusActive Status = iota
	// StatusPaused subscription may be resumed
	StatusPaused
	// StatusStop subscription cant be resumed
	StatusStop
)

func (s Status) String() string {
	return [...]string{"ACTIVE", "PAUSED", "STOP"}[s]
}
