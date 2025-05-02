package model

type ManualWork struct {
	Name  string  `json:"name"`  // Display name
	Value float64 `json:"value"` // Money earned manually
	Count int     `json:"count"`
}

func (m *ManualWork) String(gameState GameStateReader) string {
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
