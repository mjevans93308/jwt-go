package jwt

import (
	"fmt"
	"time"
)

// DefaultValidationHelper is used by Claims.Valid if none is provided
var DefaultValidationHelper = &ValidationHelper{}

// ValidationHelper is built by the parser and passed
// to Claims.Value to carry parse/validation options
type ValidationHelper struct {
	nowFunc func() time.Time
	leeway  time.Duration
}

func (h *ValidationHelper) now() time.Time {
	if h.nowFunc != nil {
		return h.nowFunc()
	}
	return TimeFunc()
}

// Before returns true if Now is before t
// Takes leeway into account
func (h *ValidationHelper) Before(t time.Time) bool {
	return h.now().Before(t.Add(-h.leeway))
}

// After returns true if Now is after t
// Takes leeway into account
func (h *ValidationHelper) After(t time.Time) bool {
	return h.now().After(t.Add(h.leeway))
}

// ValidateExpiresAt returns an error if the expiration time is invalid
// Takes leeway into account
func (h *ValidationHelper) ValidateExpiresAt(exp *Time) error {
	// 'exp' claim is not set. ignore.
	if exp == nil {
		return nil
	}

	// Expiration has passed
	if h.After(exp.Time) {
		delta := h.now().Sub(exp.Time)
		return &ExpiredError{h.now().Unix(), delta}
	}

	// Expiration has not passed
	return nil
}

// ValidateNotBefore returns an error if the nbf time has not been reached
// Takes leeway into account
func (h *ValidationHelper) ValidateNotBefore(nbf *Time) error {
	// 'nbf' claim is not set. ignore.
	if nbf == nil {
		return nil
	}

	// Nbf hasn't been reached
	if h.Before(nbf.Time) {
		return fmt.Errorf("token is not valid yet")
	}
	// Nbf has been reached. valid.
	return nil
}
