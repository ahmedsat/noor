package noor

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ahmedsat/bayaan"
	"github.com/ahmedsat/madar"
	"github.com/go-gl/gl/v4.5-core/gl"
)

type Shader uint32

// CreateShaderProgram creates a shader program from vertex and fragment shader sources.
func CreateShaderProgram(vertexShaderSource, fragmentShaderSource string) (sh Shader, err error) {
	bayaan.Info("Creating shader program")

	// Create the shader program
	sh = Shader(gl.CreateProgram())

	// Compile and attach vertex shader
	if err = compileShaderAndAttach(uint32(sh), vertexShaderSource, gl.VERTEX_SHADER); err != nil {
		return 0, errors.Join(err, errors.New("failed to compile vertex shader"))
	}
	bayaan.Trace("Vertex shader compiled and attached")

	// Compile and attach fragment shader
	if err = compileShaderAndAttach(uint32(sh), fragmentShaderSource, gl.FRAGMENT_SHADER); err != nil {
		return 0, errors.Join(err, errors.New("failed to compile fragment shader"))
	}
	bayaan.Trace("Fragment shader compiled and attached")

	// Link the shader program
	gl.LinkProgram(uint32(sh))
	if err = checkProgramLinkStatus(uint32(sh)); err != nil {
		return 0, errors.Join(err, errors.New("failed to link shader program"))
	}
	bayaan.Info("Shader program linked successfully")

	return sh, nil
}

// CreateShaderProgramFromFiles loads shader sources from files and creates a shader program.
func CreateShaderProgramFromFiles(vertexShaderPath, fragmentShaderPath string) (sh Shader, err error) {
	// Load vertex shader source from file
	vertexShaderSource, err := loadShaderSourceFromFile(vertexShaderPath)
	if err != nil {
		return 0, errors.Join(err, fmt.Errorf("failed to load vertex shader from file: %s", vertexShaderPath))
	}
	bayaan.Trace("Vertex shader source loaded from file: %s", vertexShaderPath)

	// Load fragment shader source from file
	fragmentShaderSource, err := loadShaderSourceFromFile(fragmentShaderPath)
	if err != nil {
		return 0, errors.Join(err, fmt.Errorf("failed to load fragment shader from file: %s", fragmentShaderPath))
	}
	bayaan.Trace("Fragment shader source loaded from file: %s", fragmentShaderPath)

	// Create shader program using loaded sources
	return CreateShaderProgram(vertexShaderSource, fragmentShaderSource)
}

// compileShaderAndAttach compiles a shader and attaches it to the program.
func compileShaderAndAttach(program uint32, source string, shaderType uint32) error {
	shader := gl.CreateShader(shaderType)
	defer gl.DeleteShader(shader) // Ensure shader cleanup after attachment

	// Set the source code and compile the shader
	cSources, free := gl.Strs(source + "\x00")
	defer free() // Free the C strings after use
	gl.ShaderSource(shader, 1, cSources, nil)
	gl.CompileShader(shader)

	// Check for compilation errors
	if err := checkShaderCompileStatus(shader); err != nil {
		return fmt.Errorf("failed to compile shader (type: %d): %w", shaderType, err)
	}

	gl.AttachShader(program, shader)
	bayaan.Trace("Shader compiled and attached to program (shader type: %d)", shaderType)

	return nil
}

// checkShaderCompileStatus checks if the shader was compiled successfully.
func checkShaderCompileStatus(shader uint32) error {
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return fmt.Errorf("shader compilation failed: %v", log)
	}
	return nil
}

// checkProgramLinkStatus checks if the program was linked successfully.
func checkProgramLinkStatus(program uint32) error {
	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return fmt.Errorf("program linking failed: %v", log)
	}
	return nil
}

// loadShaderSourceFromFile reads shader source from a file.
func loadShaderSourceFromFile(filePath string) (string, error) {
	source, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read shader file: %v", err)
	}
	return string(source), nil
}

// SetUniform1f sets a float uniform value.
func (sh *Shader) SetUniform1f(name string, value float32) {
	location := sh.getUniformLocation(name)
	gl.Uniform1f(location, value)
	bayaan.Trace("Set float uniform: %s = %f", name, value)
}

// SetUniform1i sets an integer uniform value.
func (sh *Shader) SetUniform1i(name string, value int32) {
	location := sh.getUniformLocation(name)
	gl.Uniform1i(location, value)
	bayaan.Trace("Set int uniform: %s = %d", name, value)
}

// SetUniform3f sets a vec3 uniform value (e.g., for color or position).
func (sh *Shader) SetUniform3f(name string, x, y, z float32) {
	location := sh.getUniformLocation(name)
	gl.Uniform3f(location, x, y, z)
	bayaan.Trace("Set vec3 uniform: %s = (%f, %f, %f)", name, x, y, z)
}

// SetUniform4f sets a vec4 uniform value.
func (sh *Shader) SetUniform4f(name string, x, y, z, w float32) {
	location := sh.getUniformLocation(name)
	gl.Uniform4f(location, x, y, z, w)
	bayaan.Trace("Set vec4 uniform: %s = (%f, %f, %f, %f)", name, x, y, z, w)
}

// SetUniformMatrix4fv sets a 4x4 matrix uniform value.
func (sh *Shader) SetUniformMatrix4fv(name string, matrix madar.Matrix4X4) {
	t := matrix.Transpose()
	location := sh.getUniformLocation(name)
	gl.UniformMatrix4fv(location, 1, false, &t[0])
	bayaan.Trace("Set mat4 uniform: %s", name)
}

// SetUniformMatrix3fv sets a 3x3 matrix uniform value.
func (sh *Shader) SetUniformMatrix3fv(name string, matrix madar.Matrix3X3) {
	t := matrix.Transpose()
	location := sh.getUniformLocation(name)
	gl.UniformMatrix3fv(location, 1, false, &t[0])
	bayaan.Trace("Set mat3 uniform: %s", name)
}

// getUniformLocation retrieves the location of a uniform variable in the shader program.
func (sh *Shader) getUniformLocation(name string) int32 {
	nameCStr := gl.Str(name + "\x00")
	location := gl.GetUniformLocation(uint32(*sh), nameCStr)
	if location == -1 {
		bayaan.Warn("Uniform %s does not exist in shader", name)
	}
	return location
}

// Activate activates the shader program for rendering.
func (sh *Shader) Activate() {
	gl.UseProgram(uint32(*sh))
}

// Delete deletes the shader program.
func (sh *Shader) Delete() {
	gl.DeleteProgram(uint32(*sh))
}
