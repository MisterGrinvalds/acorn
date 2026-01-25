package io

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"gopkg.in/yaml.v3"
)

// Error codes for structured error responses.
const (
	ErrCodeInvalidInput     = "INVALID_INPUT"
	ErrCodeNotFound         = "NOT_FOUND"
	ErrCodePermissionDenied = "PERMISSION_DENIED"
	ErrCodeInternal         = "INTERNAL_ERROR"
	ErrCodeTimeout          = "TIMEOUT"
	ErrCodeValidation       = "VALIDATION_ERROR"
	ErrCodeUnsupported      = "UNSUPPORTED_FORMAT"
)

// IOError represents a structured error response.
type IOError struct {
	Code      string            `json:"code" yaml:"code"`
	Message   string            `json:"message" yaml:"message"`
	Details   map[string]string `json:"details,omitempty" yaml:"details,omitempty"`
	Timestamp time.Time         `json:"timestamp" yaml:"timestamp"`
}

// Error implements the error interface.
func (e *IOError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// WithDetails adds a detail to the error.
func (e *IOError) WithDetails(key, value string) *IOError {
	if e.Details == nil {
		e.Details = make(map[string]string)
	}
	e.Details[key] = value
	return e
}

// NewError creates a new IOError with the given code and message.
func NewError(code, message string) *IOError {
	return &IOError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().UTC(),
	}
}

// NewInvalidInputError creates an invalid input error.
func NewInvalidInputError(message string) *IOError {
	return NewError(ErrCodeInvalidInput, message)
}

// NewNotFoundError creates a not found error.
func NewNotFoundError(message string) *IOError {
	return NewError(ErrCodeNotFound, message)
}

// NewValidationError creates a validation error.
func NewValidationError(message string) *IOError {
	return NewError(ErrCodeValidation, message)
}

// NewInternalError creates an internal error.
func NewInternalError(message string) *IOError {
	return NewError(ErrCodeInternal, message)
}

// ErrorResponse wraps an error for structured output.
type ErrorResponse struct {
	Error *IOError `json:"error" yaml:"error"`
}

// ErrorHandler manages error output in structured formats.
type ErrorHandler struct {
	config *IOConfig
	writer io.Writer
}

// NewErrorHandler creates an ErrorHandler.
func NewErrorHandler(cfg *IOConfig, w io.Writer) *ErrorHandler {
	return &ErrorHandler{
		config: cfg,
		writer: w,
	}
}

// HandleError processes and outputs an error, returning it for command return.
func (h *ErrorHandler) HandleError(err error) error {
	ioErr := h.toIOError(err)

	if h.config.OutputFormat.IsStructured() {
		if writeErr := h.writeStructuredError(ioErr); writeErr != nil {
			// If we can't write structured error, fall back to plain text
			fmt.Fprintf(h.writer, "Error: %s\n", ioErr.Message)
		}
	} else {
		// Plain text error output
		fmt.Fprintf(h.writer, "Error: %s\n", ioErr.Message)
		if len(ioErr.Details) > 0 {
			for k, v := range ioErr.Details {
				fmt.Fprintf(h.writer, "  %s: %s\n", k, v)
			}
		}
	}

	return ioErr
}

// toIOError converts any error to IOError.
func (h *ErrorHandler) toIOError(err error) *IOError {
	if ioErr, ok := err.(*IOError); ok {
		return ioErr
	}

	return &IOError{
		Code:      ErrCodeInternal,
		Message:   err.Error(),
		Timestamp: time.Now().UTC(),
	}
}

// writeStructuredError writes error in JSON/YAML format.
func (h *ErrorHandler) writeStructuredError(err *IOError) error {
	response := ErrorResponse{Error: err}

	switch h.config.OutputFormat {
	case FormatJSON, FormatNDJSON:
		enc := json.NewEncoder(h.writer)
		if h.config.Pretty {
			enc.SetIndent("", "  ")
		}
		return enc.Encode(response)
	case FormatYAML:
		enc := yaml.NewEncoder(h.writer)
		enc.SetIndent(2)
		defer enc.Close()
		return enc.Encode(response)
	default:
		return fmt.Errorf("unsupported format for error output: %s", h.config.OutputFormat)
	}
}

// WrapError wraps a standard error with code if not already an IOError.
func WrapError(err error, code string) *IOError {
	if ioErr, ok := err.(*IOError); ok {
		return ioErr
	}
	return NewError(code, err.Error())
}
