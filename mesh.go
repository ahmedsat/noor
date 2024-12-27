package noor

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/ahmedsat/bayaan"
	"github.com/go-gl/gl/v4.5-core/gl"
)

// Mesh represents a 3D mesh with vertex data and rendering information
type Mesh struct {
	vao           uint32
	vbo           uint32
	ebo           uint32
	drawingMethod drawingMethod
	verticesCount int32
	primitiveType uint32
}

type drawingMethod uint8

const (
	drawArrays drawingMethod = iota
	drawElements
)

// PrimitiveType defines how vertices should be interpreted for rendering
type PrimitiveType uint32

const (
	Triangles     PrimitiveType = gl.TRIANGLES
	TriangleStrip PrimitiveType = gl.TRIANGLE_STRIP
	TriangleFan   PrimitiveType = gl.TRIANGLE_FAN
	Lines         PrimitiveType = gl.LINES
	LineStrip     PrimitiveType = gl.LINE_STRIP
	Points        PrimitiveType = gl.POINTS
)

// CreateMeshInfo contains all necessary information to create a mesh
type CreateMeshInfo struct {
	Vertices      []float32
	Indices       []uint32
	Sizes         []int32
	PrimitiveType PrimitiveType
}

// CreateMesh creates a new mesh from the provided information
func CreateMesh(info CreateMeshInfo) (m Mesh, err error) {
	if len(info.Vertices) == 0 {
		return m, fmt.Errorf("no vertices provided")
	}
	if len(info.Sizes) == 0 {
		return m, fmt.Errorf("no attribute sizes provided")
	}

	vertexSize := int32(0)
	for _, v := range info.Sizes {
		if v <= 0 {
			return m, fmt.Errorf("invalid attribute size: %d", v)
		}
		vertexSize += v
	}

	m.verticesCount = int32(len(info.Vertices)) / vertexSize
	if m.verticesCount*vertexSize != int32(len(info.Vertices)) {
		return m, fmt.Errorf("vertex data size mismatch")
	}

	// Set primitive type or default to triangles
	if info.PrimitiveType == 0 {
		m.primitiveType = uint32(Triangles)
	} else {
		m.primitiveType = uint32(info.PrimitiveType)
	}

	gl.GenVertexArrays(1, &m.vao)
	gl.BindVertexArray(m.vao)

	gl.GenBuffers(1, &m.vbo)
	gl.GenBuffers(1, &m.ebo)

	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(info.Vertices)*4, gl.Ptr(info.Vertices), gl.STATIC_DRAW)

	if len(info.Indices) > 0 {
		m.drawingMethod = drawElements
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(info.Indices)*4, gl.Ptr(info.Indices), gl.STATIC_DRAW)
		m.verticesCount = int32(len(info.Indices))
	}

	if err := m.setAttributes(info.Sizes, vertexSize); err != nil {
		m.Delete()
		return m, err
	}

	return m, nil
}

func (m *Mesh) setAttributes(sizes []int32, vertexSize int32) error {
	offset := uintptr(0)
	for i := 0; i < len(sizes); i++ {
		gl.VertexAttribPointerWithOffset(uint32(i), sizes[i], gl.FLOAT, false, vertexSize*4, offset*4)
		if err := getGLError(); err != nil {
			return fmt.Errorf("error setting vertex attribute %d: %v", i, err)
		}
		gl.EnableVertexAttribArray(uint32(i))
		offset += uintptr(sizes[i])
	}
	return nil
}

// Draw renders the mesh using the current shader and settings
func (m *Mesh) Draw() {
	gl.BindVertexArray(m.vao)

	switch m.drawingMethod {
	case drawArrays:
		gl.DrawArrays(m.primitiveType, 0, m.verticesCount)
	case drawElements:
		gl.DrawElements(m.primitiveType, m.verticesCount, gl.UNSIGNED_INT, nil)
	}
}

// Delete frees all OpenGL resources associated with the mesh
func (m *Mesh) Delete() {
	if m.vao != 0 {
		gl.DeleteVertexArrays(1, &m.vao)
		m.vao = 0
	}
	if m.vbo != 0 {
		gl.DeleteBuffers(1, &m.vbo)
		m.vbo = 0
	}
	if m.ebo != 0 {
		gl.DeleteBuffers(1, &m.ebo)
		m.ebo = 0
	}
}

// SetPrimitiveType changes how the mesh vertices are interpreted during rendering
func (m *Mesh) SetPrimitiveType(primitiveType PrimitiveType) {
	m.primitiveType = uint32(primitiveType)
}

// LoadMeshOptions contains options for loading an OBJ file
type LoadMeshOptions struct {
	FlipUVs     bool
	CalcNormals bool
	Scale       float32
}

// LoadMesh loads a mesh from an OBJ file
func LoadMesh(filePath string, options *LoadMeshOptions) (m Mesh, err error) {
	if options == nil {
		options = &LoadMeshOptions{
			FlipUVs:     false,
			CalcNormals: true,
			Scale:       1.0,
		}
	}

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return m, fmt.Errorf("failed to read file: %v", err)
	}

	str := string(bytes)
	lines := strings.Split(str, "\n")

	v := []float32{}   // vertex positions
	vt := []float32{}  // texture coordinates
	vn := []float32{}  // vertex normals
	faces := [][]int{} // store face vertex indices for normal calculation

	cmi := CreateMeshInfo{
		Vertices:      []float32{},
		Indices:       []uint32{},
		Sizes:         []int32{3, 2, 3}, // position, uv, normal
		PrimitiveType: Triangles,
	}

	for lineNum, line := range lines {
		columns := strings.Fields(line)
		if len(columns) == 0 {
			continue
		}

		switch columns[0] {
		case "v":
			if len(columns) < 4 {
				return m, fmt.Errorf("line %d: invalid vertex format: %s", lineNum+1, line)
			}
			coords := make([]float32, 3)
			for i := 0; i < 3; i++ {
				val, err := strconv.ParseFloat(columns[i+1], 32)
				if err != nil {
					return m, fmt.Errorf("line %d: invalid vertex coordinate: %v", lineNum+1, err)
				}
				coords[i] = float32(val) * options.Scale
			}
			v = append(v, coords...)

		case "vt":
			if len(columns) < 3 {
				return m, fmt.Errorf("line %d: invalid texture coordinate format: %s", lineNum+1, line)
			}
			coords := make([]float32, 2)
			for i := 0; i < 2; i++ {
				val, err := strconv.ParseFloat(columns[i+1], 32)
				if err != nil {
					return m, fmt.Errorf("line %d: invalid texture coordinate: %v", lineNum+1, err)
				}
				coords[i] = float32(val)
			}
			if options.FlipUVs {
				coords[1] = 1.0 - coords[1]
			}
			vt = append(vt, coords...)

		case "vn":
			if len(columns) < 4 {
				return m, fmt.Errorf("line %d: invalid normal format: %s", lineNum+1, line)
			}
			coords := make([]float32, 3)
			for i := 0; i < 3; i++ {
				val, err := strconv.ParseFloat(columns[i+1], 32)
				if err != nil {
					return m, fmt.Errorf("line %d: invalid normal coordinate: %v", lineNum+1, err)
				}
				coords[i] = float32(val)
			}
			vn = append(vn, coords...)

		case "f":
			if len(columns) < 4 {
				return m, fmt.Errorf("line %d: invalid face format (less than 3 vertices): %s", lineNum+1, line)
			}

			faceIndices := []int{} // store vertex indices for this face
			for i := 2; i < len(columns)-1; i++ {
				f := []string{columns[1], columns[i], columns[i+1]}
				for _, vertex := range f {
					vi, vti, vni := parseFace(vertex)

					// Store vertex index for normal calculation
					if vi > 0 {
						faceIndices = append(faceIndices, int(vi-1))
					}

					// Position
					if vi == 0 {
						cmi.Vertices = append(cmi.Vertices, 0, 0, 0)
					} else if vi > 0 && vi*3 <= len(v) {
						cmi.Vertices = append(cmi.Vertices,
							v[((vi-1)*3)+0], v[((vi-1)*3)+1], v[((vi-1)*3)+2],
						)
					} else {
						return m, fmt.Errorf("line %d: invalid vertex index: %d", lineNum+1, vi)
					}

					// Texture coordinates
					if vti == 0 {
						cmi.Vertices = append(cmi.Vertices, 0, 0)
					} else if vti > 0 && vti*2 <= len(vt) {
						cmi.Vertices = append(cmi.Vertices,
							vt[((vti-1)*2)+0], vt[((vti-1)*2)+1],
						)
					} else {
						return m, fmt.Errorf("line %d: invalid texture coordinate index: %d", lineNum+1, vti)
					}

					// Normals
					if vni == 0 {
						if options.CalcNormals {
							// We'll calculate normals after processing all faces
							cmi.Vertices = append(cmi.Vertices, 0, 0, 0)
						} else {
							cmi.Vertices = append(cmi.Vertices, 0, 0, 0)
						}
					} else if vni > 0 && vni*3 <= len(vn) {
						cmi.Vertices = append(cmi.Vertices,
							vn[((vni-1)*3)+0], vn[((vni-1)*3)+1], vn[((vni-1)*3)+2],
						)
					} else {
						return m, fmt.Errorf("line %d: invalid normal index: %d", lineNum+1, vni)
					}
				}
			}
			if len(faceIndices) >= 3 {
				faces = append(faces, faceIndices)
			}
		}
	}

	// Calculate normals if needed
	if options.CalcNormals && len(vn) == 0 {
		// Create a slice to store accumulated normals for each vertex
		vertexNormals := make([]Vector3, len(v)/3)

		// Calculate normal for each face and accumulate
		for _, face := range faces {
			if len(face) < 3 {
				continue
			}

			// Get vertices of the face
			v1 := Vector3{v[face[0]*3], v[face[0]*3+1], v[face[0]*3+2]}
			v2 := Vector3{v[face[1]*3], v[face[1]*3+1], v[face[1]*3+2]}
			v3 := Vector3{v[face[2]*3], v[face[2]*3+1], v[face[2]*3+2]}

			// Calculate face normal using cross product
			edge1 := Vector3{v2.X - v1.X, v2.Y - v1.Y, v2.Z - v1.Z}
			edge2 := Vector3{v3.X - v1.X, v3.Y - v1.Y, v3.Z - v1.Z}
			normal := edge1.Cross(edge2).Normalize()

			// Accumulate normal for each vertex of the face
			for _, idx := range face {
				vertexNormals[idx] = vertexNormals[idx].Add(normal)
			}
		}

		// Normalize accumulated normals and update vertex data
		normalOffset := 5 // Position(3) + UV(2) = 5
		vertexSize := 8   // Position(3) + UV(2) + Normal(3) = 8

		for i, normal := range vertexNormals {
			normalized := normal.Normalize()
			for j := 0; j < len(faces); j++ {
				baseIndex := j * 3 * vertexSize
				for k := 0; k < 3; k++ {
					vertIndex := baseIndex + k*vertexSize
					if vertIndex+normalOffset+2 < len(cmi.Vertices) {
						// Check if this vertex matches the current normal index
						x := cmi.Vertices[vertIndex]
						y := cmi.Vertices[vertIndex+1]
						z := cmi.Vertices[vertIndex+2]
						if x == v[i*3] && y == v[i*3+1] && z == v[i*3+2] {
							cmi.Vertices[vertIndex+normalOffset] = normalized.X
							cmi.Vertices[vertIndex+normalOffset+1] = normalized.Y
							cmi.Vertices[vertIndex+normalOffset+2] = normalized.Z
						}
					}
				}
			}
		}
	}

	return CreateMesh(cmi)
}

// Vector3 represents a 3D vector
type Vector3 struct {
	X, Y, Z float32
}

// Cross computes the cross product of two vectors
func (v Vector3) Cross(other Vector3) Vector3 {
	return Vector3{
		X: v.Y*other.Z - v.Z*other.Y,
		Y: v.Z*other.X - v.X*other.Z,
		Z: v.X*other.Y - v.Y*other.X,
	}
}

// Normalize returns a normalized version of the vector
func (v Vector3) Normalize() Vector3 {
	length := float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
	if length == 0 {
		return Vector3{0, 1, 0}
	}
	return Vector3{v.X / length, v.Y / length, v.Z / length}
}

// Add adds two vectors
func (v Vector3) Add(other Vector3) Vector3 {
	return Vector3{v.X + other.X, v.Y + other.Y, v.Z + other.Z}
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return i
}

func parseFace(vertexStr string) (vi, vti, vni int) {

	vertexComponents := strings.Split(vertexStr, "/")

	switch len(vertexComponents) {
	case 1:
		vi = atoi(vertexComponents[0])
	case 2:
		vi = atoi(vertexComponents[0])
		vti = atoi(vertexComponents[1])
	case 3:
		vi = atoi(vertexComponents[0])
		vti = atoi(vertexComponents[1])
		vni = atoi(vertexComponents[2])
	default:
		bayaan.Error("can not parse vertex", bayaan.Fields{
			"vertexStr": vertexStr,
		})
	}

	return
}

// getGLError checks for OpenGL errors
func getGLError() error {
	if err := gl.GetError(); err != 0 {
		return fmt.Errorf("OpenGL error: %d", err)
	}
	return nil
}
