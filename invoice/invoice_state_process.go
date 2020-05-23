package invoice

import (
	"time"
)

type ProcessedState struct {
}

func (s *ProcessedState) State(i *Invoice) State {
	now := time.Now()
	if i.DueDate.Before(now) {
		return Failed
	}
	return WaitForPayment
}

func (s *ProcessedState) Publish(i *Invoice) error {
	return InvoiceError{InvoiceErrorInvalidStateTransition}
}

func (s *ProcessedState) Process(i *Invoice) error {
	return nil
}

func (s *ProcessedState) Pay(i *Invoice, transactionID string) error {
	if i.Payment != nil {
		i.Payment.TransactionID = transactionID
	}

	return i.SetState(&PaidState{})
}

func (s *ProcessedState) Fail(i *Invoice) error {
	return i.SetState(&FailedState{})
}

func (s *ProcessedState) Reset(i *Invoice) error {
	if s.State(i) == Failed {
		return i.SetState(&DraftState{})
	}
	return InvoiceError{InvoiceErrorInvalidStateTransition}
}
