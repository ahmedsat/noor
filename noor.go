// search for Subjects:
// - compute shaders

package noor

import (
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var window *glfw.Window

type Options struct {
	Width             int
	Height            int
	Title             string
	IsResizable       bool
	DefaultExtButtons glfw.Key
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

	return
}

func Run(draw func()) {

	for !window.ShouldClose() {
		if internalOptions.DefaultExtButtons != 0x0 && window.GetKey(internalOptions.DefaultExtButtons) == glfw.Press {
			window.SetShouldClose(true)
		}

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		draw()

		window.SwapBuffers()
		glfw.PollEvents()
	}

	glfw.Terminate()
}
