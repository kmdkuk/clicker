package model

type Upgrade struct {
	Name               string                     `json:"name"`
	Cost               float64                    `json:"cost"`
	Effect             func(float64) float64      `json:"-"` // Exclude from JSON encoding
	IsPurchased        bool                       `json:"is_purchased"`
	IsTargetManualWork bool                       `json:"is_target_manual_work"`
	TargetBuilding     int                        `json:"target_building"`
	IsReleased         func(GameStateReader) bool `json:"-"` // Exclude from JSON encoding
}
