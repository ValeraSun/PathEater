package systems

type HealthSystem struct {
	getter     componentsGetter
	eventQueue  chan events.Event
}

func NewHealthSystem(getter componentsGetter, subscriber subscriber) *HealthSystem {
	s := &HealthSystem{
		getter: getter,
		eventQueue: make(chan events.Event, 100),
	}
	subscriber.Subscribe("damageShip", s.OnEvent)
	subscriber.Subscribe("damageCosmoAlient", s.OnEvent)
	return s
}

func (*HealthSystem) Update(dt float32) error {
	comps := s.getter.GetEntitiesByComponent("health")
	for {
		var err error
		select {
		case e := <-s.eventQueue:
			typ := e.Type()

			switch {
			case typ == "damageShip":
				ev, _ := e.(*events.DamageShipEvent)
				c, _ := s.getter.GetComponent(ev.ID, "health")
				hp := c.(*HealthComponent)
				isDead := hp.Damage(ev.Damage)
				if isDead {
					//обработка проигрыша
				}
			case typ == "damageCosmoAlien":
				ev, _ := e.(*events.DamageCosmoAlienEvent)
				c, _ := s.getter.GetComponent(ev.ID, "health")
				hp := c.(*HealthComponent)
				isDead := hp.Damage(ev.Damage)
				if isDead {
					//обработка смерти пришельца
				}
			}
		default:
			return nil
		}

		if err != nil {
			return err
		}
	}
}

func (s *HealthSystem) OnEvent(event events.Event) error {
	select {
	case s.eventQueue <- event:
		e := <- s.eventQueue
		ev := e.()
	default:
		fmt.Printf("Преполена очередь %v\n", s)
	}

	return nil
}
