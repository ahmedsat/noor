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
	Position       madar.Vector3
	LookAt         madar.Vector3
	Up             madar.Vector3
	ProjectionType ProjectionType
	Near           float32
	Far            float32
	Fov            float32
	Aspect         float32
}

// Getter and setter methods remain the same, so I'll omit them for brevity

func (c *Camera) ViewMatrix() madar.Matrix4 {
	forward := madar.Vector3{
		c.LookAt[0] - c.Position[0],
		c.LookAt[1] - c.Position[1],
		c.LookAt[2] - c.Position[2],
	}
	forward = forward.Normalize()

	right := forward.Cross(c.Up).Normalize()
	up := right.Cross(forward)

	return madar.Matrix4{
		right[0], up[0], -forward[0], 0,
		right[1], up[1], -forward[1], 0,
		right[2], up[2], -forward[2], 0,
		-right.Dot(c.Position), -up.Dot(c.Position), forward.Dot(c.Position), 1,
	}
}

func (c *Camera) ProjectionMatrix() madar.Matrix4 {
	switch c.ProjectionType {
	case Orthographic:
		return c.orthographicMatrix()
	case Perspective:
		return c.perspectiveMatrix()
	default:
		return c.orthographicMatrix()
	}
}

func (c *Camera) orthographicMatrix() madar.Matrix4 {
	right := 1.0 * c.Aspect
	left := -right
	top := float32(1.0)
	bottom := -top

	return madar.OrthographicMatrix(left, right, bottom, top, c.Near, c.Far)
}

func (c *Camera) perspectiveMatrix() madar.Matrix4 {
	fov := c.Fov * float32(math.Pi) / 180 // Convert to radians
	return madar.PerspectiveMatrix(fov, c.Aspect, c.Near, c.Far)
}

func (c *Camera) SetAspect(aspect float32) {
	c.Aspect = aspect
}
