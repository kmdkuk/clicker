package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type Game struct {
	// Add game state variables here
}

func NewGame() *Game {
	return &Game{
		// Initialize game state variables here
	}
}

func (g *Game) Update() error {
	// Update game state logic here
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255}) // Fill the screen with black
	// Draw game elements here
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480 // Set the game's screen size
}