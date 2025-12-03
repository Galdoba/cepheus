package modifier

type Modifier struct {
	category    string
	description string
	value       int
}

func NewModifier(cat, desc string, val int) Modifier {
	return Modifier{category: cat, description: desc, value: val}
}

func (m Modifier) Category() string {
	return m.category
}

func (m Modifier) Description() string {
	return m.description
}

func (m Modifier) Value() int {
	return m.value
}
