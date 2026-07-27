package components

type OxygenComponent struct {
	Oxygen    int
	MaxOxygen int
}

func NewOxygenComponent(maxOxygen int) *OxygenComponent {
	return &OxygenComponent{
		Oxygen:    maxOxygen,
		MaxOxygen: maxOxygen,
	}
}

func (*OxygenComponent) Type() string { return "oxygen" }

type isGasp bool

func (o *OxygenComponent) Leak(loss int) isGasp {
	o.Oxygen = o.Oxygen - loss
	if o.Oxygen <= 0 {
		return true
	}
	return false
}
