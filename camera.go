package noor

import (
	"github.com/ahmedsat/madar"
	"github.com/ahmedsat/noor/input"
)

type ProjectionType int
type CameraMode int

const (
	Orthographic ProjectionType = iota
	Perspective
)

const (
	Fixed CameraMode = iota
	Free
	ThirdPerson
	FirstPerson
	Orbit
)

type Camera struct {
	Position         madar.Vector3
	Direction        madar.Vector3
	Up               madar.Vector3
	Right            madar.Vector3
	Projection       ProjectionType
	Mode             CameraMode
	Zoom             float32
	Width            float32
	Height           float32
	FOV              float32
	Near             float32
	Far              float32
	ProjectionMatrix madar.Matrix4X4
	ViewMatrix       madar.Matrix4X4
	MovementSpeed    float32
	MouseSensitivity float32
	Damping          float32
	Target           madar.Vector3
	OrbitDistance    float32
	OrbitSpeed       madar.Vector3
	isDirty          bool
}

// Creates a new Camera object
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
		MovementSpeed:    5.0,
		MouseSensitivity: 0.002,
		Damping:          0.9,
		Target:           madar.Vector3{X: 0, Y: 0, Z: 0},
		OrbitDistance:    5,
		OrbitSpeed:       madar.Vector3{X: 0.5, Y: 0.5, Z: 0},
		isDirty:          true,
	}
	cam.Update()
	return cam
}

// Update the camera's matrices and apply changes
func (c *Camera) Update() {
	if !c.isDirty {
		return
	}
	c.updateVectors()
	c.updateProjectionMatrix()
	c.updateViewMatrix()
	c.isDirty = false
}

// Applies damping to smooth out movements over time
func (c *Camera) ApplyDamping(value float32) float32 {
	return value * c.Damping
}

// Handles movement by direction vector
func (c *Camera) Move(direction madar.Vector3, distance float32) {
	distance = c.ApplyDamping(distance)
	c.Position = c.Position.Add(direction.Normalize().Scale(distance))
	c.isDirty = true
}

// Handles free movement (FPS-like camera)
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

// Orbit Camera handling
func (c *Camera) handleOrbitCamera(deltaTime float32) {
	// Handle orbit rotation around the target
	orbitX := c.OrbitSpeed.X * deltaTime
	orbitY := c.OrbitSpeed.Y * deltaTime

	rotationMatrix := madar.RotationMatrix4X4(orbitY, orbitX, 0)
	offset := c.Position.Sub(c.Target)

	// Calculate the new camera position
	newPosition := rotationMatrix.MultiplyVector3(offset).Add(c.Target)
	c.Position = newPosition

	// Update the look direction
	c.LookAt(c.Target)

	if input.IsKeyHeld(input.KeyQ) {
		c.ZoomIn(0.1)
	}
	if input.IsKeyHeld(input.KeyE) {
		c.ZoomIn(-0.1)
	}
	c.isDirty = true
}

// Handles Third-Person camera (chase cam)
func (c *Camera) handleThirdPersonCamera(deltaTime float32) {
	// Ensure the camera is behind the target
	offset := c.Direction.Scale(c.OrbitDistance)
	c.Position = c.Target.Sub(offset)
	c.LookAt(c.Target)
	c.isDirty = true
}

// Updates the projection matrix based on perspective/orthographic view
func (c *Camera) updateProjectionMatrix() {
	aspect := c.Width / c.Height
	if c.Projection == Perspective {
		c.ProjectionMatrix = madar.PerspectiveMatrix4X4(c.FOV, aspect, c.Near, c.Far)
	} else {
		size := c.Zoom * 10
		c.ProjectionMatrix = madar.OrthographicMatrix4X4(-size*aspect, size*aspect, -size, size, c.Near, c.Far)
	}
}

// Updates the view matrix based on camera position and target
func (c *Camera) updateViewMatrix() {
	center := c.Position.Add(c.Direction)
	c.ViewMatrix = madar.LookAtMatrix4X4(c.Position, center, c.Up)
}

// Rotates the camera by applying a pitch, yaw, roll
func (c *Camera) Rotate(pitch, yaw, roll float32) {
	rotationMatrix := madar.RotationMatrix4X4(pitch, yaw, roll)
	c.Direction = rotationMatrix.MultiplyVector3(c.Direction).Normalize()
	c.Right = rotationMatrix.MultiplyVector3(c.Right).Normalize()
	c.Up = c.Right.Cross(c.Direction).Normalize()
	c.isDirty = true
}

// LookAt forces the camera to look at a specific target
func (c *Camera) LookAt(target madar.Vector3) {
	c.Direction = target.Sub(c.Position).Normalize()
	c.updateVectors()
}

// Updates the camera's directional vectors (Right, Up)
func (c *Camera) updateVectors() {
	c.Right = c.Direction.Cross(c.Up).Normalize()
	c.Up = c.Right.Cross(c.Direction).Normalize()
}

// ////////////////////////////////////

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

func (c *Camera) SetPosition(x, y, z float32) {
	c.Position = madar.Vector3{X: x, Y: y, Z: z}
	c.isDirty = true
}

func (c *Camera) SetDirection(x, y, z float32) {
	c.Direction = madar.Vector3{X: x, Y: y, Z: z}.Normalize()
	c.isDirty = true
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

func (c *Camera) SetMode(mode CameraMode) {
	c.Mode = mode
	switch c.Mode {
	case Free:
		c.MovementSpeed = 0.1
	case Orbit:
		c.OrbitDistance = 5.0
		c.OrbitSpeed = madar.Vector3{X: 0.5, Y: 0.5, Z: 0}
	}
	c.Update()
}

func (c *Camera) handleFirstPersonCamera(deltaTime float32) {
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

func (c *Camera) SetTarget(x, y, z float32) {
	c.Target = madar.Vector3{X: x, Y: y, Z: z}
	c.Update()
}

func (c *Camera) SetOrbitDistance(distance float32) {
	c.OrbitDistance = distance
	c.Update()
}
