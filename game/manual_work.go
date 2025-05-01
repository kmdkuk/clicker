package game

type ManualWork struct {
	name  string  // Display name
	value float64 // Money earned manually
	count int
}

func (m *ManualWork) String() string {
	return m.name
}

func (m *ManualWork) Work(upgrades []Upgrade) float64 {
	m.count++
	return m.Value(upgrades)
}

func (m *ManualWork) Value(upgrades []Upgrade) float64 {
	value := m.value
	for _, upgrade := range upgrades {
		if upgrade.isTargetManualWork && upgrade.isPurchased {
			value = upgrade.effect(value)
		}
	}
	return value
}
