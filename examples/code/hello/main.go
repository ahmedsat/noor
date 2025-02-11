package main

import (
	"image/color"

	"github.com/ahmedsat/noor"
)

var vertices = []noor.Vertex{
	{Position: [3]float32{0, 0.5, 0}},
	{Position: [3]float32{0.5, -0.5, 0}},
	{Position: [3]float32{-0.5, -0.5, 0}},
}

func main() {
	n := noor.New(800, 600, "Hello, Noor!", color.Black).UnwrapOrPanic()

	// set the background color
	// ? if you don't call this, the background will be black
	n.SetBackground(color.RGBA{R: 0x20, G: 0x30, B: 0x30, A: 0xff})

	mesh := noor.NewMesh(vertices, nil, noor.DrawTriangles)

	for !n.ShouldClose() {
		mesh.Draw()
	}
}
