package noor

import (
	"github.com/go-gl/gl/v4.5-core/gl"
)

type Mesh struct {
	vao uint32
	drawingMethod
	verticesCount int32
}

type drawingMethod uint8

const (
	drawArrays drawingMethod = iota
	drawElements
)

type CreateMeshInfo struct {
	Vertices []float32
	Indices  []uint32
	Sizes    []int32
}

func CreateMesh(info CreateMeshInfo) (m Mesh) {

	vertexSize := int32(0)
	for _, v := range info.Sizes {
		vertexSize += v
	}

	m.verticesCount = int32(len(info.Vertices)) / vertexSize

	gl.GenVertexArrays(1, &m.vao)
	gl.BindVertexArray(m.vao)

	var vbo, ebo uint32
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(info.Vertices)*4, gl.Ptr(info.Vertices), gl.STATIC_DRAW)

	if len(info.Indices) != 0 {
		m.drawingMethod = drawElements
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(info.Indices)*4, gl.Ptr(info.Indices), gl.STATIC_DRAW)
		m.verticesCount = int32(len(info.Indices))
		m.drawingMethod = drawElements
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

	switch m.drawingMethod {
	case drawArrays:
		gl.DrawArrays(gl.TRIANGLES, 0, m.verticesCount)
	case drawElements:
		gl.DrawElements(gl.TRIANGLES, m.verticesCount, gl.UNSIGNED_INT, nil)
	}

}
