package dto

import (
	"fmt"

	"github.com/kmdkuk/clicker/presentation/formatter"
)

type ManualWork struct {
	Name  string
	Value float64
}

func (m *ManualWork) String() string {
	return fmt.Sprintf("%s: %s", m.Name, formatter.FormatCurrency(m.Value, "$"))
}
func (m *ManualWork) GetName() string {
	return m.Name
}
