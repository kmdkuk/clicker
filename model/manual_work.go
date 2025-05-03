package model

import "fmt"

type ManualWork struct {
	Name      string  `json:"name"`  // The name displayed for manual work
	BaseValue float64 `json:"value"` // The amount of money earned per manual action
	Count     int     `json:"count"`
}

func (m *ManualWork) String(gameState GameStateReader) string {
	return fmt.Sprintf("%s: $%.2f", m.Name, m.GetValue(gameState.GetUpgrades()))
}

func (m *ManualWork) Work(upgrades []Upgrade) float64 {
	m.Count++
	return m.GetValue(upgrades)
}

func (m *ManualWork) GetValue(upgrades []Upgrade) float64 {
	value := m.BaseValue
	for _, upgrade := range upgrades {
		if upgrade.IsTargetManualWork && upgrade.IsPurchased {
			value = upgrade.Effect(value)
		}
	}
	return value
}
