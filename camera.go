package noor

import (
	"math"

	"github.com/ahmedsat/madar"
	"github.com/ahmedsat/noor/input"
)

type ProjectionType int

const (
	Orthographic ProjectionType = iota
	Perspective
)

type CameraMode int

const (
	Free CameraMode = iota
	ThirdPerson
	FirstPerson
	Orbit
)

type Camera struct {
	Position  madar.Vector3
	Direction madar.Vector3
	Up        madar.Vector3
	Right     madar.Vector3

	Projection ProjectionType
	Mode       CameraMode
	Zoom       float32

	Width  float32
	Height float32

	FOV  float32
	Near float32
	Far  float32

	ProjectionMatrix madar.Matrix4X4
	ViewMatrix       madar.Matrix4X4

	// New fields
	Target           madar.Vector3
	OrbitDistance    float32
	OrbitSpeed       madar.Vector3
	MovementSpeed    float32
	RotationSpeed    float32
	MouseSensitivity float32
}

func NewCamera(width, height float32) *Camera {
	cam := &Camera{
		Position:         madar.Vector3{X: 0, Y: 0, Z: 3},
		Direction:        madar.Vector3{X: 0, Y: 0, Z: -1},
		Up:               madar.Vector3{X: 0, Y: 1, Z: 0},
		Projection:       Perspective,
		Mode:             Free,
		Zoom:             1,
		Width:            width,
		Height:           height,
		FOV:              45,
		Near:             0.1,
		Far:              100,
		Target:           madar.Vector3{X: 0, Y: 0, Z: 0},
		OrbitDistance:    5,
		OrbitSpeed:       madar.Vector3{X: 0.1, Y: 0.1, Z: 0.1},
		MovementSpeed:    0.1,
		RotationSpeed:    0.1,
		MouseSensitivity: 0.002,
	}
	cam.updateVectors()
	cam.updateProjectionMatrix()
	cam.updateViewMatrix()
	return cam
}

// ... [Keep the existing methods: updateVectors, updateProjectionMatrix, updateViewMatrix, Update, SetPosition, SetDirection, SetProjection, SetZoom, SetFOV] ...

func (c *Camera) Rotate(pitch, yaw, roll float32) {
	rotationMatrix := madar.RotationMatrix4X4(pitch, yaw, roll)
	c.Direction = rotationMatrix.MultiplyVector3(c.Direction).Normalize()
	c.Right = rotationMatrix.MultiplyVector3(c.Right).Normalize()
	c.Up = c.Right.Cross(c.Direction).Normalize()
	c.Update()
}

func (c *Camera) Move(direction madar.Vector3, distance float32) {
	c.Position = c.Position.Add(direction.Normalize().Scale(distance))
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

func (c *Camera) ZoomIn(amount float32) {
	c.Zoom += amount
	if c.Zoom < 0.1 {
		c.Zoom = 0.1
	}
	c.Update()
}

func (c *Camera) OrbitAround(target madar.Vector3, pitchDelta, yawDelta float32) {
	c.Target = target
	distance := c.Position.Sub(target).Length()

	// Calculate current spherical coordinates
	currentDir := c.Position.Sub(target).Normalize()
	pitch := float32(math.Asin(float64(currentDir.Y)))
	yaw := float32(math.Atan2(float64(currentDir.Z), float64(currentDir.X)))

	// Update angles
	pitch += pitchDelta * c.OrbitSpeed.Y
	yaw += yawDelta * c.OrbitSpeed.X

	// Clamp pitch to avoid flipping
	pitch = madar.Clamp(pitch, -math.Pi/2+0.1, math.Pi/2-0.1)

	// Calculate new position
	newPos := madar.Vector3{
		X: float32(math.Cos(float64(pitch)) * math.Cos(float64(yaw))),
		Y: float32(math.Sin(float64(pitch))),
		Z: float32(math.Cos(float64(pitch)) * math.Sin(float64(yaw))),
	}

	c.Position = target.Add(newPos.Scale(distance))
	c.SetDirection(-newPos.X, -newPos.Y, -newPos.Z)
}

func (c *Camera) SetTarget(target madar.Vector3) {
	c.Target = target
	c.Update()
}

func (c *Camera) SetMode(mode CameraMode) {
	c.Mode = mode
	c.Update()
}

func (c *Camera) HandleInput(deltaTime float32) {
	switch c.Mode {
	case Free:
		c.handleFreeCamera(deltaTime)
	case ThirdPerson:
		c.handleThirdPersonCamera(deltaTime)
	case FirstPerson:
		c.handleFirstPersonCamera(deltaTime)
	case Orbit:
		c.handleOrbitCamera(deltaTime)
	}
}

func (c *Camera) handleFreeCamera(deltaTime float32) {
	speed := c.MovementSpeed * deltaTime

	if input.IsKeyHeld(input.KeyW) {
		c.MoveForward(speed)
	}
	if input.IsKeyHeld(input.KeyS) {
		c.MoveForward(-speed)
	}
	if input.IsKeyHeld(input.KeyA) {
		c.MoveRight(-speed)
	}
	if input.IsKeyHeld(input.KeyD) {
		c.MoveRight(speed)
	}
	if input.IsKeyHeld(input.KeySpace) {
		c.MoveUp(speed)
	}
	if input.IsKeyHeld(input.KeyLeftShift) {
		c.MoveUp(-speed)
	}

	mouseDelta := input.GetMouseDelta()
	c.Rotate(-mouseDelta.Y*c.MouseSensitivity, -mouseDelta.X*c.MouseSensitivity, 0)
}

func (c *Camera) handleThirdPersonCamera(deltaTime float32) {
	// Implement third-person camera logic here
	// This could involve orbiting around a character, following a target, etc.
}

func (c *Camera) handleFirstPersonCamera(deltaTime float32) {
	// Similar to free camera, but without vertical movement
	speed := c.MovementSpeed * deltaTime

	if input.IsKeyHeld(input.KeyW) {
		c.MoveForward(speed)
	}
	if input.IsKeyHeld(input.KeyS) {
		c.MoveForward(-speed)
	}
	if input.IsKeyHeld(input.KeyA) {
		c.MoveRight(-speed)
	}
	if input.IsKeyHeld(input.KeyD) {
		c.MoveRight(speed)
	}

	mouseDelta := input.GetMouseDelta()
	c.Rotate(-mouseDelta.Y*c.MouseSensitivity, -mouseDelta.X*c.MouseSensitivity, 0)
}

func (c *Camera) handleOrbitCamera(deltaTime float32) {
	mouseDelta := input.GetMouseDelta()
	c.OrbitAround(c.Target, -mouseDelta.Y*c.MouseSensitivity, -mouseDelta.X*c.MouseSensitivity)

	scroll := input.GetMouseScroll()
	c.OrbitDistance -= scroll.Y * c.MovementSpeed
	if c.OrbitDistance < 1 {
		c.OrbitDistance = 1
	}
}

func (c *Camera) LookAt(target madar.Vector3) {
	c.Direction = target.Sub(c.Position).Normalize()
	c.Right = c.Up.Cross(c.Direction).Normalize()
	c.Up = c.Direction.Cross(c.Right)
	c.Update()
}

func (c *Camera) SetOrbitSpeed(x, y, z float32) {
	c.OrbitSpeed = madar.Vector3{X: x, Y: y, Z: z}
}

func (c *Camera) SetMovementSpeed(speed float32) {
	c.MovementSpeed = speed
}

func (c *Camera) SetRotationSpeed(speed float32) {
	c.RotationSpeed = speed
}

func (c *Camera) SetMouseSensitivity(sensitivity float32) {
	c.MouseSensitivity = sensitivity
}

// ///////////////////////////////////////

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
