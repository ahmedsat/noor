package noor

import (
	"math"
	"time"

	"github.com/ahmedsat/bayaan"
	"github.com/ahmedsat/madar"
	"github.com/ahmedsat/noor/input"
)

// Add camera control keys
const (
	moveForwardKey  = input.KeyW
	moveBackwardKey = input.KeyS
	moveLeftKey     = input.KeyA
	moveRightKey    = input.KeyD
	moveUpKey       = input.KeyE
	moveDownKey     = input.KeyQ

	// Camera mode control
	toggleMouseLockKey = input.KeyLeftAlt
	resetCameraKey     = input.KeyR
)

// Add camera control options
type CameraControls struct {
	// Movement speeds
	BaseMovementSpeed float32
	SprintMultiplier  float32

	// Mouse sensitivity
	MouseSensitivity float32

	// Zoom settings
	ZoomSpeed float32
	MinZoom   float32
	MaxZoom   float32

	// FOV settings
	FOVSpeed float32
	MinFOV   float32
	MaxFOV   float32

	// Orbit settings
	OrbitSpeed     float32
	MinOrbitRadius float32
	MaxOrbitRadius float32
}

type ProjectionType int

func (p ProjectionType) String() string {
	switch p {
	case Orthographic:
		return "Orthographic"
	case Perspective:
		return "Perspective"
	default:
		return "Unknown"
	}
}

type CameraMode int

func (c CameraMode) String() string {
	switch c {
	case Fixed:
		return "Fixed"
	case Free:
		return "Free"
	case ThirdPerson:
		return "ThirdPerson"
	case FirstPerson:
		return "FirstPerson"
	case Orbit:
		return "Orbit"
	default:
		return "Unknown"
	}
}

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

const (
	DefaultBaseMovementSpeed = 1
	DefaultSprintMultiplier  = 2
	DefaultFOV               = 45.0
	DefaultMinFOV            = 5
	DefaultMaxFOV            = 110
	DefaultZoom              = 1.0
	DefaultNear              = 0.1
	DefaultFar               = 100.0
	DefaultYaw               = 0.0
	DefaultPitch             = 0.0
	DefaultRoll              = 0.0
	DefaultSpeed             = 2.5
	DefaultSensitivity       = 0.1
	DefaultOrbitRadius       = 5.0
	DefaultMouseSensitivity  = 1
	DefaultZoomSpeed         = 1
	DefaultMinZoom           = 0.1
	DefaultMaxZoom           = 100
	DefaultFOVSpeed          = 1
	DefaultOrbitSpeed        = 1
	DefaultMinOrbitRadius    = 0.001
	DefaultMaxOrbitRadius    = 10000
)

type Camera struct {
	position  madar.Vector3
	direction madar.Vector3
	up        madar.Vector3
	right     madar.Vector3
	forward   madar.Vector3

	projection ProjectionType
	mode       CameraMode
	width      float32
	height     float32
	fOV        float32
	near       float32
	far        float32
	zoom       float32

	yaw   float32
	pitch float32
	roll  float32

	movementSpeed  float32
	sensitivity    float32
	constrainPitch bool

	viewMatrix       madar.Matrix4X4
	projectionMatrix madar.Matrix4X4

	target      madar.Vector3
	orbitRadius float32

	loggerID string

	controls      CameraControls
	isMouseLocked bool
}

type CreateCameraInfo struct {
	Position    madar.Vector3
	Direction   madar.Vector3
	Up          madar.Vector3
	Target      madar.Vector3
	Projection  ProjectionType
	Mode        CameraMode
	Width       float32
	Height      float32
	FOV         float32
	Near        float32
	Far         float32
	Zoom        float32
	OrbitRadius float32
}

func NewCamera(info CreateCameraInfo) *Camera {
	loggerID := "camera_" + time.Now().Format("150405")
	bayaan.Info("[%s] Initializing new camera instance", loggerID)

	c := &Camera{
		position:       info.Position,
		direction:      info.Direction,
		up:             info.Up,
		projection:     info.Projection,
		mode:           info.Mode,
		width:          info.Width,
		height:         info.Height,
		fOV:            If(info.FOV != 0, info.FOV, DefaultFOV),
		near:           If(info.Near != 0, info.Near, DefaultNear),
		far:            If(info.Far != 0, info.Far, DefaultFar),
		zoom:           If(info.Zoom != 0, info.Zoom, DefaultZoom),
		target:         info.Target,
		orbitRadius:    If(info.OrbitRadius != 0, info.OrbitRadius, DefaultOrbitRadius),
		yaw:            DefaultYaw,
		pitch:          DefaultPitch,
		roll:           DefaultRoll,
		movementSpeed:  DefaultSpeed,
		sensitivity:    DefaultSensitivity,
		constrainPitch: true,
		loggerID:       loggerID,
	}

	if c.direction.IsZero() {
		bayaan.Debug("[%s] Direction not specified, using default direction (0,0,-1)", c.loggerID)
		c.direction = madar.Vector3{X: 0, Y: 0, Z: -1}
	}
	if c.up.IsZero() {
		bayaan.Debug("[%s] Up vector not specified, using default up vector (0,1,0)", c.loggerID)
		c.up = madar.Vector3{X: 0, Y: 1, Z: 0}
	}
	bayaan.Info("[%s] Camera initialized with mode: %v, projection: %v", c.loggerID, c.mode, c.projection)

	c.controls = CameraControls{
		BaseMovementSpeed: 5.0,
		SprintMultiplier:  2.0,
		MouseSensitivity:  0.1,
		ZoomSpeed:         0.1,
		MinZoom:           0.1,
		MaxZoom:           10.0,
		FOVSpeed:          2.0,
		MinFOV:            10.0,
		MaxFOV:            90.0,
		OrbitSpeed:        1.0,
		MinOrbitRadius:    1.0,
		MaxOrbitRadius:    20.0,
	}

	c.Update()
	return c
}

func (c *Camera) Update() {
	bayaan.Trace("[%s] Updating camera state", c.loggerID)

	if c.mode == Free || c.mode == FirstPerson {
		c.updateDirectionFromEuler()
	}

	if c.direction.IsZero() {
		bayaan.Error("[%s] Direction is zero nothing will be darwin on the screen", c.loggerID)
	}

	c.right = c.up.Cross(c.direction).Normalize()
	c.forward = c.direction.Cross(c.right).Normalize()

	switch c.mode {
	case Orbit:
		bayaan.Debug("[%s] Updating orbit camera position", c.loggerID)
		c.updateOrbitCamera()
	case ThirdPerson:
		bayaan.Debug("[%s] Updating third-person camera position", c.loggerID)
		c.updateThirdPersonCamera()
	}

	c.updateViewMatrix()
	c.updateProjectionMatrix()
}

func (c *Camera) updateDirectionFromEuler() {
	bayaan.Trace("[%s] Updating direction from Euler angles - Yaw: %.2f, Pitch: %.2f", c.loggerID, c.yaw, c.pitch)

	if c.constrainPitch {
		if c.pitch > 89.0 {
			bayaan.Debug("[%s] Constraining pitch from %.2f to 89.0", c.loggerID, c.pitch)
			c.pitch = 89.0
		}
		if c.pitch < -89.0 {
			bayaan.Debug("[%s] Constraining pitch from %.2f to -89.0", c.loggerID, c.pitch)
			c.pitch = -89.0
		}

	}

	yawRad := degreesToRadians(c.yaw)
	pitchRad := degreesToRadians(c.pitch)

	// c.direction = madar.Vector3{
	// 	X: float32(math.Cos(float64(yawRad)) * math.Cos(float64(pitchRad))),
	// 	Y: float32(math.Sin(float64(pitchRad))),
	// 	Z: float32(math.Sin(float64(yawRad)) * math.Cos(float64(pitchRad))),
	// }.Normalize()
	c.direction = c.direction.Rotate(yawRad, pitchRad, c.roll).Normalize()
}

func (c *Camera) updateOrbitCamera() {

	c.position = c.target.Add(c.direction.MultiplyScalar(c.orbitRadius))
}

func (c *Camera) updateThirdPersonCamera() {

	idealOffset := c.direction.MultiplyScalar(-c.orbitRadius)
	c.position = c.target.Add(idealOffset)
}

func (c *Camera) updateViewMatrix() {
	c.viewMatrix = madar.LookAtMatrix4X4(c.position, c.position.Add(c.direction), c.up)
}

func (c *Camera) updateProjectionMatrix() {
	aspect := c.width / c.height
	switch c.projection {
	case Perspective:
		c.projectionMatrix = madar.PerspectiveMatrix4X4(c.fOV, aspect, c.near, c.far)
	case Orthographic:
		c.projectionMatrix = madar.OrthographicMatrix4X4(-c.zoom*aspect, c.zoom*aspect, -c.zoom, c.zoom, c.near, c.far)
	}
}

func (c *Camera) ProcessMouseMovement(xOffset, yOffset float32, constrainPitch bool) {
	bayaan.Trace("[%s] Processing mouse movement - X offset: %.2f, Y offset: %.2f", c.loggerID, xOffset, yOffset)

	xOffset *= c.sensitivity
	yOffset *= c.sensitivity

	c.yaw += xOffset
	c.pitch += yOffset

	c.Update()
}

func (c *Camera) ProcessMouseScroll(yOffset float32) {
	bayaan.Trace("[%s] Processing mouse scroll - Y offset: %.2f", c.loggerID, yOffset)

	if c.projection == Perspective {
		oldFOV := c.fOV
		c.fOV -= yOffset
		if c.fOV < 1.0 {
			bayaan.Debug("[%s] FOV constrained from %.2f to 1.0", c.loggerID, c.fOV)
			c.fOV = 1.0
		}
		if c.fOV > 90.0 {
			bayaan.Debug("[%s] FOV constrained from %.2f to 90.0", c.loggerID, c.fOV)
			c.fOV = 90.0
		}
		bayaan.Debug("[%s] FOV changed from %.2f to %.2f", c.loggerID, oldFOV, c.fOV)
	} else {
		oldZoom := c.zoom
		c.zoom -= yOffset * 0.05
		if c.zoom < 0.1 {
			bayaan.Debug("[%s] Zoom constrained from %.2f to 0.1", c.loggerID, c.zoom)
			c.zoom = 0.1
		}
		bayaan.Debug("[%s] Zoom changed from %.2f to %.2f", c.loggerID, oldZoom, c.zoom)
	}
	c.Update()
}

func (c *Camera) SetOrbitRadius(radius float32) {
	c.orbitRadius = radius
	c.Update()
}

func (c *Camera) SetTarget(target madar.Vector3) {
	c.target = target
	c.Update()
}

func (c *Camera) SetConstrainPitch(constrain bool) {
	c.constrainPitch = constrain
	c.Update()
}

func (c *Camera) SetSensitivity(sensitivity float32) {
	c.sensitivity = sensitivity
}

func (c *Camera) SetMovementSpeed(speed float32) {
	c.movementSpeed = speed
}

func degreesToRadians(degrees float32) float32 {
	return degrees * math.Pi / 180.0
}

func (c *Camera) GetMovementSpeed() float32 {
	return c.movementSpeed
}

func (c *Camera) GetSensitivity() float32 {
	return c.sensitivity
}

func (c *Camera) GetTarget() madar.Vector3 {
	return c.target
}

func (c *Camera) GetOrbitRadius() float32 {
	return c.orbitRadius
}

func (c *Camera) GetViewMatrix() madar.Matrix4X4 {
	return c.viewMatrix
}

func (c *Camera) GetProjectionMatrix() madar.Matrix4X4 {
	return c.projectionMatrix
}

func (c *Camera) SetPosition(v madar.Vector3) {
	c.position = v
	c.Update()
}

func (c *Camera) SetDirection(v madar.Vector3) {
	c.direction = v
	c.Update()
}

func (c *Camera) SetUp(v madar.Vector3) {
	c.up = v
	c.Update()
}

func (c *Camera) SetZoom(zoom float32) {
	c.zoom = zoom
	c.Update()
}

func (c *Camera) SetFOV(fov float32) {
	c.fOV = fov
	c.Update()
}

func (c *Camera) SetNear(near float32) {
	c.near = near
	c.Update()
}

func (c *Camera) SetFar(far float32) {
	c.far = far
	c.Update()
}

func (c *Camera) SetWidth(width float32) {
	c.width = width
	c.Update()
}

func (c *Camera) SetHeight(height float32) {
	c.height = height
	c.Update()
}

func (c *Camera) GetPosition() madar.Vector3 {
	return c.position
}

func (c *Camera) GetDirection() madar.Vector3 {
	return c.direction
}

func (c *Camera) GetUp() madar.Vector3 {
	return c.up
}

func (c *Camera) GetZoom() float32 {
	return c.zoom
}

func (c *Camera) GetFOV() float32 {
	return c.fOV
}

func (c *Camera) GetNear() float32 {
	return c.near
}

func (c *Camera) GetFar() float32 {
	return c.far
}

func (c *Camera) GetWidth() float32 {
	return c.width
}

func (c *Camera) GetHeight() float32 {
	return c.height
}

func (c *Camera) GetProjection() ProjectionType {
	return c.projection
}

func (c *Camera) GetMode() CameraMode {
	return c.mode
}

func (c *Camera) Move(direction madar.Vector3) {
	bayaan.Debug("[%s] Moving camera by vector (%.2f, %.2f, %.2f)", c.loggerID, direction.X, direction.Y, direction.Z)
	oldPos := c.position
	c.position = c.position.Add(direction)
	bayaan.Trace("[%s] Position changed from (%.2f, %.2f, %.2f) to (%.2f, %.2f, %.2f)",
		c.loggerID, oldPos.X, oldPos.Y, oldPos.Z, c.position.X, c.position.Y, c.position.Z)
	c.Update()
}

func (c *Camera) MoveForward(distance float32) {
	c.Update()
	c.position = c.position.Add(c.forward.MultiplyScalar(distance))
	c.Update()
}

func (c *Camera) MoveRight(distance float32) {
	c.Update()
	c.position = c.position.Add(c.right.MultiplyScalar(distance))
	c.Update()
}

func (c *Camera) MoveUp(distance float32) {
	c.Update()
	c.position = c.position.Add(c.up.MultiplyScalar(distance))
	c.Update()
}

func (c *Camera) Rotate(direction madar.Vector3) {
	c.Update()
	c.direction = c.direction.Add(direction)
	c.Update()
}

func (c *Camera) LookAt(target madar.Vector3) {
	bayaan.Debug("[%s] Looking at target (%.2f, %.2f, %.2f)", c.loggerID, target.X, target.Y, target.Z)
	oldDir := c.direction
	c.direction = target.Sub(c.position).Normalize()
	bayaan.Trace("[%s] Direction changed from (%.2f, %.2f, %.2f) to (%.2f, %.2f, %.2f)",
		c.loggerID, oldDir.X, oldDir.Y, oldDir.Z, c.direction.X, c.direction.Y, c.direction.Z)
	c.Update()
}

func (c *Camera) SetMode(mode CameraMode) {
	if mode >= CameraMode(5) {
		err := bayaan.Error("[%s] Invalid camera mode: %v", c.loggerID, mode)
		panic(err)
	}

	bayaan.Info("[%s] Changing camera mode from %v to %v", c.loggerID, c.mode, mode)
	c.mode = mode
	c.Update()
}

func (c *Camera) SetProjection(projection ProjectionType) {
	if projection >= ProjectionType(2) {
		err := bayaan.Error("[%s] Invalid projection type: %v", c.loggerID, projection)
		panic(err)
	}

	bayaan.Info("[%s] Changing projection from %v to %v", c.loggerID, c.projection, projection)
	c.projection = projection
	c.Update()
}

// ProcessInput handles all input for the camera

func (c *Camera) ProcessInput(deltaTime float32) {
	c.handleKeyboardInput(deltaTime)
	c.handleMouseInput(deltaTime)
}

func (c *Camera) handleKeyboardInput(deltaTime float32) {
	// Calculate movement speed (sprint if shift is held)
	movementSpeed := c.controls.BaseMovementSpeed
	if input.IsKeyHeld(input.KeyLeftShift) {
		movementSpeed *= c.controls.SprintMultiplier
	}
	movementSpeed *= deltaTime

	// Process movement based on camera mode
	switch c.mode {
	case Free, FirstPerson:
		c.handleFreeMovement(movementSpeed)
	case Orbit:
		c.handleOrbitMovement(movementSpeed)
	case ThirdPerson:
		c.handleThirdPersonMovement(movementSpeed)
	}

	// Toggle mouse lock
	if input.IsKeyPressed(toggleMouseLockKey) {
		c.isMouseLocked = !c.isMouseLocked
		if c.isMouseLocked {
			input.LockMouse()
			bayaan.Debug("[%s] Mouse locked", c.loggerID)
		} else {
			input.UnlockMouse()
			bayaan.Debug("[%s] Mouse unlocked", c.loggerID)
		}
	}

	// Reset camera
	if input.IsKeyPressed(resetCameraKey) {
		c.reset()
	}
}

func (c *Camera) handleFreeMovement(movementSpeed float32) {
	// Forward/Backward
	if input.IsKeyHeld(moveForwardKey) {
		c.MoveForward(movementSpeed)
	}
	if input.IsKeyHeld(moveBackwardKey) {
		c.MoveForward(-movementSpeed)
	}

	// Left/Right
	if input.IsKeyHeld(moveLeftKey) {
		c.MoveRight(-movementSpeed)
	}
	if input.IsKeyHeld(moveRightKey) {
		c.MoveRight(movementSpeed)
	}

	// Up/Down
	if input.IsKeyHeld(moveUpKey) {
		c.MoveUp(movementSpeed)
	}
	if input.IsKeyHeld(moveDownKey) {
		c.MoveUp(-movementSpeed)
	}
}

func (c *Camera) handleOrbitMovement(float32) {
	// Adjust orbit radius with scroll
	scroll := input.GetMouseScroll()
	newRadius := c.orbitRadius - scroll.Y*c.controls.ZoomSpeed
	c.orbitRadius = clamp(newRadius, c.controls.MinOrbitRadius, c.controls.MaxOrbitRadius)

	// Orbit around target when right mouse button is held
	if input.IsMouseButtonHeld(input.MouseRight) {
		mouseDelta := input.GetMouseDelta()
		c.yaw += mouseDelta.X * c.controls.OrbitSpeed * c.controls.MouseSensitivity
		c.pitch += -mouseDelta.Y * c.controls.OrbitSpeed * c.controls.MouseSensitivity
	}
}

func (c *Camera) handleThirdPersonMovement(movementSpeed float32) {
	// Move target instead of camera
	if input.IsKeyHeld(moveForwardKey) {
		c.target = c.target.Add(c.forward.MultiplyScalar(movementSpeed))
	}
	if input.IsKeyHeld(moveBackwardKey) {
		c.target = c.target.Add(c.forward.MultiplyScalar(-movementSpeed))
	}
	if input.IsKeyHeld(moveLeftKey) {
		c.target = c.target.Add(c.right.MultiplyScalar(-movementSpeed))
	}
	if input.IsKeyHeld(moveRightKey) {
		c.target = c.target.Add(c.right.MultiplyScalar(movementSpeed))
	}
}

func (c *Camera) handleMouseInput(deltaTime float32) {
	// Only process mouse look when mouse is locked or left button is held
	if c.isMouseLocked || input.IsMouseButtonHeld(input.MouseLeft) {
		mouseDelta := input.GetMouseDelta()

		// Update camera angles
		c.yaw += mouseDelta.X * c.controls.MouseSensitivity * deltaTime
		c.pitch += -mouseDelta.Y * c.controls.MouseSensitivity * deltaTime

		// Constrain pitch
		c.pitch = clamp(c.pitch, -89.0, 89.0)
	}

	// Handle zoom/FOV with mouse scroll
	scroll := input.GetMouseScroll()
	if c.projection == Perspective {
		newFOV := c.fOV - scroll.Y*c.controls.FOVSpeed*deltaTime
		c.fOV = clamp(newFOV, c.controls.MinFOV, c.controls.MaxFOV)
	} else {
		newZoom := c.zoom - scroll.Y*c.controls.ZoomSpeed*deltaTime
		c.zoom = clamp(newZoom, c.controls.MinZoom, c.controls.MaxZoom)
	}
}

// Helper function for value clamping

func clamp(value, min, max float32) float32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Reset camera to default values

func (c *Camera) reset() {
	bayaan.Info("[%s] Resetting camera to default values", c.loggerID)

	c.position = madar.Vector3{X: 0, Y: 0, Z: 0}
	c.direction = madar.Vector3{X: 0, Y: 0, Z: -1}
	c.up = madar.Vector3{X: 0, Y: 1, Z: 0}
	c.yaw = DefaultYaw
	c.pitch = DefaultPitch
	c.fOV = DefaultFOV
	c.zoom = DefaultZoom
	c.orbitRadius = DefaultOrbitRadius

	c.Update()
}

// Add method to update camera controls

func (c *Camera) SetControls(controls CameraControls) {
	bayaan.Debug("[%s] Updating camera controls", c.loggerID)
	c.controls = CameraControls{
		BaseMovementSpeed: If(controls.BaseMovementSpeed != 0, controls.BaseMovementSpeed, If(c.controls.BaseMovementSpeed != 0, c.controls.BaseMovementSpeed, DefaultBaseMovementSpeed)),
		SprintMultiplier:  If(controls.SprintMultiplier != 0, controls.SprintMultiplier, If(c.controls.SprintMultiplier != 0, c.controls.SprintMultiplier, DefaultSprintMultiplier)),
		MouseSensitivity:  If(controls.MouseSensitivity != 0, controls.MouseSensitivity, If(c.controls.MouseSensitivity != 0, c.controls.MouseSensitivity, DefaultMouseSensitivity)),
		ZoomSpeed:         If(controls.ZoomSpeed != 0, controls.ZoomSpeed, If(c.controls.ZoomSpeed != 0, c.controls.ZoomSpeed, DefaultZoomSpeed)),
		MinZoom:           If(controls.MinZoom != 0, controls.MinZoom, If(c.controls.MinZoom != 0, c.controls.MinZoom, DefaultMinZoom)),
		MaxZoom:           If(controls.MaxZoom != 0, controls.MaxZoom, If(c.controls.MaxZoom != 0, c.controls.MaxZoom, DefaultMaxZoom)),
		FOVSpeed:          If(controls.FOVSpeed != 0, controls.FOVSpeed, If(c.controls.FOVSpeed != 0, c.controls.FOVSpeed, DefaultFOVSpeed)),
		MinFOV:            If(controls.MinFOV != 0, controls.MinFOV, If(c.controls.MinFOV != 0, c.controls.MinFOV, DefaultMinFOV)),
		MaxFOV:            If(controls.MaxFOV != 0, controls.MaxFOV, If(c.controls.MaxFOV != 0, c.controls.MaxFOV, DefaultMaxFOV)),
		OrbitSpeed:        If(controls.OrbitSpeed != 0, controls.OrbitSpeed, If(c.controls.OrbitSpeed != 0, c.controls.OrbitSpeed, DefaultOrbitSpeed)),
		MinOrbitRadius:    If(controls.MinOrbitRadius != 0, controls.MinOrbitRadius, If(c.controls.MinOrbitRadius != 0, c.controls.MinOrbitRadius, DefaultMinOrbitRadius)),
		MaxOrbitRadius:    If(controls.MaxOrbitRadius != 0, controls.MaxOrbitRadius, If(c.controls.MaxOrbitRadius != 0, c.controls.MaxOrbitRadius, DefaultMaxOrbitRadius)),
	}
}
func (c *Camera) Cleanup() {
	bayaan.Info("[%s] Cleaning up camera resources", c.loggerID)
}

func If[T any](condition bool, True, False T) T {
	if condition {
		return True
	}
	return False
}
