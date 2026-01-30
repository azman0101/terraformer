package okta

import (
	"regexp"
	"testing"
)

func TestNormalizeResourceNameWithRandom(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		rand     bool
		expected string // Note: for rand=true, we can't predict exact output but can check prefix/suffix/format
	}{
		{
			name:     "basic string",
			input:    "SimpleName",
			rand:     false,
			expected: "simplename",
		},
		{
			name:     "special chars",
			input:    "Name-With_Special@Chars!",
			rand:     false,
			expected: "name_with_special_chars!",
		},
		{
			name:     "leading numbers",
			input:    "123Name",
			rand:     false,
			expected: "name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeResourceNameWithRandom(tt.input, tt.rand)
			if got != tt.expected {
				t.Errorf("normalizeResourceNameWithRandom(%q, %v) = %q; want %q", tt.input, tt.rand, got, tt.expected)
			}
		})
	}

	// Test rand=true case separately
	t.Run("random suffix", func(t *testing.T) {
		input := "BaseName"
		got := normalizeResourceNameWithRandom(input, true)
		match, _ := regexp.MatchString(`^basename_[a-z0-9]{4}$`, got)
		if !match {
			t.Errorf("normalizeResourceNameWithRandom(%q, true) = %q; did not match expected pattern", input, got)
		}
	})
}

func BenchmarkNormalizeResourceNameWithRandom(b *testing.B) {
	input := "Some-Complex@Name_Structure#123"
	for i := 0; i < b.N; i++ {
		normalizeResourceNameWithRandom(input, false)
	}
}

func BenchmarkNormalizeResourceNameWithRandom_WithRand(b *testing.B) {
	input := "Some-Complex@Name_Structure#123"
	for i := 0; i < b.N; i++ {
		normalizeResourceNameWithRandom(input, true)
	}
}
