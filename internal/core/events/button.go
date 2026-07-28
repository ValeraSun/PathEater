package events

type WeaponEvent struct {
	Left     bool
	Right    bool
	Shooting bool
}

func (*WeaponEvent) Type() string { return "weapon" }
