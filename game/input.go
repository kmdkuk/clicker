package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type InputHandler struct{}

func (ih *InputHandler) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		// Handle escape key press
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		// Handle enter key press
	}
	// Add more input handling as needed
}

func (ih *InputHandler) MousePosition() (int, int) {
	x, y := ebiten.CursorPosition()
	return x, y
}

func (ih *InputHandler) IsMouseButtonPressed(button ebiten.MouseButton) bool {
	return ebiten.IsMouseButtonPressed(button)
}