package hl7

// Field is a slice of components.
type Field []Component

// MarshalHL7 is used to convert the field back into HL7 format.
func (f Field) MarshalHL7() []byte {
	buf := []byte{}

	for i, c := range f {
		if i > 0 {
			buf = append(buf, compSep)
		}
		buf = append(buf, c.MarshalHL7()...)
	}
	return buf
}

// GetComponent is used to get the component at the given index.
func (f Field) GetComponent(idx int) (Component, bool) {
	if idx >= len(f) {
		return nil, false
	}
	return f[idx], true
}

// GetSubComponent is used to get the sub-component at the given index.
func (f Field) GetSubComponent(compIdx, subCompIdx int) (SubComponent, bool) {
	if comp, ok := f.GetComponent(compIdx); ok {
		return comp.GetSubComponent(subCompIdx)
	}
	return nil, false
}

func newField(compSep, subCompSep, escape byte, data []byte) Field {
	var (
		field Field
		start int
	)
	for i := range data {
		if data[i] == compSep {
			field = append(field, newComponent(subCompSep, escape, data[start:i]))
			start = i + 1
		}
		if i == len(data)-1 {
			field = append(field, newComponent(subCompSep, escape, data[start:]))
		}
	}
	return field
}
