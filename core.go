package noor

import (
	_ "embed"
	"fmt"
	"image/color"
	"runtime"
	"strings"
	"time"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// embed default shaders in the binary using go:embed
//
//go:embed assets/shaders/default.vert
var DefaultVertexShader string

//go:embed assets/shaders/default.frag
var DefaultFragmentShader string

type Noor struct {
	*glfw.Window
	*Scene
}

func New(width, height int, title string, bg color.Color) Result[Noor] {

	if !IsLockedToThread() {
		fmt.Println("WARNING: The current goroutine is not locked to a thread. This may cause issues with the OpenGL context.")
		fmt.Println("         Use noor.IsLockedToThread() to check if the current goroutine is locked to a thread or not.")
		fmt.Println("         If you are not sure what this means, Just type runtime.LockOSThread() before calling noor.New().")
	}

	noor := Noor{}

	noor.Scene = NewScene()

	var err error

	if err = glfw.Init(); err != nil {
		return Err[Noor](err)
	}

	setWindowHints()

	noor.Window, err = glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return Err[Noor](err)
	}

	noor.Window.MakeContextCurrent()

	noor.Window.SetInputMode(glfw.StickyKeysMode, glfw.True)

	if err = gl.Init(); err != nil {
		return Err[Noor](err)
	}

	gl.Enable(gl.DEPTH_TEST)

	r, g, b, a := bg.RGBA()
	gl.ClearColor(float32(r)/float32(0xffff), float32(g)/float32(0xffff), float32(b)/float32(0xffff), float32(a)/float32(0xffff))

	return Ok[Noor](noor)
}

func setWindowHints() {
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
}

func (n *Noor) Loop(update func(float32)) {

	lastFrameTime := time.Now()

	for !n.Window.ShouldClose() {
		currentFrameTime := time.Now()
		deltaTime := currentFrameTime.Sub(lastFrameTime).Seconds()
		lastFrameTime = currentFrameTime

		if n.Window.GetKey(glfw.KeyEscape) == glfw.Press {
			n.Window.SetShouldClose(true)
		}

		update(float32(deltaTime))

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		n.Render()

		glfw.PollEvents()
		n.Window.SwapBuffers()

	}

}

func (n *Noor) SetBackground(bg color.Color) {
	r, g, b, a := bg.RGBA()
	gl.ClearColor(float32(r)/float32(0xffff), float32(g)/float32(0xffff), float32(b)/float32(0xffff), float32(a)/float32(0xffff))
}

func (n *Noor) Close() {

	n.Window.SetShouldClose(true)

	n.Window.Destroy()
	glfw.Terminate()
}

func IsLockedToThread() bool {
	buf := make([]byte, 1<<16)
	n := runtime.Stack(buf, false)
	stack := string(buf[:n])

	return strings.Contains(stack, "locked to thread")
}
