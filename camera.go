package noor

import "github.com/ahmedsat/madar"

type Projection interface {
	GetProjectionMatrix(c Camera) madar.Matrix4X4
	Update()
}

type Perspective struct {
	Fov    float32
	Aspect float32
	Near   float32
	Far    float32

	matrix madar.Matrix4X4
}

func (p *Perspective) GetProjectionMatrix(c Camera) madar.Matrix4X4 {
	return p.matrix
}

func (p *Perspective) Update() {
	p.matrix = madar.PerspectiveMatrix4X4(p.Fov, p.Aspect, p.Near, p.Far)
}

type Camera struct {
	// editable by user
	Position, Direction, Up madar.Vector3
	Projection              Projection

	// managed by noor6
	target     madar.Vector3
	viewMatrix madar.Matrix4X4
}

func NewCamera(position, direction, up madar.Vector3, projection Projection) *Camera {
	c := &Camera{
		Position:   position,
		Direction:  direction,
		Up:         up,
		Projection: projection,
	}

	c.Update(0)

	return c
}

func (c *Camera) GetViewMatrix() madar.Matrix4X4 { return c.viewMatrix }

func (c *Camera) Update(dt float32) {

	c.target = c.Position.Add(c.Direction)

	c.viewMatrix = madar.LookAtMatrix4X4(c.Position, c.target, c.Up)

	c.Projection.Update()
}

func (c *Camera) GetProjectionMatrix() madar.Matrix4X4 { return c.Projection.GetProjectionMatrix(*c) }

func (c *Camera) Cleanup() {}

func (c *Camera) MoveForward(distance float32) {
	c.Position = c.Position.Add(c.Direction.Scale(distance))
	c.Update(0)
}
func (c *Camera) MoveBackward(distance float32) {
	c.Position = c.Position.Sub(c.Direction.Scale(distance))
	c.Update(0)
}
func (c *Camera) MoveLeft(distance float32) {
	c.Position = c.Position.Add(c.Up.Cross(c.Direction).Normalize().Scale(distance))
	c.Update(0)
}
func (c *Camera) MoveRight(distance float32) {
	c.Position = c.Position.Sub(c.Up.Cross(c.Direction).Normalize().Scale(distance))
	c.Update(0)
}
func (c *Camera) MoveUp(distance float32) {
	c.Position = c.Position.Add(c.Up.Normalize().Scale(distance))
	c.Update(0)
}
func (c *Camera) MoveDown(distance float32) {
	c.Position = c.Position.Sub(c.Up.Normalize().Scale(distance))
	c.Update(0)
}

func (c *Camera) Rotate(rotation madar.Vector3) {
	rotationMatrix := madar.RotationMatrix4X4(rotation)
	c.Direction = rotationMatrix.MultiplyVector3(c.Direction)
	c.Up = rotationMatrix.MultiplyVector3(c.Up)
	c.Update(0)
}

func (c *Camera) LookAt(target madar.Vector3) {
	c.Direction = target.Sub(c.Position).Normalize()

	globalUp := madar.Vector3{X: 0, Y: 1, Z: 0}
	right := c.Direction.Cross(globalUp).Normalize()
	c.Up = right.Cross(c.Direction).Normalize()

	c.Update(0)
}
