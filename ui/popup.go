package ui

import "github.com/kmdkuk/clicker/input"

type Popup struct {
	Message string // Message to display
	Active  bool   // Whether the popup is active
}

func (p *Popup) Show(message string) {
	p.Message = message
	p.Active = true
}

func (p *Popup) Close() {
	p.Active = false
}

func (p *Popup) HandleInput(keyType input.KeyType) {
	if keyType == input.KeyTypeDecision {
		p.Close()
	}
}
