package game

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

func (p *Popup) HandleInput(keyType KeyType) {
	if keyType == KeyTypeDecision {
		p.Close()
	}
}
