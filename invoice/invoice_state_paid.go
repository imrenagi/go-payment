package invoice

type PaidState struct {
}

func (s *PaidState) State(i *Invoice) State {
	return Paid
}

func (s *PaidState) Publish(i *Invoice) error {
	return InvoiceError{InvoiceErrorInvalidStateTransition}
}

func (s *PaidState) Process(i *Invoice) error {
	return InvoiceError{InvoiceErrorInvalidStateTransition}
}

func (s *PaidState) Pay(i *Invoice, transactionID string) error {
	return nil
}

func (s *PaidState) Fail(i *Invoice) error {
	return InvoiceError{InvoiceErrorInvalidStateTransition}
}

func (s *PaidState) Reset(i *Invoice) error {
	return InvoiceError{InvoiceErrorInvalidStateTransition}
}
