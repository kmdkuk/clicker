package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Handler is an interface for handling input
type Handler interface {
	Update()
	GetPressedKey() KeyType
}

func NewHandler() Handler {
	return &DefaultHandler{
		pressedKey: ebiten.KeyMeta, // Initialize with a default key
	}
}

// DefaultHandler is the default implementation of InputHandler
type DefaultHandler struct {
	pressedKey ebiten.Key // Stores the pressed key
	wheeldx    float64
	wheeldy    float64
}

// Update method to record the pressed key
func (ih *DefaultHandler) Update() {
	ih.pressedKey = ebiten.KeyMeta // Initialize ebiten.Key(0) => 'A'
	// Record the pressed key
	for _, key := range inpututil.AppendJustPressedKeys(nil) {
		ih.pressedKey = key
		break // Record only the first pressed key
	}
	ih.wheeldx, ih.wheeldy = ebiten.Wheel()
}

// GetPressedKey method to classify and retrieve the pressed key
func (ih *DefaultHandler) GetPressedKey() KeyType {
	if ih.wheeldx > 0 {
		return KeyTypeRight
	}
	if ih.wheeldx < 0 {
		return KeyTypeLeft
	}
	if ih.wheeldy > 0 {
		return KeyTypeUp
	}
	if ih.wheeldy < 0 {
		return KeyTypeDown
	}
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
