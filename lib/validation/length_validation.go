package validation

import (
	"strconv"

	"github.com/go-playground/validator/v10"
)

// Validate length of string
func MaxStringLength(field validator.FieldLevel) bool {
	fl := field.Field().String()
	param := field.Param()

	maxLen, err := strconv.Atoi(param)
	if err != nil {
		return false
	}

	return len(fl) <= maxLen
}

func MinStringLength(field validator.FieldLevel) bool {
	fl := field.Field().String()
	param := field.Param()

	minLen, err := strconv.Atoi(param)
	if err != nil {
		return false
	}

	return len(fl) >= minLen
}
