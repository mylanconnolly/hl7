package hl7

// Field is a slice of components.
type Field []Component

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
