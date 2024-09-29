package noor

import (
	"errors"
	"fmt"

	"github.com/ahmedsat/bayaan"
	"github.com/ahmedsat/madar"
	"github.com/go-gl/gl/v4.5-core/gl"
)

// Vertex structure: position, color, texCoord
type Vertex struct {
	Position madar.Vector3 // 3D position
	Color    madar.Vector3 // RGB color
	TexCoord madar.Vector2 // 2D texture coordinates
}

// NewVertex creates a new Vertex with given position, color, and texture coordinates
func NewVertex(position madar.Vector3, color madar.Vector3, texCoord madar.Vector2) (v Vertex) {
	bayaan.Trace("Creating new vertex: Position: %v, Color: %v, TexCoord: %v", position, color, texCoord)
	return Vertex{position, color, texCoord}
}

const VertexSize = 3*4 + 3*4 + 2*4 // position (3 floats) + color (3 floats) + texCoord (2 floats)

type Mesh struct {
	Vertices []Vertex
	Indices  []uint32
	VBO      uint32
	EBO      uint32
}

// NewMesh creates a new Mesh and uploads vertex and index data to GPU
func NewMesh(vertices []Vertex, indices []uint32) (*Mesh, error) {
	if !isInitialized {
		bayaan.Error("Renderer is not initialized, cannot create a new mesh")
		return nil, errUnInitialized
	}

	m := &Mesh{
		Vertices: vertices,
		Indices:  indices,
	}

	bayaan.Info("Creating a new mesh with %d vertices and %d indices", len(vertices), len(indices))

	vertexData := []float32{}
	for _, v := range vertices {
		vertexData = append(vertexData, v.Position[:]...) // Add position
		vertexData = append(vertexData, v.Color[:]...)    // Add color
		vertexData = append(vertexData, v.TexCoord[:]...) // Add texCoord
	}

	// Generate and bind the VBO
	gl.GenBuffers(1, &m.VBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertexData)*4, gl.Ptr(vertexData), gl.STATIC_DRAW)
	bayaan.Trace("VBO created and data loaded to GPU: %d bytes", len(vertexData)*4)

	// Generate and bind the EBO (if needed)
	if len(indices) > 0 {
		gl.GenBuffers(1, &m.EBO)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.EBO)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)
		bayaan.Trace("EBO created and data loaded to GPU: %d bytes", len(indices)*4)
	} else {
		bayaan.Info("No indices provided for this mesh, EBO not created")
	}

	// Check for OpenGL errors
	if err := checkOpenGLError("NewMesh"); err != nil {
		bayaan.Error("OpenGL error occurred while creating a new mesh: %s", err)
		return nil, err
	}

	bayaan.Info("Mesh created successfully")
	return m, nil
}

func setupVertexAttributes() {
	bayaan.Trace("Setting up vertex attribute pointers")
	// Set vertex attribute pointers
	gl.EnableVertexAttribArray(0) // Position
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, int32(VertexSize), 0)
	bayaan.Trace("Position attribute set (index 0)")

	gl.EnableVertexAttribArray(1) // Color
	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, int32(VertexSize), 3*4)
	bayaan.Trace("Color attribute set (index 1)")

	gl.EnableVertexAttribArray(2) // TexCoord
	gl.VertexAttribPointerWithOffset(2, 2, gl.FLOAT, false, int32(VertexSize), 6*4)
	bayaan.Trace("Texture coordinate attribute set (index 2)")
}

type Material struct {
	Shader
	Textures []Texture
}

type Object struct {
	Mesh        *Mesh
	Material    *Material
	VAO         uint32
	ModelMatrix madar.Matrix4 // Transformation for rendering
}

// NewObject creates a new Object, sets up VAO, and binds its Mesh and Material
func NewObject(mesh *Mesh, material *Material, matrix madar.Matrix4) (*Object, error) {
	if !isInitialized {
		bayaan.Error("Renderer is not initialized, cannot create a new object")
		return nil, errUnInitialized
	}

	o := &Object{
		Material:    material,
		ModelMatrix: matrix,
		Mesh:        mesh,
	}

	// Generate and bind the VAO
	gl.GenVertexArrays(1, &o.VAO)
	gl.BindVertexArray(o.VAO)
	bayaan.Trace("VAO created and bound: %d", o.VAO)

	// Setup vertex attributes
	setupVertexAttributes()

	// Setup Material (shader and textures)
	gl.UseProgram(o.Material.program)
	bayaan.Trace("Shader program %d bound", o.Material.program)

	o.Material.Shader.SetUniformMatrix4fv("model", o.ModelMatrix)
	bayaan.Trace("Model matrix set for the object")

	// Setup textures
	for i, texture := range o.Material.Textures {
		setTexture(o.Material.program, texture, i)
		bayaan.Trace("Texture %s (unit %d) bound to shader", texture.Name, i)
	}

	// Check for OpenGL errors
	if err := checkOpenGLError("NewObject"); err != nil {
		bayaan.Error("OpenGL error occurred while creating a new object: %s", err)
		return nil, err
	}

	bayaan.Info("Object created successfully")
	return o, nil
}

// Utility function to bind and set textures
func setTexture(program uint32, texture Texture, unit int) {
	gl.ActiveTexture(gl.TEXTURE0 + uint32(unit))
	gl.BindTexture(gl.TEXTURE_2D, texture.Handle)
	gl.Uniform1i(gl.GetUniformLocation(program, gl.Str(texture.Name+"\x00")), int32(unit))
}

// Draw renders the object with its mesh and material
func (o *Object) Draw(c Camera) error {
	if !isInitialized {
		bayaan.Error("Renderer is not initialized, cannot draw the object")
		return errUnInitialized
	}

	// Bind VAO and textures
	gl.BindVertexArray(o.VAO)
	for i, texture := range o.Material.Textures {
		setTexture(o.Material.program, texture, i)
	}

	o.Material.Shader.SetUniformMatrix4fv("view", c.ViewMatrix())
	o.Material.Shader.SetUniformMatrix4fv("projection", c.ProjectionMatrix())

	// Draw either using indices or vertices
	if len(o.Mesh.Indices) > 0 {
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, o.Mesh.EBO)
		gl.DrawElements(gl.TRIANGLES, int32(len(o.Mesh.Indices)), gl.UNSIGNED_INT, nil)

	} else {
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(o.Mesh.Vertices)))

	}

	return checkOpenGLError("Draw")
}

func (o *Object) UpdateModelMatrix(update func()) {
	update()
	gl.UseProgram(o.Material.program)
	o.Material.Shader.SetUniformMatrix4fv("model", o.ModelMatrix)
	bayaan.Trace("Model matrix set for the object")
}

func (o *Object) UpdateCamera(camera Camera) {
	gl.UseProgram(o.Material.program)

	view := camera.ViewMatrix()
	o.Material.Shader.SetUniformMatrix4fv("view", view)
	bayaan.Trace("View matrix set for the object")
	projection := camera.ProjectionMatrix()
	o.Material.Shader.SetUniformMatrix4fv("projection", projection)
	bayaan.Trace("projection matrix set for the object")

}

// checkOpenGLError checks for OpenGL errors after each call for debugging purposes
func checkOpenGLError(context string) error {
	if errCode := gl.GetError(); errCode != gl.NO_ERROR {
		errMsg := "OpenGL error in " + context + ": " + getGLErrorString(errCode)
		bayaan.Error(errMsg)
		return errors.New(errMsg)
	}
	return nil
}

// getGLErrorString returns a string representation of OpenGL error codes
func getGLErrorString(errCode uint32) string {
	switch errCode {
	case gl.INVALID_ENUM:
		return "INVALID_ENUM"
	case gl.INVALID_VALUE:
		return "INVALID_VALUE"
	case gl.INVALID_OPERATION:
		return "INVALID_OPERATION"
	case gl.STACK_OVERFLOW:
		return "STACK_OVERFLOW"
	case gl.STACK_UNDERFLOW:
		return "STACK_UNDERFLOW"
	case gl.OUT_OF_MEMORY:
		return "OUT_OF_MEMORY"
	default:
		return fmt.Sprintf("Unknown error (0x%x)", errCode)
	}
}
