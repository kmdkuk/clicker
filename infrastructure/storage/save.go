package storage

import (
	"fmt"

	"github.com/kmdkuk/clicker/game/level"
	"github.com/kmdkuk/clicker/infrastructure/state"
)

type Save struct {
	Money      float64   `json:"money"`
	Buildings  []int     `json:"buildings"`
	Upgradings []upgrade `json:"upgradings"`
	ManualWork int       `json:"manual_work"`
}

type upgrade struct {
	ID          string `json:"id"`
	IsPurchased bool   `json:"is_purchased"`
}

func ConverToSave(gameState state.GameState) Save {
	buildings := make([]int, len(gameState.GetBuildings()))
	upgradings := make([]upgrade, len(gameState.GetUpgrades()))

	for i, b := range gameState.GetBuildings() {
		buildings[i] = b.Count
	}

	for i, u := range gameState.GetUpgrades() {
		upgradings[i].ID = u.ID
		upgradings[i].IsPurchased = u.IsPurchased
	}

	return Save{
		Money:      gameState.GetMoney(),
		Buildings:  buildings,
		Upgradings: upgradings,
		ManualWork: gameState.GetManualWork().Count,
	}
}

func (s *Save) ConvertToGameState() (state.GameState, error) {
	gameState := state.NewGameState()
	gameState.UpdateMoney(s.Money)
	if err := gameState.SetManualWorkCount(s.ManualWork); err != nil {
		return gameState, err
	}
	for i, b := range s.Buildings {
		if err := gameState.SetBuildingCount(i, b); err != nil {
			return gameState, err
		}
	}
	for _, u := range s.Upgradings {
		if err := gameState.SetUpgradesIsPurchasedWithID(u.ID, u.IsPurchased); err != nil {
			return gameState, err
		}
	}
	return gameState, nil
}

func (s *Save) Validation() error {
	if s.Money < 0 {
		return fmt.Errorf("invalid money value: %f", s.Money)
	}
	if len(s.Buildings) > len(level.NewBuildings()) {
		return fmt.Errorf("invalid buildings count: %d", len(s.Buildings))
	}
	if len(s.Upgradings) > len(level.NewUpgrades()) {
		return fmt.Errorf("invalid upgradings count: %d", len(s.Upgradings))
	}
	if s.ManualWork < 0 {
		return fmt.Errorf("invalid manual work count: %d", s.ManualWork)
	}
	return nil
}
