package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

type buttonState glfw.Action

var (
	window *glfw.Window
)

func Init(w *glfw.Window) {
	window = w
}
