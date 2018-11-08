package hl7

import (
	"reflect"
	"testing"
	"time"
)

func TestNewSubComponent(t *testing.T) {
	tests := []struct {
		name   string
		data   []byte
		escape byte
		want   SubComponent
	}{
		{"empty (nil)", []byte(nil), '\\', SubComponent(nil)},
		{"empty (not nil)", []byte{}, '\\', SubComponent{}},
		{"not empty", []byte("MSH|..."), '\\', SubComponent("MSH|...")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newSubComponent(tt.escape, tt.data)

			if !reflect.DeepEqual(tt.want, got) {
				t.Fatalf("Got: %#v, want: %#v", got, tt.want)
			}
		})
	}
}

func TestSubComponentInt(t *testing.T) {
	tests := []struct {
		name    string
		value   SubComponent
		want    int
		wantErr bool
	}{
		{"integer", SubComponent("1"), 1, false},
		{"letter", SubComponent("a"), 0, true},
		{"date", SubComponent("20060102"), 20060102, false},
		{"empty", SubComponent(""), 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.value.Int()

			if tt.wantErr && err == nil {
				t.Fatal("Wanted error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("Don't want error, got: %v", err)
			}
			if tt.want != got {
				t.Fatalf("Want: %d, got: %d", tt.want, got)
			}
		})
	}
}

func TestSubComponentString(t *testing.T) {
	tests := []struct {
		name  string
		value SubComponent
		want  string
	}{
		{"integer", SubComponent("1"), "1"},
		{"letter", SubComponent("a"), "a"},
		{"date", SubComponent("20060102"), "20060102"},
		{"empty", SubComponent(""), ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.value.String()

			if tt.want != got {
				t.Fatalf("Want: %s, got: %s", tt.want, got)
			}
		})
	}
}

func TestSubComponentTime(t *testing.T) {
	tests := []struct {
		name    string
		value   SubComponent
		want    time.Time
		wantErr bool
	}{
		{"date", SubComponent("20120505"), time.Date(2012, 5, 5, 0, 0, 0, 0, time.UTC), false},
		{"date with hour", SubComponent("2012050509"), time.Date(2012, 5, 5, 9, 0, 0, 0, time.UTC), false},
		{"date with time", SubComponent("201205050925"), time.Date(2012, 5, 5, 9, 25, 0, 0, time.UTC), false},
		{"date with seconds", SubComponent("20120505092505"), time.Date(2012, 5, 5, 9, 25, 5, 0, time.UTC), false},
		{"date with fractional seconds", SubComponent("20120505092505.1"), time.Date(2012, 5, 5, 9, 25, 5, 100000000, time.UTC), false},
		{"date with fractional seconds", SubComponent("20120505092505.12"), time.Date(2012, 5, 5, 9, 25, 5, 120000000, time.UTC), false},
		{"date with fractional seconds", SubComponent("20120505092505.123"), time.Date(2012, 5, 5, 9, 25, 5, 123000000, time.UTC), false},
		{"date with fractional seconds", SubComponent("20120505092505.1234"), time.Date(2012, 5, 5, 9, 25, 5, 123400000, time.UTC), false},
		{"invalid format", SubComponent("2012-05-05"), time.Time{}, true},
		{"invalid number of characters", SubComponent("2"), time.Time{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.value.Time()

			if tt.wantErr && err == nil {
				t.Fatal("Wanted error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("Don't want error, got: %v", err)
			}
			if !reflect.DeepEqual(tt.want, got) {
				t.Fatalf("Want: %#v, got: %#v", tt.want.Format("2006-01-02 15:04:05.000000000"), got.Format("2006-01-02 15:04:05.000000000"))
			}
		})
	}
}
