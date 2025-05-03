package game

import (
	"context"
	"log"
	"time"

	"github.com/kmdkuk/clicker/config"
	"github.com/kmdkuk/clicker/infrastructure/state"
	"github.com/kmdkuk/clicker/infrastructure/storage"
	"github.com/kmdkuk/clicker/presentation"
	"github.com/kmdkuk/clicker/presentation/input"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	config       *config.Config  // Game configuration
	gameState    state.GameState // Game state
	storage      storage.Storage
	inputHandler input.Handler         // Handler to manage input processing
	renderer     presentation.Renderer // Update Renderer to use the presentation package
}

func NewGame(c *config.Config, gameState state.GameState, storage storage.Storage, renderer presentation.Renderer, inputHandler input.Handler) *Game {
	return &Game{
		config:       c,
		gameState:    gameState,
		storage:      storage,
		inputHandler: inputHandler,
		renderer:     renderer,
	}
}

func (g *Game) StartAutoSave(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Printf("Auto-save stopped")
				return
			case <-ticker.C:
				// Auto-save the game state
				err := g.storage.SaveGameState(g.gameState)
				if err != nil {
					log.Printf("Auto-save failed: %v", err)
				}
			}
		}
	}()
}

func (g *Game) Update() error {
	g.inputHandler.Update() // Update input handler

	g.gameState.UpdateBuildings(time.Now())

	// Update game state
	g.renderer.HandleInput(g.inputHandler.GetPressedKey())

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.renderer.Draw(screen) // Delegate drawing to renderer
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480 // Set the game screen size
}

// GetTotalGenerateRate calculates the total money generation rate from all unlocked buildings
func (g *Game) GetTotalGenerateRate() float64 {
	totalRate := 0.0
	for _, building := range g.gameState.GetBuildings() {
		if building.IsUnlocked() {
			totalRate += building.TotalGenerateRate(g.gameState.GetUpgrades())
		}
	}
	return totalRate
}
