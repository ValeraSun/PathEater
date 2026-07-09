package events

type EventMove struct {
    X float32
    Y float32
    Z float32
}

func (e *EventMove) Type() string{ return "Move" }

func CreateEventMove(X float32, Y float32, Z float32) EventMove {
	return EventMove{
		X: X,
		Y: Y,
		Z: Z,
	}
}

type EventUseItem struct{
	ItemID string
}

func (e *EventUseItem) Type() string{ return "UseItem" }

func CreateEventUseItem(itemID string) EventUseItem {
	return EventUseItem{
		ItemID: itemID,
	}
}

type EventExit struct{}

func (e *EventExit) Type() string{ return "Exit" }

func CreateEventExit() EventExit {
	return EventExit{}
}