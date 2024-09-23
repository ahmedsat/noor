package noor

import (
	"unsafe"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type Vertex struct {
	Position [3]float32
	Color    [3]float32
	TexCoord [2]float32
}

const VertexSize = int(unsafe.Sizeof(Vertex{}))

type Mesh struct {
	Vertices []Vertex
	Indices  []uint32
}

func NewMesh(vertices []Vertex, indices []uint32) (m *Mesh, err error) {
	if !isInitialized {
		return m, unInitializedError
	}

	m = new(Mesh)
	m.Vertices = vertices
	m.Indices = indices

	// Generate and bind the VBO
	var VBO uint32
	gl.GenBuffers(1, &VBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*VertexSize, gl.Ptr(vertices), gl.STATIC_DRAW)

	// Set vertex attribute pointers
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, int32(VertexSize), 0)

	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, int32(VertexSize), 3*4)

	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointerWithOffset(2, 2, gl.FLOAT, false, int32(VertexSize), 6*4)

	// Generate and bind the EBO
	var EBO uint32
	if len(indices) > 0 {
		gl.GenBuffers(1, &EBO)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)
	}

	return
}

type Material struct {
	Shader
	Textures []Texture
}

type Object struct {
	Mesh        *Mesh
	Material    *Material
	VAO         uint32
	ModelMatrix [16]float32 // Transformation for rendering
}

func NewObject(vertices []Vertex, indices []uint32, material *Material) (o *Object, err error) {

	if !isInitialized {
		return o, unInitializedError
	}

	o = new(Object)
	// Generate and bind the VAO
	gl.GenVertexArrays(1, &o.VAO)
	gl.BindVertexArray(o.VAO)

	o.Mesh, err = NewMesh(vertices, indices)
	if err != nil {
		return
	}

	o.Material = material

	// Set model matrix
	o.ModelMatrix = [16]float32{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}

	// setup material shader uniforms
	gl.UseProgram(o.Material.program)
	gl.UniformMatrix4fv(gl.GetUniformLocation(o.Material.program, gl.Str("model\x00")), 1, false, &o.ModelMatrix[0])

	// setup material textures
	for i, texture := range o.Material.Textures {
		gl.Uniform1i(gl.GetUniformLocation(o.Material.program, gl.Str(texture.Name+"\x00")), int32(i))
	}

	return
}

func (o *Object) Draw() (err error) {

	if !isInitialized {
		return unInitializedError
	}

	gl.BindVertexArray(o.VAO)

	for i, texture := range o.Material.Textures {
		gl.ActiveTexture(gl.TEXTURE0 + uint32(i))
		gl.BindTexture(gl.TEXTURE_2D, texture.Handle)
	}

	gl.UseProgram(o.Material.program)

	if len(o.Mesh.Indices) > 0 {
		gl.DrawElements(gl.TRIANGLES, int32(len(o.Mesh.Indices)), gl.UNSIGNED_INT, nil)
	} else {
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(o.Mesh.Vertices)))
	}

	return
}
