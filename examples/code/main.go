package main

import (
	"image/color"
	"runtime"

	"github.com/ahmedsat/noor"
)

func init() {
	runtime.LockOSThread()
}

func main() {

	n := noor.New(800, 600, "Hello, Shader!", color.Black).UnwrapOrPanic()
	defer n.Close()
	n.SetBackground(color.RGBA{R: 0x20, G: 0x30, B: 0x30, A: 0xff})

	// // ? notice that if loading our shader fails `noor` will us its default shaders
	// shader := noor.CreateShaderProgramFromFiles(
	// 	"examples/assets/shaders/example.vert",
	// 	"examples/assets/shaders/example.frag",
	// ).UnwrapOrPanic()
	// defer shader.Delete()

	tex, err := noor.NewTextureFromFile("examples/assets/textures/wall.jpg", noor.DefaultTextureParameters())
	if err != nil {
		panic(err)
	}

	mesh := noor.NewMesh(vertices, indices, noor.DrawTriangles)

	obj := noor.NewObject("obj", mesh)
	obj.AddTexture(tex)

	n.AddObject(obj)

	n.Loop(func(deltaTime float32) {
		obj.Rotation.X += deltaTime * 100
		obj.Rotation.Y += deltaTime * 100
		obj.Rotation.Z += deltaTime * 100
	})

}
