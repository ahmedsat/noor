package main

import (
	"image/color"
	"runtime"

	"github.com/ahmedsat/madar"
	"github.com/ahmedsat/noor"
)

func init() {
	runtime.LockOSThread()
}

func main() {

	n := noor.New(800, 600, "Hello, Shader!", color.Black).UnwrapOrPanic()
	defer n.Close()
	n.SetBackground(color.RGBA{R: 0x20, G: 0x30, B: 0x30, A: 0xff})

	// ? notice that if loading our shader fails `noor` will us its default shaders
	shader := noor.CreateShaderProgramFromFiles(
		"examples/assets/shaders/example.vert",
		"examples/assets/shaders/example.frag",
	).UnwrapOrPanic()
	defer shader.Delete()
	n.SetShader(shader)

	mesh := noor.NewMesh(vertices, indices, noor.DrawTriangles)

	r := float32(0)
	var mat madar.Matrix

	for !n.ShouldClose() {

		r += 1
		if r > 360 {
			r = 0
		}

		mat = madar.RotationMatrix(r, r, r)
		n.Shader.SetUniformMatrixFloat32("uMat", mat.GL())

		mesh.Draw()
	}
}
