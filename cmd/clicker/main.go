package main

import (
	"log"

	"github.com/kmdkuk/clicker/config"
	"github.com/kmdkuk/clicker/game"

	"github.com/hajimehoshi/ebiten/v2"
	flag "github.com/spf13/pflag"
)

func main() {
	cfg := config.NewConfig()
	flag.BoolVarP(&cfg.EnableDebug, "debug", "d", false, "Enable debug mode")
	flag.Parse()
	g := game.NewGame(cfg)
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Clicker")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
