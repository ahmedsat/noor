package noor

import (
	"github.com/ahmedsat/madar"
)

type ProjectionType int

const (
	Orthographic ProjectionType = iota
	Perspective
)

type Camera struct {
	Position  madar.Vector3
	Direction madar.Vector3
	Up        madar.Vector3
	Right     madar.Vector3

	Projection ProjectionType
	Zoom       float32

	Width  float32
	Height float32

	FOV  float32
	Near float32
	Far  float32

	ProjectionMatrix madar.Matrix4X4
	ViewMatrix       madar.Matrix4X4
}

func NewCamera(width, height float32) *Camera {
	cam := &Camera{
		Position:   madar.Vector3{X: 0, Y: 0, Z: 3},
		Direction:  madar.Vector3{X: 0, Y: 0, Z: -1},
		Up:         madar.Vector3{X: 0, Y: 1, Z: 0},
		Projection: Perspective,
		Zoom:       1,
		Width:      width,
		Height:     height,
		FOV:        45,
		Near:       0.1,
		Far:        100,
	}
	cam.updateVectors()
	cam.updateProjectionMatrix()
	cam.updateViewMatrix()
	return cam
}

func (c *Camera) updateVectors() {
	c.Right = c.Direction.Cross(c.Up).Normalize()
	c.Up = c.Right.Cross(c.Direction).Normalize()
}

func (c *Camera) updateProjectionMatrix() {
	aspect := c.Width / c.Height
	if c.Projection == Perspective {
		c.ProjectionMatrix = madar.PerspectiveMatrix4X4(c.FOV, aspect, c.Near, c.Far)
	} else {
		size := c.Zoom * 10
		c.ProjectionMatrix = madar.OrthographicMatrix4X4(-size*aspect, size*aspect, -size, size, c.Near, c.Far)
	}
}

func (c *Camera) updateViewMatrix() {
	center := c.Position.Add(c.Direction)
	c.ViewMatrix = madar.LookAtMatrix4X4(c.Position, center, c.Up)
}

func (c *Camera) Update() {
	c.updateVectors()
	c.updateProjectionMatrix()
	c.updateViewMatrix()
}

func (c *Camera) SetPosition(x, y, z float32) {
	c.Position = madar.Vector3{X: x, Y: y, Z: z}
	c.Update()
}

func (c *Camera) SetDirection(x, y, z float32) {
	c.Direction = madar.Vector3{X: x, Y: y, Z: z}.Normalize()
	c.Update()
}

func (c *Camera) SetProjection(projType ProjectionType) {
	c.Projection = projType
	c.Update()
}

func (c *Camera) SetZoom(zoom float32) {
	c.Zoom = zoom
	c.Update()
}

func (c *Camera) SetFOV(fov float32) {
	c.FOV = fov
	c.Update()
}

func (c *Camera) Rotate(yaw, pitch, roll float32) {
	c.Direction = c.Direction.Rotate(yaw, pitch, roll).Normalize()
	c.Right = c.Direction.Cross(c.Up).Normalize()
	c.Up = c.Right.Cross(c.Direction).Normalize()
	c.Update()
}

func (c *Camera) Move(direction madar.Vector3, distance float32) {
	c.Position = c.Position.Add(direction.Scale(distance))
	c.Update()
}

func (c *Camera) MoveForward(distance float32) {
	c.Move(c.Direction, distance)
}

func (c *Camera) MoveRight(distance float32) {
	c.Move(c.Right, distance)
}

func (c *Camera) MoveUp(distance float32) {
	c.Move(c.Up, distance)
}
