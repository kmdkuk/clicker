package components

import (
	"github.com/kmdkuk/clicker/presentation/input"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Popup struct {
	Message string // The message displayed in the popup
	Active  bool   // Indicates whether the popup is currently active
}

func NewPopup() *Popup {
	return &Popup{
		Message: "",
		Active:  false,
	}
}

func (p *Popup) Draw(screen *ebiten.Image) {
	if p.IsActive() {
		ebitenutil.DebugPrintAt(screen, "Popup: "+p.Message, 10, 200)
	}
}

func (p *Popup) HandleInput(keyType input.KeyType) {
	if p.IsActive() && keyType == input.KeyTypeDecision {
		p.Close()
	}
}

func (p *Popup) IsActive() bool {
	return p.Active
}

func (p *Popup) GetMessage() string {
	return p.Message
}

func (p *Popup) Show(message string) {
	p.Message = message
	p.Active = true
}

func (p *Popup) Close() {
	p.Active = false
}
