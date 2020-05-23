package payment

import "fmt"

// New Errors
var (
	ErrNotFound     = fmt.Errorf("not found")
	ErrInternal     = fmt.Errorf("internal error")
	ErrDatabase     = fmt.Errorf("database error")
	ErrBadRequest   = fmt.Errorf("bad request")
	ErrCantProceed  = fmt.Errorf("can't proceed")
	ErrUnauthorized = fmt.Errorf("unauthorized access")
	ErrForbidden    = fmt.Errorf("forbidden access")
)
