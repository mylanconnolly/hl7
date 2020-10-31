package hl7

import (
	"bytes"
	"errors"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
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
				assert.Nil(t, err)
			}

			_, err := reader.ReadMessage()
			assert.Error(t, err)
		})
	}

	t.Run("read error is propagated", func(t *testing.T) {
		buf := bytes.NewBufferString("MSH|....")
		reader := NewReader(iotest.TimeoutReader(buf))
		_, err := reader.ReadMessage()
		assert.Error(t, err)
	})
}

func TestReaderEachMessage(t *testing.T) {
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
			i := 0

			err := reader.EachMessage(func(msg *Message) error {
				i++
				return nil
			})

			assert.Equal(t, tt.count, i)
			assert.Nil(t, err)
		})
	}

	t.Run("read error is propagated", func(t *testing.T) {
		buf := bytes.NewBufferString("MSH|....")
		reader := NewReader(iotest.TimeoutReader(buf))

		err := reader.EachMessage(func(msg *Message) error {
			return nil
		})

		assert.Error(t, err)
	})

	t.Run("error is propagated", func(t *testing.T) {
		buf := bytes.NewBufferString("MSH|....")
		reader := NewReader(buf)

		err := reader.EachMessage(func(msg *Message) error {
			return errors.New("foo")
		})

		assert.Error(t, err)
	})
}
