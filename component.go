package hl7

// Component is used to represent a slice of sub-components.
type Component []SubComponent

// GetSubComponent is used to get the sub-component at the given index.
func (c Component) GetSubComponent(idx int) (SubComponent, bool) {
	if idx >= len(c) {
		return nil, false
	}
	return c[idx], true
}

func newComponent(subCompSep, escape byte, data []byte) Component {
	var (
		comp  Component
		start int
	)

	for i := range data {
		if data[i] == subCompSep {
			comp = append(comp, newSubComponent(escape, data[start:i]))
			start = i + 1
		}
		if i == len(data)-1 {
			comp = append(comp, newSubComponent(escape, data[start:]))
		}
	}
	return comp
}
