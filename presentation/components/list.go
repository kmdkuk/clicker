package components

import (
	"image/color"

	"github.com/kmdkuk/clicker/application/dto"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
	Items        []ListItem
	Visible      bool
	x            int
	y            int
	scrollPos    int // Current scroll position
	viewportSize int // Number of items visible at once
}

func NewList(defaultVisible bool, x, y int) *List {
	// Default to showing 10 items at once
	return NewListWithViewport(defaultVisible, x, y, 10)
}

// NewListWithViewport creates a new list with a specified viewport size
func NewListWithViewport(defaultVisible bool, x, y, viewportSize int) *List {
	return &List{
		Items:        []ListItem{},
		Visible:      defaultVisible,
		x:            x,
		y:            y,
		scrollPos:    0,
		viewportSize: viewportSize,
	}
}

func (l *List) Draw(screen *ebiten.Image, cursor int) {
	if !l.Visible {
		return
	}

	// Don't render anything for empty lists
	if len(l.Items) == 0 {
		return
	}

	// Adjust scroll position if cursor moves outside the viewport
	if cursor < l.scrollPos {
		l.scrollPos = cursor
	} else if cursor >= l.scrollPos+l.viewportSize {
		l.scrollPos = cursor - l.viewportSize + 1
	}

	// Ensure scroll position stays within valid range
	if l.scrollPos < 0 {
		l.scrollPos = 0
	}
	maxScroll := len(l.Items) - l.viewportSize
	if maxScroll < 0 {
		maxScroll = 0
	}
	if l.scrollPos > maxScroll {
		l.scrollPos = maxScroll
	}

	// Calculate ending index for display range
	endIdx := l.scrollPos + l.viewportSize
	if endIdx > len(l.Items) {
		endIdx = len(l.Items)
	}

	// Draw only items within the current viewport
	for i := l.scrollPos; i < endIdx; i++ {
		item := l.Items[i]
		displayY := l.y + (i-l.scrollPos)*20

		if i == cursor {
			ebitenutil.DebugPrintAt(screen, "> "+item.String(), l.x, displayY)
		} else {
			ebitenutil.DebugPrintAt(screen, "  "+item.String(), l.x, displayY)
		}
	}

	// Draw scrollbar instead of indicators
	l.drawScrollBar(screen, endIdx)
}

// drawScrollBar draws a scrollbar on the right edge of the screen
func (l *List) drawScrollBar(screen *ebiten.Image, endIdx int) {
	// Get screen dimensions
	screenWidth := screen.Bounds().Dx()

	// Skip drawing scrollbar if all items fit in viewport
	if len(l.Items) <= l.viewportSize {
		return
	}

	// Constants for scrollbar styling
	const (
		scrollbarWidth  = 8
		scrollbarMargin = 5
	)

	// Calculate scrollbar position
	scrollbarX := float64(screenWidth - scrollbarWidth - scrollbarMargin)
	scrollbarY := float64(l.y)

	// Calculate scrollbar height based on visible range
	visibleCount := endIdx - l.scrollPos
	listHeight := float64(visibleCount * 20)

	// Draw scrollbar background (track)
	backgroundColor := color.RGBA{R: 80, G: 80, B: 80, A: 180}
	vector.DrawFilledRect(screen, float32(scrollbarX), float32(scrollbarY),
		float32(scrollbarWidth), float32(listHeight), backgroundColor, false)

	// Calculate handle size and position
	// Handle size is proportional to the visible portion of the list
	totalItems := len(l.Items)
	handleRatio := float64(visibleCount) / float64(totalItems)
	handleHeight := listHeight * handleRatio
	if handleHeight < 10 {
		handleHeight = 10 // Minimum handle height
	}

	// Position the handle based on scroll position
	scrollRatio := float64(l.scrollPos) / float64(totalItems-visibleCount)
	handleY := scrollbarY + scrollRatio*(listHeight-handleHeight)

	// Draw scrollbar handle
	handleColor := color.RGBA{R: 180, G: 180, B: 180, A: 255}
	vector.DrawFilledRect(screen, float32(scrollbarX), float32(handleY),
		float32(scrollbarWidth), float32(handleHeight), handleColor, false)
}

// Scroll manually adjusts the scroll position by the specified amount
func (l *List) Scroll(amount int) {
	l.scrollPos += amount

	// Prevent scrolling out of bounds
	if l.scrollPos < 0 {
		l.scrollPos = 0
	}

	maxScroll := len(l.Items) - l.viewportSize
	if maxScroll < 0 {
		maxScroll = 0
	}
	if l.scrollPos > maxScroll {
		l.scrollPos = maxScroll
	}
}

// GetVisibleRange returns the start and end indices of currently visible items (for testing)
func (l *List) GetVisibleRange() (start, end int) {
	end = l.scrollPos + l.viewportSize
	if end > len(l.Items) {
		end = len(l.Items)
	}
	return l.scrollPos, end
}
