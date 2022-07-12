package hl7

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMessage(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    *Message
		wantErr bool
	}{
		{"Empty (nil)", []byte(nil), nil, true},
		{"Empty (not nil)", []byte{}, nil, true},
		{"Too short", []byte(`MSH|^~\`), nil, true},
		{
			"Minimal example",
			[]byte(`MSH|^~\&`),
			&Message{
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
			&Message{
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

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMessageParse(t *testing.T) {
	tests := []struct {
		name   string
		data   []byte
		counts map[string]int
	}{
		{
			"one segment",
			[]byte("MSH|^~\\&"),
			map[string]int{"MSH": 1},
		},
		{
			"two segments",
			[]byte("MSH|^~\\&\rMSH|^~\\&"),
			map[string]int{"MSH": 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, _ := NewMessage(tt.data)
			msg.Parse()

			// Verify that the parsed segments match up with what we expect (this can
			// catch unexpected segments showing up).
			for stype, segments := range msg.segments {
				wantCount := tt.counts[stype]
				assert.Equal(t, wantCount, len(segments))
			}

			// Verify that all the counts we expect match with what was parsed (this
			// can catch missing segments).
			for stype, wantCount := range tt.counts {
				segments := msg.segments[stype]
				assert.Equal(t, wantCount, len(segments))
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
		{"two segments, extra whitespace", []byte("MSH|^~\\&\r\nMSH|^~\\&"), 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, _ := NewMessage(tt.data)

			for i := 0; i < tt.count; i++ {
				_, err := msg.ReadSegment()

				assert.Nil(t, err)
			}
			_, err := msg.ReadSegment()

			assert.Error(t, err)
		})
	}
}
