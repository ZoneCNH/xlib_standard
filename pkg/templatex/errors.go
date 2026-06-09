package templatex

import (
	"context"
	"errors"
)

// ErrorKind represents the category of an error.
type ErrorKind string

const (
	KindValidation  ErrorKind = "validation"
	KindConfig      ErrorKind = "config"
	KindConnection  ErrorKind = "connection"
	KindAuth        ErrorKind = "auth"
	KindTimeout     ErrorKind = "timeout"
	KindUnavailable ErrorKind = "unavailable"
	KindClosed      ErrorKind = "closed"
	KindInternal    ErrorKind = "internal"
)

// retryableKinds maps kinds that are retryable.
var retryableKinds = map[ErrorKind]bool{
	KindConnection:  true,
	KindTimeout:     true,
	KindUnavailable: true,
}

// Error represents a structured error with kind, operation context, and cause.
type Error struct {
	Kind      ErrorKind
	Op        string
	Message   string
	Err       error
	Retryable bool
}

// NewError creates a new Error. Retryable is derived from kind.
func NewError(kind ErrorKind, op, message string, err error) *Error {
	return &Error{
		Kind:      kind,
		Op:        op,
		Message:   message,
		Err:       err,
		Retryable: retryableKinds[kind],
	}
}

// WrapError wraps an existing error with kind and operation context.
func WrapError(kind ErrorKind, op string, err error) *Error {
	message := ""
	if err != nil {
		message = err.Error()
	}
	return NewError(kind, op, message, err)
}

// Error returns the string representation.
func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	s := string(e.Kind)
	if e.Op != "" {
		s += ": " + e.Op
	}
	if e.Message != "" {
		s += ": " + e.Message
	} else if e.Err != nil {
		s += ": " + e.Err.Error()
	}
	return s
}

// Unwrap returns the underlying error.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// IsKind checks if err is of the given kind(s).
func IsKind(err error, kinds ...ErrorKind) bool {
	var target *Error
	if !errors.As(err, &target) {
		return false
	}
	for _, k := range kinds {
		if target.Kind == k {
			return true
		}
	}
	return false
}

// validationError creates a validation error.
func validationError(op, message string, cause error) *Error {
	return NewError(KindValidation, op, message, cause)
}

// contextError creates an error from a context error.
func contextError(op string, cause error) *Error {
	kind := KindUnavailable
	if errors.Is(cause, context.DeadlineExceeded) {
		kind = KindTimeout
	}
	return NewError(kind, op, "", cause)
}

// errorKind returns the kind of an error, defaulting to KindInternal.
func errorKind(err error) ErrorKind {
	var target *Error
	if errors.As(err, &target) {
		return target.Kind
	}
	return KindInternal
}
