package payment

import "fmt"

var (
	// ErrNotFound used if a resource is searched for is missing/not found
	ErrNotFound = fmt.Errorf("not found")
	// ErrInternal used if any internal errors call etc is happened
	ErrInternal = fmt.Errorf("internal error")
	// ErrDatabase used if there is some issue with database
	ErrDatabase = fmt.Errorf("database error")
	// ErrBadRequest used if a user/caller sent the wrong request parameters
	ErrBadRequest = fmt.Errorf("bad request")
	// ErrCantProceed used if a function call is violating some state/logic
	ErrCantProceed = fmt.Errorf("can't proceed")
	// ErrUnauthorized used if a caller is not authenticated
	ErrUnauthorized = fmt.Errorf("unauthorized access")
	// ErrForbidden used if a caller is not authorized to access a resource
	ErrForbidden = fmt.Errorf("forbidden access")
)
