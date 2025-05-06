package presentation

import (
	"github.com/kmdkuk/clicker/presentation/input"
)

// NavigationComponent handles cursor position and page management
type Navigation struct {
	cursor     int
	page       int
	maxPages   int
	totalItems []int
}

func NewNavigation(totalItems []int) *Navigation {
	return &Navigation{
		cursor:     0,
		page:       0,
		maxPages:   2, // Total number of pages (currently fixed at 2)
		totalItems: totalItems,
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
	return n.totalItems[n.page] + 1
}

func (n *Navigation) GetCursor() int {
	return n.cursor
}

func (n *Navigation) SetCursor(cursor int) {
	if cursor < 0 {
		n.cursor = 0
	} else if cursor >= n.getTotalItems() {
		n.cursor = n.getTotalItems() - 1
	} else {
		n.cursor = cursor
	}
}

func (n *Navigation) GetPage() int {
	return n.page
}

func (n *Navigation) SetPage(page int) {
	if page < 0 {
		n.page = 0
	} else if page >= n.maxPages {
		n.page = n.maxPages - 1
	} else {
		n.page = page
	}
}
