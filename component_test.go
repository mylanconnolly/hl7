package hl7

import (
	"reflect"
	"testing"
)

func TestNewComponent(t *testing.T) {
	tests := []struct {
		name       string
		subCompSep byte
		escape     byte
		data       []byte
		want       Component
	}{
		{"empty (nil)", '&', '\\', []byte(nil), Component(nil)},
		{"empty (not nil)", '&', '\\', []byte{}, Component(nil)},
		{"one part", '&', '\\', []byte("foo"), Component{SubComponent("foo")}},
		{"two parts", '&', '\\', []byte("foo&bar"), Component{SubComponent("foo"), SubComponent("bar")}},
		{"two parts", '@', '\\', []byte("foo@bar"), Component{SubComponent("foo"), SubComponent("bar")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newComponent(tt.subCompSep, tt.escape, tt.data)

			if !reflect.DeepEqual(tt.want, got) {
				t.Fatalf("Got: %#v, want: %#v", got, tt.want)
			}
		})
	}
}
