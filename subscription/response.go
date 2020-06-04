package subscription

// CreateResponse is a wrapper for create subscription response
type CreateResponse struct {
	ID                    string
	Status                Status
	LastCreatedInvoiceURL string
}
