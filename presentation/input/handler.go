package input

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Handler is an interface for handling input
type Handler interface {
	Update()
	GetPressedKey() KeyType
	IsClicked() bool
	GetMouseCursor() (int, int)
	ResetClickState()
}

func NewHandler() Handler {
	return &DefaultHandler{
		pressedKey: ebiten.KeyMeta, // Initialize with a default key
	}
}

// DefaultHandler is the default implementation of InputHandler
type DefaultHandler struct {
	pressedKey   ebiten.Key // Stores the pressed key
	wheeldx      float64
	wheeldy      float64
	mouseX       int
	mouseY       int
	keepClicking bool
	isClicked    bool
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

	if !ih.keepClicking && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		ih.keepClicking = true
		ih.isClicked = true
		ih.mouseX, ih.mouseY = ebiten.CursorPosition()
		fmt.Printf("Mouse clicked at (%d, %d)\n", ih.mouseX, ih.mouseY)
	}
	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		ih.keepClicking = false
	}
}

func (ih *DefaultHandler) ResetClickState() {
	ih.isClicked = false
	ih.mouseX = 0
	ih.mouseY = 0
}

func (ih *DefaultHandler) IsClicked() bool {
	return ih.isClicked
}
func (ih *DefaultHandler) GetMouseCursor() (int, int) {
	return ih.mouseX, ih.mouseY
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
