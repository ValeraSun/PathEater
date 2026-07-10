package ecs

type Entity struct {
	ID   string
	Mask uint64
}

const (
    ComponentTransform uint64 = 1 << iota
	ComponentMove
	ComponentPlayer
	ComponentHealth
	ComponentAttack
	ComponentUseItem
)

type Component interface {
	ComponentType()
}

type BaseComponent struct {
	Type    uint64
	OwnerID string
}
func (b BaseComponent) ComponentType() uint64 {
    return b.Type
}

//функции работы с маской!!!----------------------------------------------------
func SetBit(mask uint64, bit uint64) uint64 {
	return mask | bit
}
func DeleteBit(mask uint64, bit uint64) uint64 {
	return mask & ^bit
}
func HasBit(mask uint64, bit uint64) bool {
	return (mask & bit) != 0
}
//функции работы с маской!!!----------------------------------------------------

func NewEntity(id string) *Entity {
    return &Entity{
        ID:   id,
        Mask: 0,
    }
}

func (e *Entity) AddComponents(components ...BaseComponent) {
	for _, comp := range components {
		e.Mask = SetBit(e.Mask, comp.Type)
	}
}

func (e *Entity) RemoveComponents(components ...BaseComponent) {
	for _, comp := range components {
		e.Mask = DeleteBit(e.Mask, comp.Type)
	}
}

func (e *Entity) HasComponents(components ...BaseComponent) bool {
	for _, comp := range components {
		if !HasBit(e.Mask, comp.Type) {
			return false
		}
	}
	return true
}