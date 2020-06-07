package subscription_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/imrenagi/go-payment/invoice"

	"github.com/imrenagi/go-payment"

	. "github.com/imrenagi/go-payment/subscription"
	sm "github.com/imrenagi/go-payment/subscription/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSubscription_Start(t *testing.T) {

	c := &sm.Controller{}

	c.On("Create", mock.Anything, mock.Anything).
		Return(&CreateResponse{
			ID:                    "12345",
			Status:                StatusActive,
			LastCreatedInvoiceURL: "http://example.com",
		}, nil)
	c.On("Gateway").Return(payment.GatewayXendit)

	t.Run("crete subscription with multiple recurrence", func(t *testing.T) {

		s := New()
		startAt := time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local)
		s.Schedule = *NewSchedule(1, IntervalUnitDay, &startAt)

		err := s.Start(context.TODO(), c)
		assert.Nil(t, err)
		assert.Equal(t, payment.GatewayXendit.String(), s.Gateway)
		assert.Equal(t, "12345", s.GatewayRecurringID)
		assert.Equal(t, StatusActive, s.Status)
		assert.Equal(t, "http://example.com", s.LastCreatedInvoice)

		assert.Equal(t, startAt, *s.Schedule.PreviousExecutionAt)
		assert.Equal(t, startAt.AddDate(0, 0, 1), *s.Schedule.NextExecutionAt)

	})

	t.Run("crete subscription with one recurrence", func(t *testing.T) {

		s := New()
		s.TotalReccurence = 1
		startAt := time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local)
		s.Schedule = *NewSchedule(1, IntervalUnitDay, &startAt)

		err := s.Start(context.TODO(), c)
		assert.Nil(t, err)
		assert.Equal(t, payment.GatewayXendit.String(), s.Gateway)
		assert.Equal(t, "12345", s.GatewayRecurringID)
		assert.Equal(t, StatusActive, s.Status)
		assert.Equal(t, "http://example.com", s.LastCreatedInvoice)

		assert.Equal(t, startAt, *s.Schedule.PreviousExecutionAt)
		assert.Nil(t, s.Schedule.NextExecutionAt)

	})

}

func TestSubscription_Pause(t *testing.T) {

	t.Run("successfully pause", func(t *testing.T) {

		c := &sm.Controller{}
		c.On("Pause", mock.Anything, mock.Anything).Return(nil)

		s := New()
		startAt := time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local)
		s.Schedule = *NewSchedule(1, IntervalUnitDay, &startAt)
		s.Status = StatusActive

		err := s.Pause(context.TODO(), c)
		assert.Nil(t, err)

		assert.Equal(t, StatusPaused, s.Status)
	})

	t.Run("cant pause if it is nnot in active", func(t *testing.T) {
		c := &sm.Controller{}
		s := New()
		startAt := time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local)
		s.Schedule = *NewSchedule(1, IntervalUnitDay, &startAt)
		s.Status = StatusPaused

		err := s.Pause(context.TODO(), c)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, payment.ErrCantProceed))
		assert.Equal(t, StatusPaused, s.Status)
	})

}

func TestSubscription_Resume(t *testing.T) {

	t.Run("can't resume if it is not paused", func(t *testing.T) {

		c := &sm.Controller{}
		c.On("Resume", mock.Anything, mock.Anything).Return(nil)

		s := New()
		s.Status = StatusStop
		startAt := time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local)
		s.Schedule = *NewSchedule(1, IntervalUnitDay, &startAt)

		err := s.Resume(context.TODO(), c)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, payment.ErrCantProceed))
		assert.Equal(t, StatusStop, s.Status)
	})

	t.Run("successfully resume", func(t *testing.T) {

		c := &sm.Controller{}
		c.On("Resume", mock.Anything, mock.Anything).Return(nil)

		s := New()
		s.Status = StatusPaused
		startAt := time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local)
		s.Schedule = *NewSchedule(1, IntervalUnitDay, &startAt)

		next := time.Now().Add(-1 * time.Hour)
		s.Schedule.NextExecutionAt = &next

		err := s.Resume(context.TODO(), c)
		assert.Nil(t, err)

		expected := next.AddDate(0, 0, 1)
		assert.Equal(t, StatusActive, s.Status)
		assert.Equal(t, expected.Second(), s.Schedule.NextExecutionAt.Second())
	})

}

func TestSubscription_Stop(t *testing.T) {
	t.Run("successfully stop", func(t *testing.T) {

		c := &sm.Controller{}
		c.On("Stop", mock.Anything, mock.Anything).Return(nil)

		s := New()
		s.Status = StatusActive
		startAt := time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local)
		s.Schedule = *NewSchedule(1, IntervalUnitDay, &startAt)

		next := time.Now().Add(-1 * time.Hour)
		s.Schedule.NextExecutionAt = &next

		assert.Equal(t, StatusActive, s.Status)
		assert.NotNil(t, s.Schedule.NextExecutionAt)

		err := s.Stop(context.TODO(), c)
		assert.Nil(t, err)

		assert.Equal(t, StatusStop, s.Status)
		assert.Nil(t, s.Schedule.NextExecutionAt)
	})
}

func TestSubscription_Save(t *testing.T) {

	t.Run("one time subscription, should save invoice if no invoice is paid/expired", func(t *testing.T) {
		s := New()
		s.TotalReccurence = 1
		startAt := time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local)
		s.Schedule = *NewSchedule(1, IntervalUnitDay, &startAt)
		s.Schedule.PreviousExecutionAt = &startAt

		next := time.Now().Add(-1 * time.Hour)
		s.Schedule.NextExecutionAt = &next

		assert.Len(t, s.Invoices, 0)

		now := time.Now()
		inv := invoice.New(now, now.Add(1*time.Hour))

		err := s.Save(inv)
		assert.Nil(t, err)
		assert.Len(t, s.Invoices, 1)

		assert.Equal(t, next.AddDate(0, 0, 1).Second(), s.Schedule.NextExecutionAt.Second())
	})

	t.Run("one time subscription, should not save invoice if paid before", func(t *testing.T) {
		s := New()
		s.TotalReccurence = 1
		assert.Len(t, s.Invoices, 0)

		now := time.Now()
		inv := invoice.New(now, now.Add(1*time.Hour))
		err := s.Save(inv)
		assert.Nil(t, err)

		err = s.Save(inv)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, payment.ErrCantProceed))

		assert.Len(t, s.Invoices, 1)
	})

	t.Run("save multiple invoice for recurring payment", func(t *testing.T) {

		s := New()
		assert.Len(t, s.Invoices, 0)

		now := time.Now()
		inv := invoice.New(now, now.Add(1*time.Hour))

		err := s.Save(inv)
		assert.Nil(t, err)
		assert.Len(t, s.Invoices, 1)

		err = s.Save(inv)
		assert.Nil(t, err)
		assert.Len(t, s.Invoices, 2)

	})

}

func TestSchedule_NextAfterPause(t *testing.T) {

	startAt := time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local)

	t.Run("one time schedule should has no next execution date", func(t *testing.T) {
		s := NewSchedule(1, IntervalUnitDay, &startAt)
		assert.Nil(t, s.NextAfterPause())
	})

	t.Run("next execution date is after now, return next execution date", func(t *testing.T) {
		nextAt := time.Now().AddDate(0, 0, 1)
		s := NewSchedule(1, IntervalUnitDay, &startAt)
		s.NextExecutionAt = &nextAt
		assert.Equal(t, nextAt, *s.NextAfterPause())
	})

	t.Run("next execution date is passed once, find the next execution date after now", func(t *testing.T) {
		now := time.Now()
		s := NewSchedule(1, IntervalUnitMonth, &now)
		start := now.AddDate(0, -2, 0)
		next := now.AddDate(0, -1, 0).Add(-5 * time.Hour)

		s.StartAt = &start
		s.NextExecutionAt = &next

		expected := now.Add(-5 * time.Hour)

		assert.Equal(t, expected.Second(), s.NextAfterPause().Second())
	})

	t.Run("next execution date is passed multiple times, find the next execution date after now", func(t *testing.T) {
		now := time.Now()
		s := NewSchedule(1, IntervalUnitMonth, &now)
		start := now.AddDate(0, -5, 0)
		next := now.AddDate(0, -4, 0)

		s.StartAt = &start
		s.NextExecutionAt = &next

		result := s.NextAfterPause()
		assert.True(t, result.After(now))
	})

}
