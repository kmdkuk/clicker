package components

import (
	"github.com/kmdkuk/clicker/application/dto"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func ConvertBuildingToListItems(buildings []dto.Building) []ListItem {
	items := make([]ListItem, len(buildings))
	for i := range buildings {
		items[i] = &buildings[i]
	}
	return items
}

func ConvertUpgradeToListItems(upgrades []dto.Upgrade) []ListItem {
	items := make([]ListItem, len(upgrades))
	for i := range upgrades {
		items[i] = &upgrades[i]
	}
	return items
}

type ListItem interface {
	String() string
}

type List struct {
	Items   []ListItem
	Visible bool
	x       int
	y       int
}

func NewList(defaultVisible bool, x, y int) *List {
	return &List{
		Items:   []ListItem{},
		Visible: defaultVisible,
		x:       x,
		y:       y,
	}
}

func (l *List) Draw(screen *ebiten.Image, cursor int) {
	if !l.Visible {
		return
	}

	for i, item := range l.Items {
		y := l.y + i*20
		if i == cursor {
			ebitenutil.DebugPrintAt(screen, "> "+item.String(), l.x, y)
		} else {
			ebitenutil.DebugPrintAt(screen, "  "+item.String(), l.x, y)
		}
	}
}
