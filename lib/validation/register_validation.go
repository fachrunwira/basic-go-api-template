package validation

import "github.com/go-playground/validator/v10"

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

func New() *CustomValidator {
	v := validator.New()
	v.RegisterValidation("max_string", MaxStringLength)
	v.RegisterValidation("min_string", MinStringLength)

	return &CustomValidator{Validator: v}
}
