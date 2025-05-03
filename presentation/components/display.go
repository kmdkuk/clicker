package components

import (
	"fmt"

	"github.com/kmdkuk/clicker/domain/model"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// DisplayComponent shows basic game information
type Display struct {
	gameState model.GameStateReader
}

func NewDisplay(gameState model.GameStateReader) *Display {
	return &Display{
		gameState: gameState,
	}
}

func (d *Display) DrawMoney(screen *ebiten.Image) {
	moneyText := fmt.Sprintf("Money: $%.2f (Total Generate Rate: $%.2f/s)",
		d.gameState.GetMoney(),
		d.gameState.GetTotalGenerateRate())
	ebitenutil.DebugPrintAt(screen, moneyText, 10, 10)
}
