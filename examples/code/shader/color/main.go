package main

import (
	"image/color"

	"github.com/ahmedsat/noor"
)

var vertices = []noor.Vertex{
	{Position: [3]float32{0, 0.5, 0}, Color: [3]float32{1, 0, 0}},
	{Position: [3]float32{0.5, -0.5, 0}, Color: [3]float32{0, 1, 0}},
	{Position: [3]float32{-0.5, -0.5, 0}, Color: [3]float32{0, 0, 1}},
}

func main() {
	n := noor.New(800, 600, "Hello, Shader!", color.Black).UnwrapOrPanic()

	// we are loading our shaders from files
	// you can also load them from strings with `noor.CreateShaderProgram`
	// ? notice that if loading our shader fails `noor` will us its default shaders
	shader := noor.CreateShaderProgramFromFiles("examples/assets/shaders/color.vert", "examples/assets/shaders/color.frag").UnwrapOrPanic()

	n.Shader.Delete() // delete the default shader
	n.Shader = shader // and replace it with our shader

	mesh := noor.NewMesh(vertices, nil, noor.DrawTriangles)

	for !n.ShouldClose() {
		mesh.Draw()
	}
}
