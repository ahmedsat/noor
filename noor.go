package noor

import (
	"image/color"
	"os"
	"strings"
	"unsafe"

	"github.com/ahmedsat/bayaan"
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func ternary[T any](condition bool, ifTrue T, ifFalse T) T {
	if condition {
		return ifTrue
	}
	return ifFalse
}

type state struct {
	window         *glfw.Window
	defaultExitKey glfw.Key
	backGround     color.Color
}

var st state

func init() {
	st.defaultExitKey = glfw.KeyEscape
	st.backGround = color.Transparent
}

func Init(width int, height int, title string, backGround color.Color) (err error) {
	if bayaan.GetLevel() == bayaan.LoggerLevelDebug {
		os.Setenv("NOOR_DEBUG_MODE", "normal debug")
	}

	bayaan.Trace("Initializing Noor package", bayaan.Fields{
		"width":  width,
		"height": height,
		"title":  title,
		"bg":     ternary(backGround != nil, backGround, st.backGround),
	})

	st.backGround = ternary(backGround != nil, backGround, st.backGround)

	if err = initGLFW(width, height, title); err != nil {
		return
	}

	if err = initGL(); err != nil {
		return
	}
	return
}

func initGLFW(width int, height int, title string) (err error) {
	bayaan.Info("Initializing GLFW", bayaan.Fields{
		"width":  width,
		"height": height,
		"title":  title,
	})
	glfw.Init()
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 5)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.True)

	st.window, err = glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return
	}

	bayaan.Info("Created GLFW window successfully", bayaan.Fields{})

	st.window.MakeContextCurrent()
	return
}

func initGL() (err error) {
	bayaan.Info("Initializing OpenGL", bayaan.Fields{})
	if err = gl.Init(); err != nil {
		return
	}
	st.window.SetFramebufferSizeCallback(func(w *glfw.Window, width, height int) {
		bayaan.Trace("Window resized", bayaan.Fields{
			"width":  width,
			"height": height,
		})
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	debug := strings.ToLower(os.Getenv("NOOR_DEBUG_MODE"))

	if debug != "" {
		bayaan.Warn("Enabling OpenGL debug", bayaan.Fields{
			"debug": debug,
		})
		switch debug {
		case "lines":
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		case "points":
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.POINT)
		default:
		}
		gl.Enable(gl.DEBUG_OUTPUT)
		gl.DebugMessageCallback(func(
			source, gltype, id, severity uint32,
			length int32,
			message string,
			userParam unsafe.Pointer) {
			bayaan.Debug("OpenGL debug message", bayaan.Fields{
				"source":    source,
				"gltype":    gltype,
				"id":        id,
				"severity":  severity,
				"length":    length,
				"message":   message,
				"userParam": userParam,
			})
		}, nil)

	}

	gl.Enable(gl.DEPTH_TEST)
	return
}

func ShouldClose() bool {
	if st.window == nil {
		return true
	}

	if st.window.GetKey(glfw.Key(st.defaultExitKey)) == glfw.Press {
		return true
	}

	glfw.PollEvents()
	st.window.SwapBuffers()

	r, g, b, a := st.backGround.RGBA()
	gl.ClearColor(float32(r)/0xffff, float32(g)/0xffff, float32(b)/0xffff, float32(a)/0xffff)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	return st.window.ShouldClose()
}

func Close() (err error) {
	st.window.Destroy()
	bayaan.Info("Window destroyed", bayaan.Fields{})
	glfw.Terminate()
	bayaan.Info("GLFW terminated", bayaan.Fields{})

	return
}
