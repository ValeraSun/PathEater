package systems

type PlayerConrtollerSystem struct {
	Update(world *World, dt float64) error
	World *ecs.World
}

func (s PlayerControllerSystem) Update(world *World, dt float64){}

func (s PlayerControllerSystem) OnEvent(event EventMove) EventHandler {
	comp, exists := s.World.[event.ID]["transform"]
	component := comp.()
	comp.velocity
}

func NewPlayerControllerSystem(w *ecs.World) PlayerConrtollerSystem {
	return PlayerConrtollerSystem{
		World: w,
	}
}