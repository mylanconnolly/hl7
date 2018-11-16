package hl7

import (
	"bufio"
	"bytes"
	"io"
	"sync"
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
	lock   sync.Mutex
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
func (r *Reader) ReadMessage() (*Message, error) {
	var buf []byte

	r.lock.Lock()

	for {
		b, err := r.reader.ReadByte()

		if err == io.EOF {
			break
		}
		if b == CR || b == FF || b == LF {
			p, err := r.reader.Peek(4)

			if err != nil {
				break
			}
			if bytes.Equal(p, []byte("MSH|")) {
				break
			}
			if bytes.Equal(p, []byte("\nMSH")) {
				r.reader.ReadByte()
				break
			}
		}
		buf = append(buf, b)
	}
	r.lock.Unlock()

	return NewMessage(buf)
}
