package noor

import (
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type DrawMode uint32

const (
	DrawTriangles DrawMode = gl.TRIANGLES
	DrawLines     DrawMode = gl.LINES
	DrawPoints    DrawMode = gl.POINTS
)

type Vertex struct {
	Position [3]float32
	Color    [3]float32
	UV       [2]float32
	Normal   [3]float32
}

func NewVertex(position [3]float32, color [3]float32, uv [2]float32, normal [3]float32) Vertex {
	return Vertex{Position: position, Color: color, UV: uv, Normal: normal}
}

func setVertexAttributes() {
	Assert(unsafe.Sizeof(Vertex{}) == 44, "Vertex size must be 44 bytes")
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.EnableVertexAttribArray(2)
	gl.EnableVertexAttribArray(3)

	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 44, 0)
	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 44, 12)
	gl.VertexAttribPointerWithOffset(2, 2, gl.FLOAT, false, 44, 24)
	gl.VertexAttribPointerWithOffset(3, 3, gl.FLOAT, false, 44, 32)
}

type Mesh struct {
	VAO          uint32
	VBO          uint32
	EBO          uint32
	DrawMode     DrawMode
	Count        int32
	DrawElements bool
}

func NewMesh(vertices []Vertex, indices []uint32, drawMode DrawMode) *Mesh {

	var vao, vbo, ebo uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	count := int32(len(vertices))

	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*44, gl.Ptr(vertices), gl.STATIC_DRAW)

	if len(indices) > 0 {
		gl.GenBuffers(1, &ebo)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)
		count = int32(len(indices))
	}

	setVertexAttributes()

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	return &Mesh{
		VAO:      vao,
		VBO:      vbo,
		EBO:      ebo,
		DrawMode: drawMode,
		Count:    count,
	}
}

func (m *Mesh) Delete() {
	gl.DeleteBuffers(1, &m.VAO)
	gl.DeleteBuffers(1, &m.VBO)
	if m.Count > 0 {
		gl.DeleteBuffers(1, &m.EBO)
	}
}

func (m *Mesh) Draw() {
	gl.BindVertexArray(m.VAO)

	if m.DrawElements {
		gl.DrawElements(uint32(m.DrawMode), m.Count, gl.UNSIGNED_INT, nil)
	} else {
		gl.DrawArrays(uint32(m.DrawMode), 0, m.Count)
	}
	gl.BindVertexArray(0)
}
