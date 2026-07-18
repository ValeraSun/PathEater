package systems

type ShipSystem struct {
	getter      componentsGetter
	broadcaster Broadcaster
}

func NewShipSystem(getter componentsGetter, broadcaster Broadcaster, subscriber subscriber) *ShipSystem {
	return &ShipSystem{
		getter:      getter,
		broadcaster: broadcaster,
	}
	subscriber.Subscribe("setShipState", s.OnEvent)
	return s
}

func (s *ShipSystem) Update(dt float32) error {
	comps := s.getter.GetEntitiesByComponent("ship")

	s.rotateWeapon(comps)

	return nil
}
