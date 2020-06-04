package subscription_test

import (
	"testing"
	"time"

	. "github.com/imrenagi/go-payment/subscription"
	"github.com/stretchr/testify/assert"
)

func TestSchedule_ScheduleNext(t *testing.T) {

	startAt := time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local)
	nextAt := startAt.AddDate(0, 0, 1)
	nextAfterAt := nextAt.AddDate(0, 0, 1)

	cases := []struct {
		Name             string
		Schedule         Schedule
		ExpectedPrevious time.Time
		ExpectedNext     time.Time
	}{
		{
			Name: "no prev and next",
			Schedule: Schedule{
				Interval:     1,
				IntervalUnit: IntervalUnitDay,
				StartAt:      &startAt,
			},
			ExpectedPrevious: startAt,
			ExpectedNext:     startAt.AddDate(0, 0, 1),
		},
		{
			Name: "has prev and next",
			Schedule: Schedule{
				Interval:            1,
				IntervalUnit:        IntervalUnitDay,
				StartAt:             &startAt,
				PreviousExecutionAt: &startAt,
				NextExecutionAt:     &nextAt,
			},
			ExpectedPrevious: nextAt,
			ExpectedNext:     nextAt.AddDate(0, 0, 1),
		},
		{
			Name: "has prev and next, but prev is different than the start",
			Schedule: Schedule{
				Interval:            1,
				IntervalUnit:        IntervalUnitDay,
				StartAt:             &startAt,
				PreviousExecutionAt: &nextAt,
				NextExecutionAt:     &nextAfterAt,
			},
			ExpectedPrevious: nextAfterAt,
			ExpectedNext:     nextAfterAt.AddDate(0, 0, 1),
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			c.Schedule.ScheduleNext()
			assert.Equal(t, c.ExpectedPrevious, *c.Schedule.PreviousExecutionAt)
			assert.Equal(t, c.ExpectedNext, *c.Schedule.NextExecutionAt)
		})
	}

}
