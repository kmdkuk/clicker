package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kmdkuk/clicker/game"
)

func main() {
	g := game.NewGame()
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Clicker")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
