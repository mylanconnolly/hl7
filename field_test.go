package hl7

import (
	"reflect"
	"testing"
)

func TestNewField(t *testing.T) {
	tests := []struct {
		name       string
		compSep    byte
		subCompSep byte
		escape     byte
		data       []byte
		want       Field
	}{
		{"empty (nil)", '^', '&', '\\', []byte(nil), Field(nil)},
		{"empty (not nil)", '^', '&', '\\', []byte{}, Field(nil)},
		{"one part", '^', '&', '\\', []byte("foo"), Field{{SubComponent("foo")}}},
		{"two parts", '^', '&', '\\', []byte("foo^bar"), Field{{SubComponent("foo")}, {SubComponent("bar")}}},
		{"two parts", '@', '&', '\\', []byte("foo@bar"), Field{{SubComponent("foo")}, {SubComponent("bar")}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newField(tt.compSep, tt.subCompSep, tt.escape, tt.data)

			if !reflect.DeepEqual(tt.want, got) {
				t.Fatalf("Got: %#v, want: %#v", got, tt.want)
			}
		})
	}
}
