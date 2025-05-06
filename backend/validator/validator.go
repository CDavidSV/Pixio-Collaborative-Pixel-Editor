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
			ok, err := applyValidationRule(rule, &result, field, v.Type().Field(i).Name)
			if err != nil {
				return result, err
			}

			if !ok {
				return result, nil
			}
		}
	}

	return result, nil
}

func applyValidationRule(rule string, result *ValidationResult, field reflect.Value, fieldName string) (bool, error) {
	switch {
	case strings.HasPrefix(rule, "max="):
		max, err := strconv.ParseFloat(strings.Split(rule, "=")[1], 64)
		if err != nil {
			return false, err
		}

		kind := field.Type().Kind()
		if kind == reflect.String {
			if !MaxChars(field.String(), int(max)) {
				result.IsValid = false
				result.Errors = append(result.Errors, newValidationError(fieldName, field.String(), fmt.Sprintf("%s must be less than %d characters long", fieldName, int(max))))
				return false, nil
			}
		} else if kind >= reflect.Int && kind <= reflect.Float64 {
			value := field.Convert(reflect.TypeOf(float64(0))).Float()
			if value > max {
				result.IsValid = false
				result.Errors = append(result.Errors, newValidationError(fieldName, fmt.Sprintf("%v", value), fmt.Sprintf("%s must be less than %v", fieldName, max)))
				return false, nil
			}
		} else {
			return false, fmt.Errorf("fields using rules \"max\" or \"min\" must be of type string or number")
		}

	case strings.HasPrefix(rule, "min="):
		min, err := strconv.ParseFloat(strings.Split(rule, "=")[1], 64)
		if err != nil {
			return false, err
		}

		kind := field.Type().Kind()
		if kind == reflect.String {
			if !MinChars(field.String(), int(min)) {
				result.IsValid = false
				result.Errors = append(result.Errors, newValidationError(fieldName, field.String(), fmt.Sprintf("%s must be at least %d characters long", fieldName, int(min))))
				return false, nil
			}
		} else if kind >= reflect.Int && kind <= reflect.Float64 {
			value := field.Convert(reflect.TypeOf(float64(0))).Float()
			if value < min {
				result.IsValid = false
				result.Errors = append(result.Errors, newValidationError(fieldName, fmt.Sprintf("%v", value), fmt.Sprintf("%s must be greater than %v", fieldName, min)))
				return false, nil
			}
		} else {
			return false, fmt.Errorf("fields using rules \"max\" or \"min\" must be of type string or number")
		}

	case strings.HasPrefix(rule, "email"):
		if field.Type().Kind() != reflect.String {
			return false, fmt.Errorf("fields using rules \"email\" must be of type string")
		}

		if !IsEmail(field.String()) {
			result.IsValid = false
			result.Errors = append(result.Errors, newValidationError(fieldName, field.String(), fmt.Sprintf("%s must be a valid email address", field.String())))
			return false, nil
		}
	case strings.HasPrefix(rule, "req"):
		if field.Type().Kind() != reflect.String {
			return false, fmt.Errorf("fields using rules \"req\" must be of type string")
		}

		if field.String() == "" {
			result.IsValid = false
			result.Errors = append(result.Errors, newValidationError(fieldName, field.String(), fmt.Sprintf("%s is required", fieldName)))
			return false, nil
		}
	case strings.HasPrefix(rule, "alphanum"):
		if field.Type().Kind() != reflect.String {
			return false, fmt.Errorf("fields using rules \"alphanum\" must be of type string")
		}

		if !IsAlphanumeric(field.String()) {
			result.IsValid = false
			result.Errors = append(result.Errors, newValidationError(fieldName, field.String(), fmt.Sprintf("%s must be alphanumeric", fieldName)))
			return false, nil
		}
	}

	return true, nil
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
