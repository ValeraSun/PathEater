package systems

type HealthSystem struct {
	getter componentsGetter
}

func NewHealthSystem(getter componentsGetter) *HealthSystem {
	return &HealthSystem{
		getter: getter,
	}
}

func (*HealthSystem) Update(dt float32) error {
	return nil
}
