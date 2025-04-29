package validator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	FieldName string
	Field     string
	Error     string
}

type ValidationResult struct {
	IsValid bool
	Errors  []ValidationError
}

type ValidationErrorResponse struct {
	Status           int               `json:"status"`
	ValidationErrors []ValidationError `json:"validationErrors"`
}

var EmailRX *regexp.Regexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func (v *ValidationResult) SendValidationError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	response := ValidationErrorResponse{
		Status:           http.StatusBadRequest,
		ValidationErrors: v.Errors,
	}

	responseJSON, _ := json.Marshal(response)
	w.Write(responseJSON)
}

func Validate(obj any) (ValidationResult, error) {
	v := reflect.ValueOf(obj)

	result := ValidationResult{
		IsValid: true,
		Errors:  []ValidationError{},
	}

	// Iterate of the the structs fields and proceed to validate
	for i := range v.NumField() {
		field := v.Field(i)
		tag := v.Type().Field(i).Tag.Get("validate")

		if tag == "" {
			continue
		}

		validationRules := strings.Split(tag, ",")
		for _, rule := range validationRules {
			err := applyValidationRule(rule, &result, field, v.Type().Field(i).Name)
			if err != nil {
				return result, err
			}
		}
	}

	return result, nil
}

func applyValidationRule(rule string, result *ValidationResult, field reflect.Value, fieldName string) error {
	switch {
	case strings.HasPrefix(rule, "max="):
		max, err := strconv.Atoi(strings.Split(rule, "=")[1])
		if err != nil {
			return err
		}

		if field.Type().Kind() != reflect.String {
			return fmt.Errorf("fields using rules \"max\" or \"min\" must be of type string")
		}

		if !MaxChars(field.String(), max) {
			result.IsValid = false
			result.Errors = append(result.Errors, newValidationError(fieldName, field.String(), fmt.Sprintf("%s must be less than %d characters long", fieldName, max)))
		}
	case strings.HasPrefix(rule, "min="):
		min, err := strconv.Atoi(strings.Split(rule, "=")[1])
		if err != nil {
			return err
		}

		if field.Type().Kind() != reflect.String {
			return fmt.Errorf("fields using rules \"max\" or \"min\" must be of type string")
		}

		if !MinChars(field.String(), min) {
			result.IsValid = false
			result.Errors = append(result.Errors, newValidationError(fieldName, field.String(), fmt.Sprintf("%s must be at least %d characters long", fieldName, min)))
		}
	case strings.HasPrefix(rule, "email"):
		if field.Type().Kind() != reflect.String {
			return fmt.Errorf("fields using rules \"email\" must be of type string")
		}

		if !IsEmail(field.String()) {
			result.IsValid = false
			result.Errors = append(result.Errors, newValidationError(fieldName, field.String(), fmt.Sprintf("%s must be a valid email address", field.String())))
		}
	case strings.HasPrefix(rule, "req"):
		if field.Type().Kind() != reflect.String {
			return fmt.Errorf("fields using rules \"req\" must be of type string")
		}

		if field.String() == "" {
			result.IsValid = false
			result.Errors = append(result.Errors, newValidationError(fieldName, field.String(), fmt.Sprintf("%s is required", fieldName)))
		}
	case strings.HasPrefix(rule, "alphanum"):
		if field.Type().Kind() != reflect.String {
			return fmt.Errorf("fields using rules \"alphanum\" must be of type string")
		}

		if !IsAlphanumeric(field.String()) {
			result.IsValid = false
			result.Errors = append(result.Errors, newValidationError(fieldName, field.String(), fmt.Sprintf("%s must be alphanumeric", fieldName)))
		}
	}

	return nil
}

func newValidationError(fieldName, field, errorString string) ValidationError {
	return ValidationError{
		FieldName: fieldName,
		Field:     field,
		Error:     errorString,
	}
}

func MaxChars(text string, maxCount int) bool {
	return len(text) <= maxCount
}

func MinChars(text string, minCount int) bool {
	return len(text) >= minCount
}

func IsEmail(email string) bool {
	return EmailRX.MatchString(email)
}

func IsAlphanumeric(text string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(text)
}
