package main

import (
	"log"
	"time"

	"github.com/kmdkuk/clicker/config"
	"github.com/kmdkuk/clicker/game"
	"github.com/kmdkuk/clicker/input"
	"github.com/kmdkuk/clicker/state"
	"github.com/kmdkuk/clicker/ui"

	"github.com/hajimehoshi/ebiten/v2"
	flag "github.com/spf13/pflag"
	"golang.org/x/net/context"
)

func main() {
	cfg := config.NewConfig()
	flag.BoolVarP(&cfg.EnableDebug, "debug", "d", false, "Enable debug mode")
	flag.Parse()
	gameState := state.NewGameState()
	storage := state.NewDefaultStorage(state.NewStorageDriver(config.DefaultSaveKey))
	if state, err := storage.LoadGameState(); err == nil {
		gameState = state
	}
	inputHandler := input.NewHandler()
	renderer := ui.NewRenderer(cfg, gameState, input.NewDecider(gameState))
	g := game.NewGame(
		cfg,
		gameState,
		storage,
		renderer,
		inputHandler,
	)
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Clicker")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g.StartAutoSave(ctx, 30*time.Second)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
