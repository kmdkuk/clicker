package game

type ManualWork struct {
	Name  string  // Display name
	Value float64 // Money earned manually
}

func (m *ManualWork) String() string {
	return m.Name
}
