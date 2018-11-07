package hl7

import (
	"bytes"
	"testing"
)

func TestReaderReadMessage(t *testing.T) {
	tests := []struct {
		name  string
		data  string
		count int
	}{
		{"empty", "", 0},
		{"one message", "MSH|....", 1},
		{"two messages (CR)", "MSH|....\rMSH|.....", 2},
		{"two messages (CRLF)", "MSH|....\r\nMSH|.....", 2},
		{"two messages (FF)", "MSH|....\fMSH|.....", 2},
		{"one message with trailing delimiter (CR)", "MSH|....\r", 1},
		{"one message with trailing delimiter (CRLF)", "MSH|....\r\n", 1},
		{"one message with trailing delimiter (FF)", "MSH|....\f", 1},
		{"two messages with trailing delimiter (CR)", "MSH|....\rMSH|.....\r", 2},
		{"two messages with trailing delimiter (CRLF)", "MSH|....\r\nMSH|.....\r\n", 2},
		{"two messages with trailing delimiter (FF)", "MSH|....\fMSH|.....\f", 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBufferString(tt.data)
			reader := NewReader(buf)

			for i := 0; i < tt.count; i++ {
				_, err := reader.ReadMessage()

				if err != nil {
					t.Fatalf("Got error `%#v` while reading message %d, want: nil", err, i+1)
				}
			}
			_, err := reader.ReadMessage()

			if err == nil {
				t.Fatal("Did not get error reading message, expected io.EOF")
			}
		})
	}
}
