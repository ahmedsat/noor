package noor

import (
	"github.com/go-gl/gl/v4.5-core/gl"
)

type Mesh struct {
	vao      uint32
	vertices []float32

	// Position madar.Vector3
}

type CreateMeshInfo struct {
	Vertices []float32
	Sizes    []int32
}

func CreateMesh(info CreateMeshInfo) (m Mesh) {
	m = Mesh{
		vertices: info.Vertices,
	}

	gl.GenVertexArrays(1, &m.vao)
	gl.BindVertexArray(m.vao)

	vbo := uint32(0)
	gl.GenBuffers(1, &vbo)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(m.vertices)*4, gl.Ptr(m.vertices), gl.STATIC_DRAW)

	vertexSize := int32(0)
	for _, v := range info.Sizes {
		vertexSize += v
	}

	setAttributes(info.Sizes, vertexSize)

	return
}

func setAttributes(sizes []int32, vertexSize int32) {
	offset := uintptr(0)
	for i := 0; i < len(sizes); i++ {
		gl.VertexAttribPointerWithOffset(uint32(i), sizes[i], gl.FLOAT, false, vertexSize*4, offset*4)
		gl.EnableVertexAttribArray(uint32(i))
		offset += uintptr(sizes[i])
	}
}

func (m *Mesh) Draw() {
	gl.BindVertexArray(m.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
}
