package noor

import (
	"github.com/ahmedsat/bayaan"
	"github.com/ahmedsat/madar"
	"github.com/go-gl/gl/v4.5-core/gl"
)

type Object struct {
	Handle      uint32
	VertexCount int
	IndexCount  int
	Shader      Shader
	Textures    []Texture

	Position, Scale, Rotation madar.Vector3
}

type Vertex struct {
	// position
	Vx, Vy, Vz float32

	// normal
	Nx, Ny, Nz float32

	// texture coordinates
	U, V float32
}

var sizes = []uint32{
	3, // positions
	3, // normals
	2, // textures coordinates
}

const (
	vertexSize = 32
)

// declared as variable to force the compiler to inline it
// and declared here for mr to not forget to edit it after editing the struct
var setAttr = func() {
	offset := 0
	for i, s := range sizes {
		gl.EnableVertexAttribArray(uint32(i))
		gl.VertexAttribPointerWithOffset(uint32(i), int32(s), gl.FLOAT, false, vertexSize, uintptr(offset)*4)
		offset += int(s)
	}
}

func NewObject(vertices []Vertex, indices []uint32, sh Shader, textures ...Texture) (obj *Object) {

	if len(vertices) == 0 {
		bayaan.Error("vertices cannot be nil or empty")
		return
	}

	obj = &Object{}

	obj.IndexCount = len(indices)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, vertexSize*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	if obj.IndexCount > 0 {
		var ebo uint32
		gl.GenBuffers(1, &ebo)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)
	}

	setAttr()

	obj.Handle = vao
	obj.VertexCount = len(vertices)
	obj.IndexCount = len(indices)
	obj.Shader = sh
	obj.Textures = textures

	obj.Scale = madar.Vector3{X: 1, Y: 1, Z: 1}

	return
}

func (obj *Object) Delete() {
	gl.DeleteVertexArrays(1, &obj.Handle)
}

func (obj *Object) Draw(c *Camera) {

	obj.Shader.Activate()
	obj.Shader.SetUniformMatrix4fv("uModel", obj.ModelMatrix())
	obj.Shader.SetUniformMatrix4fv("uView", c.ViewMatrix)
	obj.Shader.SetUniformMatrix4fv("uProjection", c.ProjectionMatrix)

	for i, tex := range obj.Textures {
		tex.Activate(obj.Shader, uint32(i))
	}

	gl.BindVertexArray(obj.Handle)
	if obj.IndexCount > 0 {
		gl.DrawElements(gl.TRIANGLES, int32(obj.IndexCount), gl.UNSIGNED_INT, gl.Ptr(nil))
	} else {
		gl.DrawArrays(gl.TRIANGLES, 0, int32(obj.VertexCount/5))
	}
}

func (obj *Object) ModelMatrix() [16]float32 {
	return madar.IdentityMatrix4X4().
		Multiply(madar.TranslationMatrix4X4(obj.Position.X, obj.Position.Y, obj.Position.Z)).
		Multiply(madar.ScaleMatrix4X4(obj.Scale.X, obj.Scale.Y, obj.Scale.Z)).
		Multiply(madar.RotationMatrix4X4(obj.Rotation.X, obj.Rotation.Y, obj.Rotation.Z))
}
