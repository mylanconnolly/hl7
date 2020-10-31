package hl7

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"no format characters", "Hello world", "Hello world"},
		{"with highlighting", `\H\Hello world\N\`, "Hello world"},
		{"with fill mode", `\.fi\Hello world\.nf\`, "Hello world"},
		{"with escaped pipes", `Hello\F\world`, "Hello|world"},
		{"with upcarets", `Hello\S\world`, "Hello^world"},
		{"with ampersands", `Hello\T\world`, "Hello&world"},
		{"with tildes", `Hello\R\world`, "Hello~world"},
		{"with escapes", `Hello\E\world`, `Hello\world`},
		{"with newlines", `Hello\.br\world`, "Hello\nworld"},
		{"with centered newlines", `Hello\.ce\world`, "Hello\nworld"},
		{"with skipped spaces", `Hello\.sk3\world`, "Hello   world"},
		{"with indent", `\.in3\Hello world`, "   Hello world"},
		{"with indent", `\.in\Hello world`, "Hello world"},
		{"with indent", `\.ti3\Hello world`, "   Hello world"},
		{"with indent", `\.ti\Hello world`, "Hello world"},
		{"with newline space", `Hello\.sp3\world`, "Hello\n   world"},
		{"with newline space", `Hello\.sp\world`, "Hello\nworld"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatString(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseRepetition(t *testing.T) {
	tests := []struct {
		name   string
		num    string
		repeat string
		want   string
	}{
		{"repeat spaces 3 times", "3", " ", "   "},
		{"repeat spaces 0 times", "0", " ", ""},
		{"invalid number", "a", " ", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseRepetition(tt.num, tt.repeat)
			assert.Equal(t, tt.want, got)
		})
	}
}
