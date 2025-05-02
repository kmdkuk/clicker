package components

import (
	"log"

	"github.com/kmdkuk/clicker/model"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func ConvertBuildingToListItems(buildings []model.Building) []ListItem {
	items := make([]ListItem, len(buildings))
	for i := range buildings {
		items[i] = &buildings[i]
	}
	return items
}

func ConvertUpgradeToListItems(upgrades []model.Upgrade) []ListItem {
	items := make([]ListItem, len(upgrades))
	for i := range upgrades {
		items[i] = &upgrades[i]
	}
	return items
}

type ListItem interface {
	String(gameState model.GameStateReader) string
}

type List struct {
	gameState model.GameStateReader
	Items     []ListItem
	Visible   bool
	x         int
	y         int
}

func NewList(gameState model.GameStateReader, items []ListItem, defaultVisible bool, x, y int) *List {
	return &List{
		gameState: gameState,
		Items:     items,
		Visible:   defaultVisible,
		x:         x,
		y:         y,
	}
}

func (l *List) Draw(screen *ebiten.Image, cursor int) {
	if l.gameState == nil {
		log.Println("gameState is nil")
		return
	}
	if !l.Visible {
		return
	}

	for i, item := range l.Items {
		y := l.y + i*20
		if i == cursor {
			ebitenutil.DebugPrintAt(screen, "> "+item.String(l.gameState), l.x, y)
		} else {
			ebitenutil.DebugPrintAt(screen, "  "+item.String(l.gameState), l.x, y)
		}
	}
}
