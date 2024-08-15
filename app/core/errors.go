package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"reflect"
	"strings"
)

// ErrorResponse represents the structure of the error response
type ErrorResponse struct {
	Errors map[string]interface{} `json:"errors"`
}

// HandleValidationErrors processes and returns validation errors
func HandleValidationErrors(c *gin.Context, err error) {
	var ve validator.ValidationErrors
	var je *json.UnmarshalTypeError

	switch {
	case errors.As(err, &ve):
		errs := make(map[string]interface{})
		for _, e := range ve {
			errs[e.Field()] = formatErrorMessage(e)
		}
		c.JSON(400, ErrorResponse{Errors: errs})
	case errors.As(err, &je):
		errs := map[string]interface{}{
			je.Field: fmt.Sprintf("Invalid value type. Expected %s", je.Type.String()),
		}
		c.JSON(400, ErrorResponse{Errors: errs})
	default:
		c.JSON(400, ErrorResponse{Errors: map[string]interface{}{"general": err.Error()}})
	}
}

// formatErrorMessage formats a single validation error message
func formatErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email address"
	case "min":
		if e.Type().Kind() == reflect.String {
			return fmt.Sprintf("This field must be at least %s characters long", e.Param())
		}
		return fmt.Sprintf("This field must be at least %s", e.Param())
	case "max":
		if e.Type().Kind() == reflect.String {
			return fmt.Sprintf("This field must be at most %s characters long", e.Param())
		}
		return fmt.Sprintf("This field must be at most %s", e.Param())
	case "e164":
		return "Invalid phone number format"
	case "oneof":
		return fmt.Sprintf("This field must be one of: %s", strings.Replace(e.Param(), " ", ", ", -1))
	case "len":
		return fmt.Sprintf("This field must be exactly %s characters long", e.Param())
	case "numeric":
		return "This field must contain only numeric characters"
	case "alphanum":
		return "This field must contain only alphanumeric characters"
	default:
		return fmt.Sprintf("Invalid value for %s", e.Field())
	}
}

// HTTPError is a custom error type for HTTP errors
type HTTPError struct {
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	return e.Message
}

func CustomErrorResponse(c *gin.Context, err error) {
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		c.JSON(httpErr.StatusCode, gin.H{"error": httpErr.Message})
	}
	log.Println("Error --> ", httpErr.StatusCode, err)
}
