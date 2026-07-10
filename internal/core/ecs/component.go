package ecs

type Transform struct {
	BaseComponent
	X, Y, Z  float32
	Rotation float32
	Scale    float32
}

func NewTransform(x, y, z float32) *Transform {
    return &Transform{
        BaseComponent: BaseComponent{Type: ComponentTransform},
        X:             x,
        Y:             y,
        Z:             z,
        Rotation:      0,
        Scale:         1.0,
    }
} 

type Move struct {
	BaseComponent
	X, Y, Z  float32
}

func NewMove(x, y, z float32) *Move{
    return &Move{
        BaseComponent: BaseComponent{Type: ComponentTransform},
        X:             x,
        Y:             y,
        Z:             z,
    }
} 

type Health struct {
    BaseComponent
    HP           int
    MaxHP        int
    Regeneration int 
    IsDead       bool
}

func NewHealth(maxHP int) *Health {
    return &Health{
        BaseComponent: BaseComponent{Type: ComponentHealth},
        HP:            maxHP,
        MaxHP:         maxHP,
        Regeneration:  0,
        IsDead:        false,
    }
}

type Player struct {
    BaseComponent
    ID         string 
    IsAdmin    bool
    Input      PlayerInput
}

type PlayerInput struct {
    MoveX float32
    MoveY float32 
    MoveZ float32 
}

func NewPlayer(id string) *Player {
    return &Player{
        BaseComponent: BaseComponent{Type: ComponentPlayer},
        ID:            id,
        IsAdmin:       false,
    }
}

type Attack struct {
    BaseComponent
    Damage       int
	ItemDamageID string
}

func NewAttack(id string, damage int) *Attack {
    return &Attack{
        BaseComponent: BaseComponent{Type: ComponentAttack},
		Damage:        damage
        ItemDamageID:  id,
    }
}

type UseItem struct {
    BaseComponent
    TypeUse string
	ItemID  string
}

func NewUseItem(id string, TypeUse string) *UseItem {
    return &UserItem{
        BaseComponent: BaseComponent{Type: ComponentUseItem},
		TypeUse:       TypeUse
        ItemID:        id,
    }
}