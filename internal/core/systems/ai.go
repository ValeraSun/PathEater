package systems

type AISystem struct {
	getter componentsGetter
}

func NewAISystem(getter componentsGetter) *AISystem {
	return &AISystem{
		getter: getter,
	}
}

func (s *AISystem) Update(dt float32) error {
	// comps := s.getter.GetEntitiesByComponent("ai")

	// enemies := make([])
	// for id, c := range comps {
	// 	if s.getter.HasComponents(id, "transform") {

	// 	}
	// }

	return nil
}
