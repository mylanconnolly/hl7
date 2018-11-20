package hl7

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// ErrUnknownTimeFormat is used to represent the case where the time format
// encountered in the HL7 file is unknown. Maybe we don't know how to parse it
// yet, or maybe the HL7 file is not complying with the spec?
var ErrUnknownTimeFormat = errors.New("unknown time format")

// SubComponent is the basic unit in HL7s. This is not strictly standards-
// compliant, since not all fields have sub-components but it is close enough.
type SubComponent []byte

// Int is used to return an integer value housed in a SubComponent.
func (s SubComponent) Int() (int, error) {
	return strconv.Atoi(string(s))
}

// String is used to return the string value housed in a SubComponent. We
// convert HL7 escape codes, linebreaks, etc. into standard values.
func (s SubComponent) String() string {
	out := strings.NewReplacer(
		`\H\`, "",
		`\N\`, "",
		`\.fi\`, "",
		`\.nf\`, "",
		`\F\`, "|",
		`\S\`, "^",
		`\T\`, "&",
		`\R\`, "~",
		`\E\`, `\`,
		`\.br\`, "\n",
		`\.ce\`, "\n",
	).Replace(string(s))

	// TODO: Implement replacers for .sp, .in, .ti, and .sk. See
	// https://www.hermetechnz.com/documentation/UltraPort/MSSQL/index.html?hl7_escape_rules.htm
	// for details. Should we use regex? Or a parser?

	return out
}

// DirtyString is the string value value without any escaping performed.
func (s SubComponent) DirtyString() string {
	return string(s)
}

// Time is used to return a date value housed in a SubComponent.
func (s SubComponent) Time() (time.Time, error) {
	switch len(s) {
	case 8:
		return time.Parse("20060102", string(s))
	case 10:
		return time.Parse("2006010215", string(s))
	case 12:
		return time.Parse("200601021504", string(s))
	case 14:
		return time.Parse("20060102150405", string(s))
	case 16:
		return time.Parse("20060102150405.0", string(s))
	case 17:
		return time.Parse("20060102150405.00", string(s))
	case 18:
		return time.Parse("20060102150405.000", string(s))
	case 19:
		return time.Parse("20060102150405.0000", string(s))
	default:
		return time.Time{}, ErrUnknownTimeFormat
	}
}

func newSubComponent(escape byte, data []byte) SubComponent {
	return SubComponent(data)
}
