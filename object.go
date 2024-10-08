package noor

import (
	"github.com/ahmedsat/bayaan"
	"github.com/ahmedsat/madar"
	"github.com/go-gl/gl/v4.5-core/gl"
)

type Object struct {
	handle      uint32
	vertexCount int
	indexCount  int
	shader      Shader
	textures    []Texture

	position, scale, rotation [3]float32
}

// vertices is [x, y, z, u, v]
func NewObject(vertices []float32, indices []uint32, sh Shader, textures ...Texture) (obj *Object) {

	if len(vertices) == 0 {
		bayaan.Error("vertices cannot be nil or empty")
		return
	}

	obj = &Object{}

	obj.indexCount = len(indices)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	if obj.indexCount > 0 {
		var ebo uint32
		gl.GenBuffers(1, &ebo)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)
	}

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 5*4, 0)

	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*4, 3*4)

	obj.handle = vao
	obj.vertexCount = len(vertices)
	obj.indexCount = len(indices)
	obj.shader = sh
	obj.textures = textures

	return
}

func (obj *Object) Delete() {
	gl.DeleteVertexArrays(1, &obj.handle)
}

func (obj *Object) Draw(c *Camera) {

	obj.shader.Activate()
	obj.shader.SetUniformMatrix4fv("model", obj.ModelMatrix())
	obj.shader.SetUniformMatrix4fv("view", c.ViewMatrix)
	obj.shader.SetUniformMatrix4fv("projection", c.ProjectionMatrix)

	for i, tex := range obj.textures {
		tex.Activate(obj.shader, uint32(i))
	}

	gl.BindVertexArray(obj.handle)
	if obj.indexCount > 0 {
		gl.DrawElements(gl.TRIANGLES, int32(obj.indexCount), gl.UNSIGNED_INT, gl.Ptr(nil))
	} else {
		gl.DrawArrays(gl.TRIANGLES, 0, int32(obj.vertexCount/5))
	}
}

func (obj *Object) ModelMatrix() [16]float32 {
	return madar.IdentityMatrix4X4().
		Multiply(madar.TranslationMatrix4X4(obj.position[0], obj.position[1], obj.position[2])).
		Multiply(madar.ScaleMatrix4X4(obj.scale[0], obj.scale[1], obj.scale[2])).
		Multiply(madar.RotationMatrix4X4(obj.rotation[0], obj.rotation[1], obj.rotation[2]))
}

func (obj *Object) SetPosition(x, y, z float32) {
	obj.position = [3]float32{x, y, z}
}

func (obj *Object) SetScale(x, y, z float32) {
	obj.scale = [3]float32{x, y, z}
}

func (obj *Object) SetRotation(x, y, z float32) {
	obj.rotation = [3]float32{x, y, z}
}

// func (obj *Object) SetRotationEuler(x, y, z float32) {
// 	obj.rotation = [3]float32{x, y, z}
// }
