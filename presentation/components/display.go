package components

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/kmdkuk/clicker/application/dto"
)

// DisplayComponent shows basic game information
type Display struct {
}

func NewDisplay() *Display {
	return &Display{}
}

func (d *Display) DrawMoney(screen *ebiten.Image, playerDTO *dto.Player) {
	moneyText := fmt.Sprintf("Money: $%.2f (Total Generate Rate: $%.2f/s)",
		playerDTO.GetMoney(),
		playerDTO.GetTotalGenerateRate())
	ebitenutil.DebugPrintAt(screen, moneyText, 10, 10)
}
