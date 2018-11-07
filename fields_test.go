package hl7

import (
	"reflect"
	"testing"
)

func TestNewFields(t *testing.T) {
	tests := []struct {
		name       string
		repeat     byte
		compSep    byte
		subCompSep byte
		escape     byte
		data       []byte
		want       Fields
	}{
		{"empty (nil)", '~', '^', '&', '\\', []byte(nil), Fields(nil)},
		{"empty (not nil)", '~', '^', '&', '\\', []byte{}, Fields(nil)},
		{"one part", '~', '^', '&', '\\', []byte("foo"), Fields{{{SubComponent("foo")}}}},
		{"two parts", '~', '^', '&', '\\', []byte("foo~bar"), Fields{{{SubComponent("foo")}}, {{SubComponent("bar")}}}},
		{"two parts", '@', '^', '&', '\\', []byte("foo@bar"), Fields{{{SubComponent("foo")}}, {{SubComponent("bar")}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newFields(tt.compSep, tt.subCompSep, tt.repeat, tt.escape, tt.data)

			if !reflect.DeepEqual(tt.want, got) {
				t.Fatalf("Got: %#v, want: %#v", got, tt.want)
			}
		})
	}
}
