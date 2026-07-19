package components

type BulletComponent struct {
	Visible bool
}

func (*BulletComponent) Type() string {
	return "bullet"
}

func NewBulletComponent() *BulletComponent {
	return &BulletComponent{
		Visible: true,
	}
}
