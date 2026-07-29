package components

type RoomComponent struct {
	HasBreakdown bool
	Oxygen       float64
	Vacuum       bool
	MinX         float64
	MaxX         float64
	MinY         float64
	MaxY         float64
}

func (*RoomComponent) Type() string {
	return "room"
}

func NewRoomComponent(minX, maxX, minY, maxY float64) *RoomComponent {
	return &RoomComponent{
		HasBreakdown: false,
		Oxygen:       100.0,
		Vacuum:       false,
		MinX:         minX,
		MaxX:         maxX,
		MinY:         minY,
		MaxY:         maxY,
	}
}
