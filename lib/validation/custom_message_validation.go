package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type (
	errorsValidation struct {
		Field validator.FieldError
	}
)

func Errors(ev validator.ValidationErrors) any {
	for _, err := range ev {
		return errorMessage(&errorsValidation{Field: err})
	}

	return nil
}

func errorMessage(ev *errorsValidation) string {
	switch ev.Field.ActualTag() {
	case "required":
		return fmt.Sprintf("Kolom '%s', masih kosong.", ev.Field.Field())
	case "email":
		return fmt.Sprintf("Kolom '%s', email yang dimasukkan tidak sesuai.", ev.Field.Field())
	case "max_string":
		return fmt.Sprintf("Kolom '%s', panjang karakter maksimal %s karakter.", ev.Field.Field(), ev.Field.Param())
	case "min_string":
		return fmt.Sprintf("Kolom '%s', panjang karakter minimal %s karakter.", ev.Field.Field(), ev.Field.Param())
	case "eqfield":
		return fmt.Sprintf("Kolom '%s', tidak memiliki nilai yang sama pada kolom '%s'.", ev.Field.Field(), ev.Field.Param())
	default:
		return fmt.Sprintf("'%s' invalid.", ev.Field.Field())
	}
}
