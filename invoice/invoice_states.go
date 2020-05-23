package invoice

type State int

const (
	Draft State = iota
	Published
	WaitForPayment
	Paid
	Failed
)

func (s State) String() string {
	return [...]string{"DRAFT", "PUBLISHED", "WAIT_FOR_PAYMENT", "PAID", "FAILED"}[s]
}

func NewState(state string) StateController {

	switch state {
	case "DRAFT":
		return &DraftState{}
	case "PUBLISHED":
		return &PublishedState{}
	case "WAIT_FOR_PAYMENT":
		return &ProcessedState{}
	case "FAILED":
		return &FailedState{}
	case "PAID":
		return &PaidState{}
	}
	return nil
}

// InvoiceStateController ...
type StateController interface {
	State(i *Invoice) State
	Publish(i *Invoice) error
	Process(i *Invoice) error
	Pay(i *Invoice, transactionID string) error
	Fail(i *Invoice) error
	Reset(i *Invoice) error
}
