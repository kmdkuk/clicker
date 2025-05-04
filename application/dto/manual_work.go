package dto

import "fmt"

type ManualWork struct {
	Name  string
	Value float64
}

func (m *ManualWork) String() string {
	return fmt.Sprintf("%s: $%.2f", m.Name, m.Value)
}
func (m *ManualWork) GetName() string {
	return m.Name
}
