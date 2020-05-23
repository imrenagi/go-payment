package invoice

import (
	"time"
)

type PublishedState struct {
}

func (s PublishedState) State(i *Invoice) State {
	now := time.Now()
	if i.DueDate.Before(now) {
		return Failed
	}
	return Published
}

func (s PublishedState) Publish(i *Invoice) error {
	return i.SetState(&PublishedState{})
}

func (s PublishedState) Process(i *Invoice) error {
	now := time.Now()
	i.InvoiceDate = now
	if i.Payment != nil {
		dur := i.Payment.WaitingDuration()
		if dur != nil {
			i.DueDate = now.Add(*dur)
			return i.SetState(&ProcessedState{})
		}
	}
	return InvoiceError{InvoiceErrorNoPaymentSet}
}

func (s PublishedState) Pay(i *Invoice, transactionID string) error {
	if i.Payment != nil {
		i.Payment.TransactionID = transactionID
	}
	return i.SetState(&PaidState{})
}

func (s PublishedState) Fail(i *Invoice) error {
	return i.SetState(&FailedState{})
}

func (s PublishedState) Reset(i *Invoice) error {
	if s.State(i) == Failed {
		return i.SetState(&DraftState{})
	}
	return InvoiceError{InvoiceErrorInvalidStateTransition}
}
