package player

type Player struct {
	ID        string
	X, Y, Z   float64
	Health    int
	IsAlive   bool
}

func NewPlayer(id string) *Player {
	return &Player{
		ID:     id,
		X:      0, Y: 0, Z: 0,
		Health: 100,
		IsAlive: true,
	}
}

func (p *Player) Move(dx, dy, float64) {
	p.X += dx
	p.Y += dy
}