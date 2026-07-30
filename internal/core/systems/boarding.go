package systems

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ValeraSun/PathEater/internal/config"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
)

type BoardingSystem struct {
	componentsGetter
	subscriber
	publisher
	spawns     []geometry.Vec3
	eventQueue chan *events.BoardingEvent
	rng        *rand.Rand
}

func NewBoardingSystem(getter componentsGetter, subscriber subscriber, publisher publisher) *BoardingSystem {
	s := &BoardingSystem{
		getter,
		subscriber,
		publisher,
		config.GetAlienSpawns(),
		make(chan *events.BoardingEvent, 100),
		rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	s.Subscribe("boarding", s.OnEvent)
	return s
}

func (s *BoardingSystem) GetAlienSpawn() geometry.Vec3 {
	i := s.rng.Intn(len(s.spawns))

	return s.spawns[i]
}

func (s *BoardingSystem) Update(dt float32) error {
	for {
		select {
		case <-s.eventQueue:

			// create := events.NewCreateAlienEvent(s.GetAlienSpawn())
			// s.Publish(create)

		default:
			return nil
		}
	}
}

func (s *BoardingSystem) OnEvent(event events.Event) error {
	e := event.(*events.BoardingEvent)

	select {
	case s.eventQueue <- e:
	default:
		fmt.Println("Преполена очередь boarding")
	}

	return nil
}
