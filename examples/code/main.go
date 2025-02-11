package main

import (
	"image/color"
	"runtime"

	"github.com/ahmedsat/noor"
)

func init() {
	runtime.LockOSThread()
}

var vertices = []noor.Vertex{
	{Position: [3]float32{0, 0.5, 0}, Color: [3]float32{1, 0, 0}, UV: [2]float32{0.5, 1}},
	{Position: [3]float32{0.5, -0.5, 0}, Color: [3]float32{0, 1, 0}, UV: [2]float32{1, 0}},
	{Position: [3]float32{-0.5, -0.5, 0}, Color: [3]float32{0, 0, 1}, UV: [2]float32{0, 0}},
}

func main() {

	n := noor.New(800, 600, "Hello, Shader!", color.Black).UnwrapOrPanic()
	defer n.Close()
	n.SetBackground(color.RGBA{R: 0x20, G: 0x30, B: 0x30, A: 0xff})

	// ? notice that if loading our shader fails `noor` will us its default shaders
	// // shader := noor.CreateShaderProgramFromFiles(
	// // 	"examples/assets/shaders/example.vert",
	// // 	"examples/assets/shaders/example.frag",
	// // ).UnwrapOrPanic()
	// // defer shader.Delete()
	// // n.SetShader(shader)

	mesh := noor.NewMesh(vertices, nil, noor.DrawTriangles)

	for !n.ShouldClose() {
		mesh.Draw()
	}
}
