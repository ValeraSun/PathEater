package systems

type TransformSystem struct {
	getter componentsGetter
}

func NewTransformSystem(getter componentsGetter) *TransformSystem {
	return &TransformSystem{
		getter: getter,
	}
}
func (*TransformSystem) Update(dt float32) error {
	return nil
}
