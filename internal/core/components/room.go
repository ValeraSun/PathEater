package components

type RoomComponent struct {
	Vacuum bool
}

func (*RoomComponent) Type() string {
	return "room"
}

func NewRoomComponent() *RoomComponent {
	return &RoomComponent{
		Vacuum: false,
	}
}
