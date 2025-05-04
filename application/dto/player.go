package dto

type Player struct {
	Money             float64
	TotalGenerateRate float64
}

func (p *Player) GetMoney() float64 {
	return p.Money
}
func (p *Player) GetTotalGenerateRate() float64 {
	return p.TotalGenerateRate
}
