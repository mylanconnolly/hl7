package hl7

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
)

func TestNewMessage(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    Message
		wantErr bool
	}{
		{"Empty (nil)", []byte(nil), Message{}, true},
		{"Empty (not nil)", []byte{}, Message{}, true},
		{"Too short", []byte(`MSH|^~\`), Message{}, true},
		{
			"Minimal example",
			[]byte(`MSH|^~\&`),
			Message{
				reader:     bufio.NewReader(bytes.NewBuffer([]byte(`MSH|^~\&`))),
				fieldSep:   '|',
				compSep:    '^',
				subCompSep: '&',
				repeat:     '~',
				escape:     '\\',
			},
			false,
		},
		{
			"Custom separators",
			[]byte("MSH....."),
			Message{
				reader:     bufio.NewReader(bytes.NewBuffer([]byte("MSH....."))),
				fieldSep:   '.',
				compSep:    '.',
				subCompSep: '.',
				repeat:     '.',
				escape:     '.',
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMessage(tt.data)

			if tt.wantErr && err == nil {
				t.Fatal("Wanted error, received nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("Got error: %#v, wanted nil", err)
			}
			if !reflect.DeepEqual(tt.want, got) {
				t.Fatalf("Got: %#v, want: %#v", got, tt.want)
			}
		})
	}
}

func TestMessageReadSegment(t *testing.T) {
	tests := []struct {
		name  string
		data  []byte
		count int
	}{
		{"one segment", []byte("MSH|^~\\&"), 1},
		{"two segments", []byte("MSH|^~\\&\rMSH|^~\\&"), 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, _ := NewMessage(tt.data)

			for i := 0; i < tt.count; i++ {
				_, err := msg.ReadSegment()

				if err != nil {
					t.Fatalf("Got error `%#v` while reading segment %d, want: nil", err, i+1)
				}
			}
			_, err := msg.ReadSegment()

			if err == nil {
				t.Fatal("Did not get error reading segment, expected io.EOF")
			}
		})
	}
}
