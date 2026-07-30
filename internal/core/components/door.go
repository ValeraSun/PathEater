package components

type DoorComponent struct {
	RoomA       string
	RoomB       string
	IsOpen      bool
	TriggerMinX float64
	TriggerMaxX float64
	TriggerMinY float64
	TriggerMaxY float64
}

func (*DoorComponent) Type() string {
	return "door"
}

func NewDoorComponent(isOpen bool, roomA, roomB string, tminx, tmaxx, tminy, tmaxy float64) *DoorComponent {
	return &DoorComponent{
		RoomA:       roomA,
		RoomB:       roomB,
		IsOpen:      isOpen,
		TriggerMinX: tminx,
		TriggerMaxX: tmaxx,
		TriggerMinY: tminy,
		TriggerMaxY: tmaxy,
	}
}
