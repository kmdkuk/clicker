package usecase

import (
	"github.com/kmdkuk/clicker/application/dto"
	"github.com/kmdkuk/clicker/infrastructure/state"
)

func NewPlayerUsecase(gameState state.GameState) *PlayerUseCase {
	return &PlayerUseCase{
		gameState: gameState,
	}
}

type PlayerUseCase struct {
	gameState state.GameState
}

func (p *PlayerUseCase) GetPlayer() *dto.Player {
	return &dto.Player{
		Money:             p.gameState.GetMoney(),
		TotalGenerateRate: p.gameState.GetTotalGenerateRate(),
	}
}
