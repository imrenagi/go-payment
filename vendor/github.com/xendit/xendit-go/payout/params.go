package payout

// CreateParams contains parameters for Create
type CreateParams struct {
	ExternalID string  `json:"external_id" validate:"required"`
	Amount     float64 `json:"amount" validate:"required"`
	Email      string  `json:"email" validate:"required"`
}

// GetParams contains parameters for Get
type GetParams struct {
	ID string `json:"id" validate:"required"`
}

// VoidParams contains parameters for Get
type VoidParams struct {
	ID string `json:"id" validate:"required"`
}
