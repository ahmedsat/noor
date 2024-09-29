package noor

import (
	"math"

	"github.com/ahmedsat/madar"
)

type ProjectionType int

const (
	Orthographic ProjectionType = iota
	Perspective
)

type Camera struct {
	position       madar.Vector3
	lookAt         madar.Vector3
	up             madar.Vector3
	projectionType ProjectionType
	near           float32
	far            float32
	fov            float32 // Only for perspective
	aspect         float32
}

func NewCamera(position, lookAt, up madar.Vector3, aspect float32, projectionType ProjectionType) *Camera {
	return &Camera{
		position:       position,
		lookAt:         lookAt,
		up:             up,
		projectionType: projectionType,
		near:           0.1,
		far:            1000.0,
		fov:            45,
		aspect:         aspect,
	}
}

// Getter and setter methods remain the same, so I'll omit them for brevity

func (c *Camera) ViewMatrix() madar.Matrix4 {
	forward := madar.Vector3{
		c.lookAt[0] - c.position[0],
		c.lookAt[1] - c.position[1],
		c.lookAt[2] - c.position[2],
	}
	forward = forward.Normalize()

	right := forward.Cross(c.up).Normalize()
	up := right.Cross(forward)

	return madar.Matrix4{
		right[0], up[0], -forward[0], 0,
		right[1], up[1], -forward[1], 0,
		right[2], up[2], -forward[2], 0,
		-right.Dot(c.position), -up.Dot(c.position), forward.Dot(c.position), 1,
	}
}

func (c *Camera) ProjectionMatrix() madar.Matrix4 {
	switch c.projectionType {
	case Orthographic:
		return c.orthographicMatrix()
	case Perspective:
		return c.perspectiveMatrix()
	default:
		return c.orthographicMatrix()
	}
}

func (c *Camera) orthographicMatrix() madar.Matrix4 {
	right := 1.0 * c.aspect
	left := -right
	top := float32(1.0)
	bottom := -top

	return madar.OrthographicMatrix(left, right, bottom, top, c.near, c.far)
}

func (c *Camera) perspectiveMatrix() madar.Matrix4 {
	fov := c.fov * float32(math.Pi) / 180 // Convert to radians
	return madar.PerspectiveMatrix(fov, c.aspect, c.near, c.far)
}

func (c *Camera) SetAspect(aspect float32) {
	c.aspect = aspect
}
