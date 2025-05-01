package model

type ManualWork struct {
	Name  string  // Display name
	Value float64 // Money earned manually
	Count int
}

func (m *ManualWork) String() string {
	return m.Name
}

func (m *ManualWork) Work(upgrades []Upgrade) float64 {
	m.Count++
	return m.UpdateValue(upgrades)
}

func (m *ManualWork) UpdateValue(upgrades []Upgrade) float64 {
	value := m.Value
	for _, upgrade := range upgrades {
		if upgrade.IsTargetManualWork && upgrade.IsPurchased {
			value = upgrade.Effect(value)
		}
	}
	return value
}
