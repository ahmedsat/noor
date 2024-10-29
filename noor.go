package noor

import (
	"errors"
	"image/color"

	"github.com/ahmedsat/noor/input"
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	window         *glfw.Window
	defaultExitKey             = input.KeyEscape
	backGround     color.Color = color.Transparent
)

type InitSettings struct {
	// window settings
	WindowWidth, WindowHeight int
	WindowTitle               string
	WindowResizable           bool
	BackGround                color.Color

	// gl settings
	GLMajorVersion, GLMinorVersion int
	GLCoreProfile, DebugLines      bool
}

// Init initializes the noor library
func Init(st InitSettings) (err error) {

	glfw.Init()

	if st.GLMajorVersion == 0 {
		st.GLMajorVersion = 4
		glfw.WindowHint(glfw.ContextVersionMajor, st.GLMajorVersion)
	}

	if st.GLMinorVersion == 0 {
		st.GLMinorVersion = 5
	}

	glfw.WindowHint(glfw.ContextVersionMinor, st.GLMinorVersion)

	if st.GLCoreProfile {
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	}

	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	glfw.WindowHint(glfw.Resizable, boolToInt(st.WindowResizable))

	if st.WindowWidth == 0 {
		st.WindowWidth = 800
	}

	if st.WindowHeight == 0 {
		st.WindowHeight = 600
	}

	if st.WindowTitle == "" {
		st.WindowTitle = "noor window"
	}

	window, err = glfw.CreateWindow(st.WindowWidth, st.WindowHeight, st.WindowTitle, nil, nil)
	if err != nil {
		err = errors.Join(err, errors.New("failed to create window"))
		return
	}

	input.Init(window)
	window.MakeContextCurrent()

	err = gl.Init()
	if err != nil {
		err = errors.Join(err, errors.New("failed to init open-gl"))
		return
	}

	gl.Viewport(0, 0, int32(st.WindowWidth), int32(st.WindowHeight))
	window.SetFramebufferSizeCallback(func(w *glfw.Window, width, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	gl.Enable(gl.DEPTH_TEST)

	if st.DebugLines {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	}

	if st.BackGround != nil {
		backGround = st.BackGround
	}

	return
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func IsWindowShouldClose() bool {
	window.SwapBuffers()
	glfw.PollEvents()
	if input.IsKeyPressed(defaultExitKey) {
		window.SetShouldClose(true)
	}
	r, g, b, a := backGround.RGBA()
	gl.ClearColor(float32(r)/0xffff, float32(g)/0xffff, float32(b)/0xffff, float32(a)/0xffff)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	return window.ShouldClose()
}

func SetDefaultExitKey(key input.Key) { defaultExitKey = key }

func Terminate() {
	window.Destroy()
	glfw.Terminate()
}

func Window() *glfw.Window {
	return window
}
