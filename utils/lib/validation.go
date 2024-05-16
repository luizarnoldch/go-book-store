package lib

import (
	"fmt"
	appError "main/utils/error"
	"strings"

	"github.com/google/uuid"
)

func ValidateUUID(s string) *appError.Error {
	err := uuid.Validate(s)
	if err != nil {
		return appError.NewValidationError(err.Error())
	}
	return nil
}

func ValidateStringNotEmpty(s string) *appError.Error {
	if strings.TrimSpace(s) == "" {
		return appError.NewValidationError("Name cannot be empty.")
	}
	return nil
}

func ValidateMaxStringCharacteres(s string, max int) *appError.Error {
	if len(s) > max {
		message := fmt.Sprintf("Description cannot exceed %d characters.", max)
		return appError.NewValidationError(message)
	}
	return nil
}

func ValidateMinStringCharacteres(s string, min int) *appError.Error {
	if len(s) < min {
		message := fmt.Sprintf("Description cannot be less thann %d characters.", min)
		return appError.NewValidationError(message)
	}
	return nil
}
