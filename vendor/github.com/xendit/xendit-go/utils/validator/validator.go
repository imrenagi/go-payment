package validator

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	gpv "github.com/go-playground/validator/v10"
	"github.com/xendit/xendit-go"
)

var (
	validator   = gpv.New()
	requiredTag = "required"
)

// ValidateRequired validate fields that are required in the struct with `valiate:"required"` tag
func ValidateRequired(ctx context.Context, value interface{}) error {
	err := validator.StructCtx(ctx, value)
	errs, ok := err.(gpv.ValidationErrors)
	if !ok {
		return err
	}

	missingFields := []string{}
	for _, e := range errs {
		if e.Tag() == requiredTag {
			missingFields = append(missingFields, e.Field())
		}
	}
	return errors.New(errorMsgFromMissingFields(missingFields))
}

// APIValidatorErr generate a xendit.Error from a validator error
func APIValidatorErr(err error) *xendit.Error {
	return &xendit.Error{
		Status:    http.StatusBadRequest,
		ErrorCode: xendit.APIValidationErrCode,
		Message:   err.Error(),
	}
}

func errorMsgFromMissingFields(fields []string) string {
	return fmt.Sprintf("Missing required fields: '%s'", strings.Join(fields, "', '"))
}
