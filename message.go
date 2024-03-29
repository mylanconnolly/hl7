package hl7

import (
	"bufio"
	"bytes"
	"io"
	"sync"
	"unicode"
)

// Constants describing possible message boundaries.
const (
	CR = '\r'   // Carriage return
	LF = '\n'   // Line feed
	FF = '\f'   // Form feed
	NB = '\x00' // Null byte
)

// Message is used to describe the parsed message.
type Message struct {
	segments   map[string][]Segment
	reader     *bufio.Reader
	lock       sync.Mutex
	fieldSep   byte
	compSep    byte
	subCompSep byte
	repeat     byte
	escape     byte
}

// Parse is used to parse the segments within the message so that they can be
// queried and iterated. This is a different paradigm from the ReadSegment
// method, which parses the segments as-needed.
func (m *Message) Parse() error {
	m.segments = map[string][]Segment{}

	for {
		segment, err := m.ReadSegment()

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		stype := segment.Type()
		m.segments[stype] = append(m.segments[stype], segment)
	}
	return nil
}

// ReadSegment is used to "read" the next segment from the message.
func (m *Message) ReadSegment() (Segment, error) {
	var buf []byte

	m.lock.Lock()

	for {
		b, err := m.reader.ReadByte()

		if err == io.EOF {
			break
		}

		// Skip all line feeds and character returns while we haven't started saving
		// bytes to the byte slice. This helps cope with messages that have a lot of
		// extra whitespace in them.
		if len(buf) == 0 && unicode.IsSpace(rune(b)) {
			continue
		}

		if b == CR || b == LF {
			break
		}

		buf = append(buf, b)
	}

	m.lock.Unlock()

	if len(buf) == 0 {
		return Segment{}, io.EOF
	}
	return newSegment(m.fieldSep, m.compSep, m.subCompSep, m.repeat, m.escape, buf), nil
}

// NewMessage takes a byte slice and returns a Message that is ready to use.
func NewMessage(data []byte) (*Message, error) {
	// The message must have at least 8 bytes in order to catch all of the
	// character definitions in the header.
	if len(data) < 8 {
		return nil, io.EOF
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
	return &m, nil
}
