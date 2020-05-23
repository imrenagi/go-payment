package invoice

type DraftState struct {
}

func (s DraftState) State(i *Invoice) State {
	return Draft
}

func (s DraftState) Publish(i *Invoice) error {
	return i.SetState(&PublishedState{})
}

func (s DraftState) Process(i *Invoice) error {
	return InvoiceError{InvoiceErrorInvalidStateTransition}
}

func (s DraftState) Pay(i *Invoice, transactionID string) error {
	return InvoiceError{InvoiceErrorInvalidStateTransition}
}

func (s DraftState) Fail(i *Invoice) error {
	return InvoiceError{InvoiceErrorInvalidStateTransition}
}

func (s DraftState) Reset(i *Invoice) error {
	return InvoiceError{InvoiceErrorInvalidStateTransition}
}
