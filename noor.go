package noor

import (
	"errors"
	"time"

	"github.com/ahmedsat/noor/input"
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// noor internal state
var (
	isInitialized = false
	window        *glfw.Window
)

// Init initializes the noor library
func Init(width, height int, title string, resizable bool) (err error) {

	glfw.Init()
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 5)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Resizable, boolToInt(resizable))

	window, err = glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return errors.Join(err, errors.New("failed to create window"))
	}
	input.Init(window)
	window.MakeContextCurrent()

	err = gl.Init()
	if err != nil {
		return errors.Join(err, errors.New("failed to init opengl"))
	}

	gl.Viewport(0, 0, int32(width), int32(height))
	window.SetFramebufferSizeCallback(func(w *glfw.Window, width, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	gl.Enable(gl.DEPTH_TEST)

	// ! for debugging
	// gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	isInitialized = true
	return
}

// terminate the noor library
func Terminate() {
	window.Destroy()
	glfw.Terminate()

	isInitialized = false
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func GetWindow() *glfw.Window {
	return window
}

func Run(draw func(), update func(float32), updateInterval time.Duration) {

	lastUpdateTime := time.Now()
	var deltaTime float32

	for !window.ShouldClose() {
		now := time.Now()

		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}

		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Draw as fast as possible
		draw()

		window.SwapBuffers()
		glfw.PollEvents()

		// Update every X seconds
		if now.Sub(lastUpdateTime) >= updateInterval {
			deltaTime = float32(now.Sub(lastUpdateTime).Seconds())
			lastUpdateTime = now

			update(deltaTime)
		}
	}
}
