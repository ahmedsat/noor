package input

import (
	"sync"
	"time"

	"github.com/ahmedsat/madar"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	mousePosition     madar.Vector2
	lastMousePosition madar.Vector2
	mouseScroll       madar.Vector2

	mutex           sync.RWMutex
	doubleClickTime = 300 * time.Millisecond
)

type Mouse glfw.MouseButton

var (
	MouseLeft   = Mouse(glfw.MouseButtonLeft)
	MouseRight  = Mouse(glfw.MouseButtonRight)
	MouseMiddle = Mouse(glfw.MouseButtonMiddle)
)

func mouseButtonCallback(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, _ glfw.ModifierKey) {
	mutex.Lock()
	defer mutex.Unlock()

	state := mouseButtonStates[Mouse(button)]
	state.pressed = action == glfw.Press
	state.justPressed = action == glfw.Press
	state.justReleased = action == glfw.Release

	if action == glfw.Press {
		state.lastPressTime = time.Now()
	}

	mouseButtonStates[Mouse(button)] = state
}

func mousePosCallback(_ *glfw.Window, x, y float64) {
	mutex.Lock()
	defer mutex.Unlock()

	lastMousePosition = mousePosition
	mousePosition = madar.Vector2{X: float32(x), Y: float32(y)}
}

func mouseScrollCallback(_ *glfw.Window, x, y float64) {
	mutex.Lock()
	defer mutex.Unlock()

	mouseScroll = madar.Vector2{X: float32(x), Y: float32(y)}
}

func keyCallback(_ *glfw.Window, key glfw.Key, _ int, action glfw.Action, _ glfw.ModifierKey) {
	mutex.Lock()
	defer mutex.Unlock()

	state := keyStates[key]
	state.pressed = action == glfw.Press || action == glfw.Repeat
	state.justPressed = action == glfw.Press
	state.justReleased = action == glfw.Release

	if action == glfw.Press {
		state.lastPressTime = time.Now()
	}

	keyStates[key] = state
}

func IsMouseButtonPressed(button Mouse) bool {
	mutex.RLock()
	defer mutex.RUnlock()
	return mouseButtonStates[button].justPressed
}

func IsMouseButtonReleased(button Mouse) bool {
	mutex.RLock()
	defer mutex.RUnlock()
	return mouseButtonStates[button].justReleased
}

func IsMouseButtonHeld(button Mouse) bool {
	mutex.RLock()
	defer mutex.RUnlock()
	return mouseButtonStates[button].pressed
}

func IsMouseButtonDoubleClicked(button Mouse) bool {
	mutex.RLock()
	defer mutex.RUnlock()
	state := mouseButtonStates[button]
	if state.justPressed && time.Since(state.lastPressTime) < doubleClickTime {
		return true
	}
	return false
}

func GetMousePosition() madar.Vector2 {
	mutex.RLock()
	defer mutex.RUnlock()
	return mousePosition
}

func GetMouseDelta() madar.Vector2 {
	mutex.RLock()
	defer mutex.RUnlock()
	r := madar.Vector2{
		X: mousePosition.X - lastMousePosition.X,
		Y: mousePosition.Y - lastMousePosition.Y,
	}
	lastMousePosition = mousePosition
	return r
}

func GetMouseScroll() madar.Vector2 {
	mutex.RLock()
	defer mutex.RUnlock()
	return mouseScroll
}

func LockMouse() {
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
}

func UnlockMouse() {
	window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
}

func Update() {
	mutex.Lock()
	defer mutex.Unlock()

	// Reset "just" states
	for button := range mouseButtonStates {
		state := mouseButtonStates[button]
		state.justPressed = false
		state.justReleased = false
		mouseButtonStates[button] = state
	}

	for key := range keyStates {
		state := keyStates[key]
		state.justPressed = false
		state.justReleased = false
		keyStates[key] = state
	}

	mouseScroll = madar.Vector2{}
}
