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

type MessageFunc func(msg *Message) error

// EachMessage is used to create a bit of a friendlier API for reading messages
// if all you're doing is (for example) importing the messages' data into a
// database or something. Each message is passed into the MessageFunc provided
// and executed in order. The downside of this implementation is that it will
// not happen concurrently, which could have some performance ramifications.
//
// Errors returned from this will not include io.EOF, so when you're done
// processing the work, only "real" errors are returned here, such as errors
// parsing the HL7 data and errors reading from the input io.Reader.
func (r *Reader) EachMessage(fn MessageFunc) error {
	for {
		msg, err := r.ReadMessage()

		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if err = fn(msg); err != nil {
			return err
		}
	}
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
		} else if err != nil {
			return nil, err
		}
		// Skip all characters that don't look like they're the beginning of a
		// message until we start storing bytes in the byte slice. This helps us
		// cope with files that have leading whitespace for whatever reason.
		if len(buf) == 0 && b != 'M' {
			continue
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
				// Get rid of the LF
				if _, err = r.reader.ReadByte(); err != nil {
					return nil, err
				}
				break
			}
		}
		buf = append(buf, b)
	}
	r.lock.Unlock()

	return NewMessage(buf)
}
