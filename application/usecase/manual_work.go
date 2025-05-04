package usecase

import (
	"log"

	"github.com/kmdkuk/clicker/application/dto"
	"github.com/kmdkuk/clicker/infrastructure/state"
)

func NewManualWorkUseCase(gameState state.GameState) *ManualWorkUseCase {
	return &ManualWorkUseCase{
		gameState: gameState,
	}
}

type ManualWorkUseCase struct {
	gameState state.GameState
}

// GetManualWork implements presentation.ManualWorkUseCase.
func (m *ManualWorkUseCase) GetManualWork() *dto.ManualWork {
	value := m.gameState.GetManualWork().GetValue(m.gameState.GetUpgrades())
	return &dto.ManualWork{
		Name:  m.gameState.GetManualWork().Name,
		Value: value,
	}
}

// ManualWorkAction implements presentation.ManualWorkUseCase.
func (m *ManualWorkUseCase) ManualWorkAction() {
	m.gameState.UpdateMoney(m.gameState.GetManualWork().Work(m.gameState.GetUpgrades()))
	if err := m.gameState.SetManualWorkCount(m.gameState.GetManualWork().Count + 1); err != nil {
		log.Printf("Failed to update manual work count: %v", err)
		return
	}
}
