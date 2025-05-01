package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type KeyType int

const (
	KeyTypeUp       KeyType = iota // Up
	KeyTypeDown                    // Down
	KeyTypeLeft                    // Left
	KeyTypeRight                   // Right
	KeyTypeDecision                // Decision
	KeyTypeNone                    // No input or other keys
)

type InputHandler struct {
	pressedKey ebiten.Key // Stores the pressed key
}

// Update method to record the pressed key
func (ih *InputHandler) Update() {
	ih.pressedKey = ebiten.Key(0) // Initialize

	// Record the pressed key
	for _, key := range inpututil.AppendJustPressedKeys(nil) {
		ih.pressedKey = key
		break // Record only the first pressed key
	}
}

// GetPressedKey method to classify and retrieve the pressed key
func (ih *InputHandler) GetPressedKey() KeyType {
	switch ih.pressedKey {
	case ebiten.KeyArrowUp, ebiten.KeyW, ebiten.KeyK:
		return KeyTypeUp // Direction key: Up
	case ebiten.KeyArrowDown, ebiten.KeyS, ebiten.KeyJ:
		return KeyTypeDown // Direction key: Down
	case ebiten.KeyArrowLeft, ebiten.KeyA, ebiten.KeyH:
		return KeyTypeLeft // Direction key: Left
	case ebiten.KeyArrowRight, ebiten.KeyD, ebiten.KeyL:
		return KeyTypeRight // Direction key: Right
	case ebiten.KeyEnter, ebiten.KeySpace:
		return KeyTypeDecision // Decision key
	default:
		return KeyTypeNone // No input or other keys
	}
}
