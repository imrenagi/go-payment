package config

import (
	"fmt"
	"strings"

	"github.com/imrenagi/go-payment"
)

// waitingTime represent the time duration with a time unit (minute, hour, day, second). This used
// for interpreting the waiting time from the payment configuration file.
type waitingTime struct {
	Duration int      `yaml:"duration"`
	Unit     timeUnit `yaml:"unit"`
}

type timeUnit int

const (
	unknownUnit timeUnit = iota
	minute
	hour
	day
	second
)

func (t timeUnit) String() string {
	return []string{"unkown", "minute", "hour", "day", "second"}[t]
}

func (t *timeUnit) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var n string
	if err := unmarshal(&n); err != nil {
		return err
	}
	*t = newTimeUnit(n)
	if *t == unknownUnit {
		return fmt.Errorf("time unit is not recognized, %w", payment.ErrBadRequest)
	}
	return nil
}

func newTimeUnit(name string) timeUnit {
	var t timeUnit
	switch strings.ToLower(name) {
	case "minute":
		t = minute
	case "hour":
		t = hour
	case "day":
		t = day
	case "second":
		t = second
	default:
		t = unknownUnit
	}
	return t
}
