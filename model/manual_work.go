package model

type ManualWork struct {
	Name  string  `json:"name"`  // The name displayed for manual work
	Value float64 `json:"value"` // The amount of money earned per manual action
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
