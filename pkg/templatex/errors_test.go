package templatex

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestNewErrorFormatsKindOpAndMessage(t *testing.T) {
	err := NewError(KindValidation, "templatex.Test", "bad input", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Kind != KindValidation {
		t.Fatalf("expected validation kind, got %q", err.Kind)
	}
	if err.Retryable {
		t.Fatal("expected non-retryable error")
	}
	if got := err.Error(); !strings.Contains(got, "validation: templatex.Test: bad input") {
		t.Fatalf("unexpected error string: %q", got)
	}
}

func TestWrapErrorPreservesCauseAndKind(t *testing.T) {
	cause := context.DeadlineExceeded
	err := WrapError(KindTimeout, "templatex.Test", cause)

	if !IsKind(err, KindTimeout) {
		t.Fatalf("expected timeout kind, got %v", err)
	}
	if !errors.Is(err, cause) {
		t.Fatalf("expected wrapped cause, got %v", err)
	}
	if !err.Retryable {
		t.Fatal("expected retryable error")
	}
}

func TestErrorHandlesNilReceiverAndCauseOnlyMessage(t *testing.T) {
	var nilError *Error
	if got := nilError.Error(); got != "" {
		t.Fatalf("nil Error.Error() = %q; want empty string", got)
	}
	if got := nilError.Unwrap(); got != nil {
		t.Fatalf("nil Error.Unwrap() = %v; want nil", got)
	}

	cause := errors.New("root cause")
	err := &Error{Kind: KindInternal, Op: "templatex.Test", Err: cause}
	if got := err.Error(); got != "internal: templatex.Test: root cause" {
		t.Fatalf("cause-only Error() = %q", got)
	}
}

func TestErrorKindFallbacksForPlainErrors(t *testing.T) {
	err := errors.New("plain")
	if IsKind(err, KindInternal) {
		t.Fatal("plain errors should not match templatex error kinds")
	}
	if got := errorKind(err); got != KindInternal {
		t.Fatalf("errorKind(plain) = %q; want %q", got, KindInternal)
	}
}

func TestContextErrorClassifiesDeadlineAsRetryableTimeout(t *testing.T) {
	err := contextError("templatex.Test", context.DeadlineExceeded)
	if !IsKind(err, KindTimeout) {
		t.Fatalf("expected timeout kind, got %v", err)
	}
	if !err.Retryable {
		t.Fatal("expected deadline errors to be retryable")
	}
}

func TestIsKindMatchesMultipleKinds(t *testing.T) {
	err := NewError(KindConnection, "test", "conn refused", nil)
	if !IsKind(err, KindTimeout, KindConnection) {
		t.Fatal("expected IsKind to match connection in multiple kinds")
	}
	if IsKind(err, KindAuth) {
		t.Fatal("expected IsKind not to match auth")
	}
}

func TestRetryableKindsCorrectness(t *testing.T) {
	retryable := []ErrorKind{KindConnection, KindTimeout, KindUnavailable}
	for _, k := range retryable {
		err := NewError(k, "test", "", nil)
		if !err.Retryable {
			t.Errorf("kind %s should be retryable", k)
		}
	}

	nonRetryable := []ErrorKind{KindValidation, KindConfig, KindAuth, KindClosed, KindInternal}
	for _, k := range nonRetryable {
		err := NewError(k, "test", "", nil)
		if err.Retryable {
			t.Errorf("kind %s should not be retryable", k)
		}
	}
}
