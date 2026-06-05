// Package validator centralizes field validation for the usecase layer. It
// accumulates field→message problems into a map that maps directly onto
// apperror.Validation's details, so the HTTP layer (Skill 5) can surface
// per-field errors in the standard error envelope (spec §5.4).
package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// Errors collects validation problems keyed by field name. The zero value is
// ready to use. Build one up with the Check* helpers, then branch on HasErrors
// and pass Map() to apperror.Validation.
type Errors struct {
	m map[string]any
}

// New returns an empty Errors.
func New() *Errors {
	return &Errors{m: map[string]any{}}
}

// Add records a problem for field (the first message recorded for a field
// wins, so the most specific check should run first).
func (e *Errors) Add(field, msg string) {
	if e.m == nil {
		e.m = map[string]any{}
	}
	if _, exists := e.m[field]; !exists {
		e.m[field] = msg
	}
}

// HasErrors reports whether any problem was recorded.
func (e *Errors) HasErrors() bool { return len(e.m) > 0 }

// Map returns the accumulated field→message map (nil when empty) for
// apperror.Validation.
func (e *Errors) Map() map[string]any {
	if len(e.m) == 0 {
		return nil
	}
	return e.m
}

// Required records a problem when the trimmed value is empty.
func (e *Errors) Required(field, value string) {
	if strings.TrimSpace(value) == "" {
		e.Add(field, "is required")
	}
}

// MaxLen records a problem when value exceeds n characters (runes).
func (e *Errors) MaxLen(field, value string, n int) {
	if utf8.RuneCountInString(value) > n {
		e.Add(field, "is too long")
	}
}

// MinLen records a problem when value has fewer than n characters (runes).
func (e *Errors) MinLen(field, value string, n int) {
	if utf8.RuneCountInString(value) < n {
		e.Add(field, "is too short")
	}
}

// NonNegative records a problem when value is below zero.
func (e *Errors) NonNegative(field string, value float64) {
	if value < 0 {
		e.Add(field, "must not be negative")
	}
}

// InRange records a problem when value falls outside [min, max].
func (e *Errors) InRange(field string, value, min, max int) {
	if value < min || value > max {
		e.Add(field, "is out of range")
	}
}

// OneOf records a problem when value is not among allowed.
func (e *Errors) OneOf(field, value string, allowed ...string) {
	for _, a := range allowed {
		if value == a {
			return
		}
	}
	e.Add(field, "is not a valid value")
}

// slugStripRe matches every run of characters that is not a lowercase letter or
// digit; such runs collapse to a single dash in Slugify.
var slugStripRe = regexp.MustCompile(`[^a-z0-9]+`)

// Slugify converts a title into a URL-safe slug: lowercase, non-alphanumeric
// runs collapsed to single dashes, leading/trailing dashes trimmed. An input
// with no alphanumeric characters yields "" — callers append a unique suffix
// regardless, so an empty base still produces a usable slug.
func Slugify(title string) string {
	s := strings.ToLower(strings.TrimSpace(title))
	s = slugStripRe.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}
