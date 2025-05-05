package components

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/kmdkuk/clicker/application/dto"
	"github.com/kmdkuk/clicker/presentation/formatter"
)

// DisplayComponent shows basic game information
type Display struct {
}

func NewDisplay() *Display {
	return &Display{}
}

func (d *Display) DrawMoney(screen *ebiten.Image, playerDTO *dto.Player) {
	moneyText := fmt.Sprintf("Money: %s (Total Generate Rate: %s/s)",
		formatter.FormatCurrency(playerDTO.GetMoney(), "$"),
		formatter.FormatCurrency(playerDTO.GetTotalGenerateRate(), "$"),
	)
	ebitenutil.DebugPrintAt(screen, moneyText, 10, 10)
}
