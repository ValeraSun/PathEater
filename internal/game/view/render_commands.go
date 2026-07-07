package view

type RenderCommand struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type CmdPlayerState struct {
	ID     string  `json:"id"`
	X, Y, Z float64 `json:"pos"`
	Health int     `json:"health"`
}

func NewCmdPlayerState(p *player.Player) *RenderCommand {
	return &RenderCommand{
		Type: "player_state",
		Data: CmdPlayerState{
			ID: p.ID, X: p.X, Y: p.Y, Z: p.Z, Health: p.Health,
		},
	}
}