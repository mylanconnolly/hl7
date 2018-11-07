package hl7

import (
	"bufio"
	"bytes"
	"io"
)

// Message is used to describe the parsed message.
type Message struct {
	reader     *bufio.Reader
	fieldSep   byte
	compSep    byte
	subCompSep byte
	repeat     byte
	escape     byte
}

// ReadSegment is used to "read" the next segment from the message.
func (m *Message) ReadSegment() (Segment, error) {
	var buf []byte

	for {
		b, err := m.reader.ReadByte()

		if err == io.EOF {
			break
		}
		if b == CR {
			break
		}
		buf = append(buf, b)
	}
	if len(buf) == 0 {
		return Segment{}, io.EOF
	}
	return newSegment(m.fieldSep, m.compSep, m.subCompSep, m.repeat, m.escape, buf), nil
}

// NewMessage takes a byte slice and returns a Message that is ready to use.
func NewMessage(data []byte) (Message, error) {
	// The message must have at least 8 bytes in order to catch all of the
	// character definitions in the header.
	if len(data) < 8 {
		return Message{}, io.EOF
	}
	reader := bytes.NewBuffer(data)

	m := Message{
		reader:     bufio.NewReader(reader),
		fieldSep:   data[3],
		compSep:    data[4],
		repeat:     data[5],
		escape:     data[6],
		subCompSep: data[7],
	}
	return m, nil
}
