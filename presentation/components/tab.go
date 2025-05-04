package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Tab represents a UI component that displays tabs for different pages
type Tab struct {
	titles     []string // Titles for each tab
	activePage int      // Currently active page index
	x, y       int      // Position of the tabs
}

// NewTab creates a new Tab component
func NewTab(titles []string, defaultPage int, x, y int) *Tab {
	return &Tab{
		titles:     titles,
		activePage: defaultPage,
		x:          x,
		y:          y,
	}
}

// Draw renders the tabs on the screen
func (t *Tab) Draw(screen *ebiten.Image) {
	for i, title := range t.titles {
		xPos := t.x + i*120 // Space tabs 120 pixels apart

		// Highlight the active tab
		if i == t.activePage {
			ebitenutil.DebugPrintAt(screen, "[x] "+title, xPos, t.y)
		} else {
			ebitenutil.DebugPrintAt(screen, "[ ] "+title, xPos, t.y)
		}
	}
}

// SetActivePage changes the active page
func (t *Tab) SetActivePage(page int) {
	if page >= 0 && page < len(t.titles) {
		t.activePage = page
	}
}

// GetActivePage returns the currently active page index
func (t *Tab) GetActivePage() int {
	return t.activePage
}
