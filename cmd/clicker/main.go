package main

import (
	"context"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	flag "github.com/spf13/pflag"

	"github.com/kmdkuk/clicker/application/usecase"
	"github.com/kmdkuk/clicker/config"
	"github.com/kmdkuk/clicker/game"
	"github.com/kmdkuk/clicker/infrastructure/state"
	"github.com/kmdkuk/clicker/infrastructure/storage"
	"github.com/kmdkuk/clicker/infrastructure/storage/driver"
	"github.com/kmdkuk/clicker/presentation"
	"github.com/kmdkuk/clicker/presentation/input"
)

func main() {
	cfg := config.NewConfig()
	flag.BoolVarP(&cfg.EnableDebug, "debug", "d", false, "Enable debug mode")
	flag.Parse()
	gameState := state.NewGameState()
	storage := storage.NewDefaultStorage(driver.NewStorageDriver(config.DefaultSaveKey))
	if state, err := storage.LoadGameState(); err == nil {
		gameState = state
	}
	renderer, err := presentation.NewRenderer(
		cfg,
		usecase.NewPlayerUsecase(gameState),
		usecase.NewManualWorkUseCase(gameState),
		usecase.NewBuildingUseCase(gameState),
		usecase.NewUpgradeUseCase(gameState),
	)
	if err != nil {
		log.Fatal(err)
	}
	inputHandler := input.NewHandler()
	g := game.NewGame(
		cfg,
		gameState,
		storage,
		renderer,
		inputHandler,
	)
	ebiten.SetWindowSize(cfg.ScreenWidth, cfg.ScreenHeight)
	ebiten.SetWindowTitle("Clicker")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g.StartAutoSave(ctx, 30*time.Second)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
