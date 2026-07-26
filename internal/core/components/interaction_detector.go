package components

import "github.com/ValeraSun/PathEater/internal/core/geometry"

type InteractionDetectorComponent struct {
	distant      float64
	cameraHeight float64
	cameraOffset float64
}

func NewInteractionDetectorComponent(distant, cameraHeight, cameraOffset float64) *InteractionDetectorComponent {
	return &InteractionDetectorComponent{
		distant:      distant,
		cameraHeight: cameraHeight,
		cameraOffset: cameraOffset,
	}
}

func (*InteractionDetectorComponent) Type() string { return "interactionDetector" }

func (c *InteractionDetectorComponent) GetRay(pos, dir geometry.Vec3) *geometry.RayCollider {
	dirY := dir
	dirY.Y = 0
	dirY = dirY.Normalize()
	pos = pos.Add(geometry.Vec3{Y: c.cameraHeight})
	pos = pos.Add(dirY.Scale(c.cameraOffset))
	return &geometry.RayCollider{
		Origin:    pos,
		Direction: dir,
		Length:    c.distant,
	}
}
