package validator

import "testing"

func TestSlugify(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"Intro to Calculus", "intro-to-calculus"},
		{"  Physics 101!!  ", "physics-101"},
		{"Amharic & English", "amharic-english"},
		{"---weird___title---", "weird-title"},
		{"ግዕዝ", ""}, // no ASCII alphanumerics → empty base
		{"UPPER lower 42", "upper-lower-42"},
	}
	for _, c := range cases {
		if got := Slugify(c.in); got != c.want {
			t.Errorf("Slugify(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestErrorsAccumulate(t *testing.T) {
	e := New()
	if e.HasErrors() {
		t.Fatal("new Errors should be empty")
	}

	e.Required("title", "   ")
	e.MaxLen("title", "ok", 1) // already has a problem for title; first wins
	e.NonNegative("price", -5)
	e.InRange("grade", 99, 1, 12)
	e.OneOf("level", "wizard", "beginner", "intermediate", "advanced")
	e.MinLen("name", "ab", 3)

	if !e.HasErrors() {
		t.Fatal("expected errors")
	}
	m := e.Map()
	for _, field := range []string{"title", "price", "grade", "level", "name"} {
		if _, ok := m[field]; !ok {
			t.Errorf("expected a problem for %q", field)
		}
	}
	// First-message-wins: the Required message should survive the later MaxLen.
	if m["title"] != "is required" {
		t.Errorf("title = %q, want first recorded message", m["title"])
	}
}

func TestErrorsValidInputProducesNilMap(t *testing.T) {
	e := New()
	e.Required("title", "Calculus")
	e.NonNegative("price", 0)
	e.InRange("grade", 12, 1, 12)
	e.OneOf("type", "video", "video", "note", "quiz", "mixed")
	if e.HasErrors() {
		t.Fatalf("expected no errors, got %v", e.Map())
	}
	if e.Map() != nil {
		t.Errorf("Map() = %v, want nil for valid input", e.Map())
	}
}
