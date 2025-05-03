package components

import (
	"github.com/kmdkuk/clicker/domain/model"
	"github.com/kmdkuk/clicker/presentation/input"
)

// NavigationComponent handles cursor position and page management
type Navigation struct {
	gameState model.GameStateReader
	cursor    int
	page      int
	maxPages  int
}

func NewNavigation(gameState model.GameStateReader) *Navigation {
	return &Navigation{
		gameState: gameState,
		cursor:    0,
		page:      0,
		maxPages:  2, // Total number of pages (currently fixed at 2)
	}
}

func (n *Navigation) HandleNavigation(keyType input.KeyType) {
	totalItems := n.getTotalItems()

	switch keyType {
	case input.KeyTypeUp:
		n.cursor = (n.cursor - 1 + totalItems) % totalItems
	case input.KeyTypeDown:
		n.cursor = (n.cursor + 1) % totalItems
	case input.KeyTypeLeft:
		n.page = (n.page - 1 + n.maxPages) % n.maxPages
	case input.KeyTypeRight:
		n.page = (n.page + 1) % n.maxPages
	}

	n.validateCursorPosition()
}

func (n *Navigation) validateCursorPosition() {
	totalItems := n.getTotalItems()
	if n.cursor < 0 {
		n.cursor = 0
	} else if n.cursor >= totalItems {
		n.cursor = totalItems - 1
	}
}

func (n *Navigation) getTotalItems() int {
	if n.page == 0 {
		return len(n.gameState.GetBuildings()) + 1 // Manual Work + Buildings
	}
	return len(n.gameState.GetUpgrades()) + 1 // Manual Work + Upgrades
}

func (n *Navigation) GetCursor() int {
	return n.cursor
}

func (n *Navigation) GetPage() int {
	return n.page
}
