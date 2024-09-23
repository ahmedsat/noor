// search for Subjects:
// - compute shaders

package noor

import (
	"fmt"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	window             *glfw.Window
	isInitialized      bool
	unInitializedError = fmt.Errorf("noor is not initialized")
)

type Options struct {
	Width             int
	Height            int
	Title             string
	IsResizable       bool
	DefaultExtButtons glfw.Key

	Background [3]float32
}

func (opts *Options) UpdateOptions(op Options) {
	if op.Width != 0 {
		opts.Width = op.Width
	}

	if op.Height != 0 {
		opts.Height = op.Height
	}

	if op.Title != "" {
		opts.Title = op.Title
	}

	if op.IsResizable {
		opts.IsResizable = op.IsResizable
	}

	if op.DefaultExtButtons != 0x0 {
		opts.DefaultExtButtons = op.DefaultExtButtons
	}

	if op.Background[0] != 0.0 || op.Background[1] != 0.0 || op.Background[2] != 0.0 {
		opts.Background = op.Background
	}

}

var internalOptions = Options{
	Width:             800,
	Height:            600,
	Title:             "Noor",
	IsResizable:       false,
	DefaultExtButtons: glfw.KeyEscape,
}

func Init(opts Options) (err error) {

	internalOptions.UpdateOptions(opts)

	if err = glfw.Init(); err != nil {
		return
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 5)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	if internalOptions.IsResizable {
		glfw.WindowHint(glfw.Resizable, glfw.True)
	} else {
		glfw.WindowHint(glfw.Resizable, glfw.False)
	}

	window, err = glfw.CreateWindow(internalOptions.Width, internalOptions.Height, internalOptions.Title, nil, nil)
	if err != nil {
		return
	}

	window.MakeContextCurrent()

	if err = gl.Init(); err != nil {
		return
	}

	// version := gl.GoStr(gl.GetString(gl.VERSION))

	gl.Viewport(0, 0, int32(internalOptions.Width), int32(internalOptions.Height))

	window.SetFramebufferSizeCallback(func(w *glfw.Window, width, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	isInitialized = true
	return
}

func Run(draw func()) (err error) {

	if !isInitialized {
		return unInitializedError
	}

	for !window.ShouldClose() {
		if internalOptions.DefaultExtButtons != 0x0 && window.GetKey(internalOptions.DefaultExtButtons) == glfw.Press {
			window.SetShouldClose(true)
		}

		gl.ClearColor(internalOptions.Background[0], internalOptions.Background[1], internalOptions.Background[2], 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		draw()

		window.SwapBuffers()
		glfw.PollEvents()
	}

	glfw.Terminate()

	return
}
