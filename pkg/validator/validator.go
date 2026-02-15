package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"pawnshop/pkg/response"
)

// Validator wraps go-playground/validator
type Validator struct {
	validate *validator.Validate
}

// New creates a new Validator
func New() *Validator {
	v := validator.New()

	// Use json tag for field names
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Register custom validations
	registerCustomValidations(v)

	return &Validator{validate: v}
}

// Validate validates a struct
func (v *Validator) Validate(s interface{}) []response.FieldError {
	err := v.validate.Struct(s)
	if err == nil {
		return nil
	}

	var errors []response.FieldError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors = append(errors, response.FieldError{
				Field:   e.Field(),
				Message: getErrorMessage(e),
			})
		}
	}

	return errors
}

// getErrorMessage returns a human-readable error message
func getErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email address"
	case "min":
		if e.Type().Kind() == reflect.String {
			return "Must be at least " + e.Param() + " characters"
		}
		return "Must be at least " + e.Param()
	case "max":
		if e.Type().Kind() == reflect.String {
			return "Must be at most " + e.Param() + " characters"
		}
		return "Must be at most " + e.Param()
	case "gte":
		return "Must be greater than or equal to " + e.Param()
	case "gt":
		return "Must be greater than " + e.Param()
	case "lte":
		return "Must be less than or equal to " + e.Param()
	case "lt":
		return "Must be less than " + e.Param()
	case "oneof":
		return "Must be one of: " + e.Param()
	case "url":
		return "Invalid URL"
	case "uuid":
		return "Invalid UUID"
	case "numeric":
		return "Must be numeric"
	case "alpha":
		return "Must contain only letters"
	case "alphanum":
		return "Must contain only letters and numbers"
	case "eqfield":
		return "Must match " + e.Param()
	case "nefield":
		return "Must not match " + e.Param()
	case "len":
		return "Must be exactly " + e.Param() + " characters"
	case "dpi":
		return "Invalid DPI format"
	case "phone_gt":
		return "Invalid phone number format"
	default:
		return "Invalid value"
	}
}

// registerCustomValidations registers custom validation rules
func registerCustomValidations(v *validator.Validate) {
	// Guatemala DPI validation (13 digits)
	v.RegisterValidation("dpi", func(fl validator.FieldLevel) bool {
		dpi := fl.Field().String()
		if len(dpi) != 13 {
			return false
		}
		for _, c := range dpi {
			if c < '0' || c > '9' {
				return false
			}
		}
		return true
	})

	// Guatemala phone validation (8 digits starting with 2-7)
	v.RegisterValidation("phone_gt", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		// Remove common separators
		phone = strings.ReplaceAll(phone, "-", "")
		phone = strings.ReplaceAll(phone, " ", "")

		if len(phone) != 8 {
			return false
		}

		first := phone[0]
		if first < '2' || first > '7' {
			return false
		}

		for _, c := range phone {
			if c < '0' || c > '9' {
				return false
			}
		}

		return true
	})

	// Password strength validation
	v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		if len(password) < 8 {
			return false
		}

		hasUpper := false
		hasLower := false
		hasDigit := false

		for _, c := range password {
			switch {
			case 'A' <= c && c <= 'Z':
				hasUpper = true
			case 'a' <= c && c <= 'z':
				hasLower = true
			case '0' <= c && c <= '9':
				hasDigit = true
			}
		}

		return hasUpper && hasLower && hasDigit
	})
}

// Global validator instance
var instance *Validator

// Get returns the global validator instance
func Get() *Validator {
	if instance == nil {
		instance = New()
	}
	return instance
}

// Validate validates using the global instance
func Validate(s interface{}) []response.FieldError {
	return Get().Validate(s)
}
