package invoice

type FailedState struct {
}

func (s *FailedState) State(i *Invoice) State {
	return Failed
}

func (s *FailedState) Publish(i *Invoice) error {
	return InvoiceError{InvoiceErrorInvalidStateTransition}
}

func (s *FailedState) Process(i *Invoice) error {
	return InvoiceError{InvoiceErrorInvalidStateTransition}
}

func (s *FailedState) Pay(i *Invoice, transactionID string) error {
	if i.Payment != nil {
		i.Payment.TransactionID = transactionID
	}

	return i.SetState(&PaidState{})
}

func (s *FailedState) Fail(i *Invoice) error {
	return nil
}

func (s *FailedState) Reset(i *Invoice) error {
	i.SetState(&DraftState{})
	return nil
}
