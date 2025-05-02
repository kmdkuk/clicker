package state

import (
	"errors"

	"github.com/kmdkuk/clicker/level"
)

type Save struct {
	Money      float64
	Buildings  []int
	Upgradings []bool
	ManualWork int
}

func ConverToSave(gameState GameState) Save {
	buildings := make([]int, len(gameState.GetBuildings()))
	upgradings := make([]bool, len(gameState.GetUpgrades()))

	for i, b := range gameState.GetBuildings() {
		buildings[i] = b.Count
	}

	for i, u := range gameState.GetUpgrades() {
		upgradings[i] = u.IsPurchased
	}

	return Save{
		Money:      gameState.GetMoney(),
		Buildings:  buildings,
		Upgradings: upgradings,
		ManualWork: gameState.GetManualWork().Count,
	}
}

func (s *Save) ConvertToGameState() (GameState, error) {
	gameState := NewGameState()
	gameState.UpdateMoney(s.Money)
	if err := gameState.SetManualWorkCount(s.ManualWork); err != nil {
		return gameState, err
	}
	for i, b := range s.Buildings {
		if err := gameState.SetBuildingCount(i, b); err != nil {
			return gameState, err
		}
	}
	for i, u := range s.Upgradings {
		if err := gameState.SetUpgradesIsPurchased(i, u); err != nil {
			return gameState, err
		}
	}
	return gameState, nil
}

func (s *Save) Validation() error {
	if s.Money < 0 {
		return errors.New("invalid money value")
	}
	if len(s.Buildings) > len(level.NewBuildings()) {
		return errors.New("invalid buildings count")
	}
	if len(s.Upgradings) > len(level.NewUpgrades()) {
		return errors.New("invalid upgrades count")
	}
	if s.ManualWork < 0 {
		return errors.New("invalid manual work value")
	}
	return nil
}
