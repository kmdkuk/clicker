package input

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Handler is an interface for handling input
type Handler interface {
	Update()
	GetPressedKey() KeyType
	IsClicked() bool
	IsMouseMoved() bool
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
	isMouseMoved bool
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

	mouseX, mouseY := ebiten.CursorPosition()
	ih.isMouseMoved = false
	if math.Abs(float64(mouseX-ih.mouseX)) > 10 || math.Abs(float64(mouseY-ih.mouseY)) > 10 {
		ih.isMouseMoved = true
		ih.mouseX = mouseX
		ih.mouseY = mouseY
	}
	if !ih.keepClicking && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		ih.keepClicking = true
		ih.isClicked = true
	}
	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		ih.keepClicking = false
	}
}

func (ih *DefaultHandler) ResetClickState() {
	ih.isClicked = false
}

func (ih *DefaultHandler) IsClicked() bool {
	return ih.isClicked
}
func (ih *DefaultHandler) IsMouseMoved() bool {
	return ih.isMouseMoved
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
