package input

import (
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
)

type buttonState struct {
	pressed       bool
	justPressed   bool
	justReleased  bool
	lastPressTime time.Time
}

var (
	window            *glfw.Window
	mouseButtonStates map[Mouse]buttonState
	keyStates         map[glfw.Key]buttonState
)

func Init(w *glfw.Window) {
	window = w
	mouseButtonStates = make(map[Mouse]buttonState)
	keyStates = make(map[glfw.Key]buttonState)

	window.SetMouseButtonCallback(mouseButtonCallback)
	window.SetCursorPosCallback(mousePosCallback)
	window.SetScrollCallback(mouseScrollCallback)
	window.SetKeyCallback(keyCallback)
}
