package hl7

import (
	"reflect"
	"testing"
)

func TestNewSegment(t *testing.T) {
	tests := []struct {
		name       string
		fieldSep   byte
		repeat     byte
		compSep    byte
		subCompSep byte
		escape     byte
		data       []byte
		want       Segment
	}{
		{"empty (nil)", '|', '~', '^', '&', '\\', []byte(nil), Segment(nil)},
		{"empty (not nil)", '|', '~', '^', '&', '\\', []byte{}, Segment(nil)},
		{"one part", '|', '~', '^', '&', '\\', []byte("foo"), Segment{{{{SubComponent("foo")}}}}},
		{"two parts", '|', '~', '^', '&', '\\', []byte("foo|bar"), Segment{{{{SubComponent("foo")}}}, {{{SubComponent("bar")}}}}},
		{"two parts", '@', '~', '^', '&', '\\', []byte("foo@bar"), Segment{{{{SubComponent("foo")}}}, {{{SubComponent("bar")}}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newSegment(tt.fieldSep, tt.compSep, tt.subCompSep, tt.repeat, tt.escape, tt.data)

			if !reflect.DeepEqual(tt.want, got) {
				t.Fatalf("Got: %#v, want: %#v", got, tt.want)
			}
		})
	}
}

func TestSegmentGetFields(t *testing.T) {
	tests := []struct {
		name    string
		segment Segment
		idx     int
		want    Fields
		wantOK  bool
	}{
		{"empty", Segment(nil), 0, nil, false},
		{"invalid index", newSegment('|', '^', '&', '~', '\\', []byte("MSH")), 1, nil, false},
		{"valid index", newSegment('|', '^', '&', '~', '\\', []byte("MSH")), 0, Fields{{{SubComponent("MSH")}}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := tt.segment.GetFields(tt.idx)

			if ok != tt.wantOK {
				t.Fatalf("Got ok: %#v, want: %#v", ok, tt.wantOK)
			}
			if !reflect.DeepEqual(tt.want, got) {
				t.Fatalf("Got: %#v, want: %#v", got, tt.want)
			}
		})
	}
}

func TestSegmentGetField(t *testing.T) {
	tests := []struct {
		name      string
		segment   Segment
		fieldsIdx int
		fieldIdx  int
		want      Field
		wantOK    bool
	}{
		{"empty", Segment(nil), 0, 0, nil, false},
		{"invalid fields index", newSegment('|', '^', '&', '~', '\\', []byte("MSH")), 1, 0, nil, false},
		{"invalid field index", newSegment('|', '^', '&', '~', '\\', []byte("MSH")), 0, 1, nil, false},
		{"valid index", newSegment('|', '^', '&', '~', '\\', []byte("MSH")), 0, 0, Field{{SubComponent("MSH")}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := tt.segment.GetField(tt.fieldsIdx, tt.fieldIdx)

			if ok != tt.wantOK {
				t.Fatalf("Got ok: %#v, want: %#v", ok, tt.wantOK)
			}
			if !reflect.DeepEqual(tt.want, got) {
				t.Fatalf("Got: %#v, want: %#v", got, tt.want)
			}
		})
	}
}

func TestSegmentGetComponent(t *testing.T) {
	tests := []struct {
		name      string
		segment   Segment
		fieldsIdx int
		fieldIdx  int
		compIdx   int
		want      Component
		wantOK    bool
	}{
		{"empty", Segment(nil), 0, 0, 0, nil, false},
		{"invalid fields index", newSegment('|', '^', '&', '~', '\\', []byte("MSH")), 1, 0, 0, nil, false},
		{"invalid field index", newSegment('|', '^', '&', '~', '\\', []byte("MSH")), 0, 1, 0, nil, false},
		{"invalid component index", newSegment('|', '^', '&', '~', '\\', []byte("MSH")), 0, 0, 1, nil, false},
		{"valid index", newSegment('|', '^', '&', '~', '\\', []byte("MSH")), 0, 0, 0, Component{SubComponent("MSH")}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := tt.segment.GetComponent(tt.fieldsIdx, tt.fieldIdx, tt.compIdx)

			if ok != tt.wantOK {
				t.Fatalf("Got ok: %#v, want: %#v", ok, tt.wantOK)
			}
			if !reflect.DeepEqual(tt.want, got) {
				t.Fatalf("Got: %#v, want: %#v", got, tt.want)
			}
		})
	}
}

func TestSegmentGetSubComponent(t *testing.T) {
	tests := []struct {
		name       string
		segment    Segment
		fieldsIdx  int
		fieldIdx   int
		compIdx    int
		subCompIdx int
		want       SubComponent
		wantOK     bool
	}{
		{"empty", Segment(nil), 0, 0, 0, 0, nil, false},
		{"invalid fields index", newSegment('|', '^', '&', '~', '\\', []byte("MSH")), 1, 0, 0, 0, nil, false},
		{"invalid field index", newSegment('|', '^', '&', '~', '\\', []byte("MSH")), 0, 1, 0, 0, nil, false},
		{"invalid component index", newSegment('|', '^', '&', '~', '\\', []byte("MSH")), 0, 0, 1, 0, nil, false},
		{"invalid sub-component index", newSegment('|', '^', '&', '~', '\\', []byte("MSH")), 0, 0, 0, 1, nil, false},
		{"valid index", newSegment('|', '^', '&', '~', '\\', []byte("MSH")), 0, 0, 0, 0, SubComponent("MSH"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := tt.segment.GetSubComponent(tt.fieldsIdx, tt.fieldIdx, tt.compIdx, tt.subCompIdx)

			if ok != tt.wantOK {
				t.Fatalf("Got ok: %#v, want: %#v", ok, tt.wantOK)
			}
			if !reflect.DeepEqual(tt.want, got) {
				t.Fatalf("Got: %#v, want: %#v", got, tt.want)
			}
		})
	}
}
