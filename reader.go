package hl7

import (
	"bufio"
	"bytes"
	"io"
)

// Constants describing possible message boundaries.
const (
	CR = '\r'
	LF = '\n'
	FF = '\f'
)

// Reader is the type used to read messages from an internal bufio.Reader.
type Reader struct {
	reader *bufio.Reader
}

// NewReader is used to return a new Reader that is ready to use.
func NewReader(reader io.Reader) *Reader {
	r := bufio.NewReader(reader)
	return &Reader{reader: r}
}

// ReadMessage is used to read the next message in the internal reader.
//
// If the reader is empty (or at io.EOF), io.EOF is returned with an empty
// message. Otherwise, error will always be nil.
func (r Reader) ReadMessage() (Message, error) {
	var buf []byte

	for {
		b, err := r.reader.ReadByte()

		if err == io.EOF {
			break
		}
		if b == CR || b == FF {
			p, err := r.reader.Peek(4)

			if err != nil {
				break
			}
			if bytes.Equal(p, []byte("MSH|")) || bytes.Equal(p, []byte("\nMSH")) {
				break
			}
		}
		buf = append(buf, b)
	}
	return NewMessage(buf)
}
