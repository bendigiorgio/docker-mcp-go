package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("validation error: field '%s' %s", v.Field, v.Message)
}

type ValidationResult struct {
	IsValid bool              `json:"is_valid"`
	Errors  []ValidationError `json:"errors,omitempty"`
}

func ValidateRequestParams(params map[string]interface{}, tool *mcp.Tool) error {
	// Check required fields
	for _, requiredField := range tool.InputSchema.Required {
		if _, exists := params[requiredField]; !exists {
			return fmt.Errorf("missing required field: %s", requiredField)
		}
	}

	// Type checking
	for field, schema := range tool.InputSchema.Properties {
		if value, exists := params[field]; exists {
			schemaMap := schema.(map[string]interface{})
			if err := ValidateFieldType(field, value, schemaMap); err != nil {
				return err
			}
		}
	}

	return nil
}

func ValidateFieldType(field string, value interface{}, schema map[string]interface{}) error {
	expectedType := schema["type"].(string)

	switch expectedType {
	case "string":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("field %s must be a string", field)
		}
	case "object":
		if _, ok := value.(map[string]interface{}); !ok {
			return fmt.Errorf("field %s must be an object", field)
		}
	case "array":
		if _, ok := value.([]interface{}); !ok {
			return fmt.Errorf("field %s must be an array", field)
		}
	case "number":
		switch value.(type) {
		case int, int32, int64, float32, float64:
			// Valid number types
		default:
			return fmt.Errorf("field %s must be a number", field)
		}
	case "boolean":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("field %s must be a boolean", field)
		}
	}

	return nil
}

func isEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		return v.Len() == 0
	default:
		return false
	}
}

// FormatValidationError formats validation errors for the standard response
func FormatValidationError(errors []ValidationError) string {
	messages := make([]string, len(errors))
	for i, err := range errors {
		messages[i] = err.Error()
	}
	return strings.Join(messages, "; ")
}
