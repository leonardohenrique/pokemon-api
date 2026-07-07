package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidationError representa um erro de validação de um campo específico.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Validate roda as validações de um struct e retorna erros formatados
// e legíveis em vez da mensagem crua da lib.
func Validate(s interface{}) []ValidationError {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var errors []ValidationError
	for _, err := range err.(validator.ValidationErrors) {
		errors = append(errors, ValidationError{
			Field:   strings.ToLower(err.Field()),
			Message: formatMessage(err),
		})
	}
	return errors
}

func formatMessage(fe validator.FieldError) string {
	field := strings.ToLower(fe.Field())

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must have at least %s", field, fe.Param())
	case "max":
		return fmt.Sprintf("%s must have at most %s", field, fe.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, fe.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}
