package noor

import (
	"fmt"
	"image/color"
	"runtime"

	"github.com/ahmedsat/bayaan"
	"github.com/ahmedsat/noor/input"
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	window         *glfw.Window
	defaultExitKey             = input.KeyEscape
	backGround     color.Color = color.Transparent
)

type LogLevel int

const (
	LogDebug LogLevel = iota
	LogInfo
	LogWarning
	LogError
	LogFatal
)

type InitSettings struct {
	WindowWidth, WindowHeight int
	WindowTitle               string
	WindowResizable           bool
	BackGround                color.Color

	GLMajorVersion, GLMinorVersion int
	GLCoreProfile, DebugLines      bool

	EnableMultiSampling bool
	VSyncEnabled        bool
	CursorDisabled      bool
}

func Init(st InitSettings) (err error) {
	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		return fmt.Errorf("failed to initialize GLFW: %w", err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, defaultInt(st.GLMajorVersion, 4))
	glfw.WindowHint(glfw.ContextVersionMinor, defaultInt(st.GLMinorVersion, 5))

	if st.GLCoreProfile {
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	}
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	glfw.WindowHint(glfw.Resizable, boolToInt(st.WindowResizable))

	if st.EnableMultiSampling {
		glfw.WindowHint(glfw.Samples, 4)
	}

	width := defaultInt(st.WindowWidth, 800)
	height := defaultInt(st.WindowHeight, 600)
	title := defaultString(st.WindowTitle, "noor window")

	window, err = glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create window: %w", err)
	}

	window.MakeContextCurrent()

	if st.VSyncEnabled {
		bayaan.Debug("Enabling VSync", bayaan.Fields{
			"swapInterval": 1,
		})
		glfw.SwapInterval(1)
	}

	input.Init(window)

	if st.CursorDisabled {
		window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	}

	if err := gl.Init(); err != nil {
		return fmt.Errorf("failed to initialize OpenGL: %w", err)
	}

	gl.Viewport(0, 0, int32(width), int32(height))
	window.SetFramebufferSizeCallback(func(w *glfw.Window, width, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
		bayaan.Trace("Window resized", bayaan.Fields{
			"width":  width,
			"height": height,
		})

	})

	gl.Enable(gl.DEPTH_TEST)

	if st.EnableMultiSampling {
		gl.Enable(gl.MULTISAMPLE)
	}

	if st.DebugLines {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		bayaan.Trace("Polygon mode set to lines", bayaan.Fields{
			"mode": gl.LINE,
		})
	}

	if st.BackGround != nil {
		backGround = st.BackGround
		bayaan.Trace("Background color set", bayaan.Fields{
			"color": backGround,
		})

	}
	return nil
}

func IsWindowShouldClose() bool {
	window.SwapBuffers()
	glfw.PollEvents()

	if input.IsKeyPressed(defaultExitKey) {
		window.SetShouldClose(true)
	}

	r, g, b, a := backGround.RGBA()
	gl.ClearColor(
		float32(r)/0xffff,
		float32(g)/0xffff,
		float32(b)/0xffff,
		float32(a)/0xffff,
	)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	return window.ShouldClose()
}

func defaultInt(value, defaultVal int) int {
	if value == 0 {
		return defaultVal
	}
	return value
}

func defaultString(value, defaultVal string) string {
	if value == "" {
		return defaultVal
	}
	return value
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func SetDefaultExitKey(key input.Key) {
	defaultExitKey = key
	bayaan.Info("Default exit key changed", bayaan.Fields{
		"key": key,
	})
}

func Terminate() {
	bayaan.Info("Terminating Noor package", bayaan.Fields{})
	window.Destroy()
	glfw.Terminate()
}

func Window() *glfw.Window {
	return window
}
