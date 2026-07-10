package events

type EventState struct {
	MoveFront bool
	MoveRight bool 
	MoveBack  bool 
	MoveLeft  bool
	Interact  bool 
	Attack    bool 
	//Direction Vector3
}

func (e *EventState) Type() string{ return "playerState" }

func NewPlayerStateEvent(id string, moveState MoveState) EventState {
	return EventState{
		ID:        id,
		MoveFront: moveState.MoveFront,
		MoveRight: moveState.MoveRight,
		MoveBack:  moveState.MoveBack,
		MoveLeft:  moveState.MoveLeft,
		Interact:  moveState.Interact,
		Attack:    moveState.Attack,
	}
}