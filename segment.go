package hl7

// Segment is a slice of fields.
type Segment []Fields

// MarshalHL7 is used to convert the segment back into HL7 format.
func (s Segment) MarshalHL7() []byte {
	buf := []byte{}

	for i, f := range s {
		if i > 0 {
			buf = append(buf, fieldSep)
		}
		buf = append(buf, f.MarshalHL7()...)
	}
	return buf
}

// Type is used to return the type of segment this is.
func (s Segment) Type() string {
	if subComp, ok := s.GetSubComponent(0, 0, 0, 0); ok {
		return subComp.String()
	}
	return ""
}

// GetFields is used to get the fields at the given index
func (s Segment) GetFields(idx int) (Fields, bool) {
	if idx >= len(s) {
		return nil, false
	}
	return s[idx], true
}

// GetField is used to get the field at the given index
func (s Segment) GetField(fieldsIdx, fieldIdx int) (Field, bool) {
	if fields, ok := s.GetFields(fieldsIdx); ok {
		return fields.GetField(fieldIdx)
	}
	return nil, false
}

// GetComponent is used to get the component at the given index
func (s Segment) GetComponent(fieldsIdx, fieldIdx, compIdx int) (Component, bool) {
	if fields, ok := s.GetFields(fieldsIdx); ok {
		return fields.GetComponent(fieldIdx, compIdx)
	}
	return nil, false
}

// GetSubComponent is used to get the sub-component at the given index
func (s Segment) GetSubComponent(fieldsIdx, fieldIdx, compIdx, subCompIdx int) (SubComponent, bool) {
	if fields, ok := s.GetFields(fieldsIdx); ok {
		return fields.GetSubComponent(fieldIdx, compIdx, subCompIdx)
	}
	return nil, false
}

func newSegment(fieldSep, compSep, subCompSep, repeat, escape byte, data []byte) Segment {
	var (
		segment Segment
		start   int
	)
	for i := range data {
		if data[i] == fieldSep {
			segment = append(segment, newFields(compSep, subCompSep, repeat, escape, data[start:i]))
			start = i + 1
		}
		if i == len(data)-1 {
			segment = append(segment, newFields(compSep, subCompSep, repeat, escape, data[start:]))
		}
	}
	return segment
}
