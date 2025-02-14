package noor

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type Shader uint32

func CreateShaderProgram(vertexShaderSource, fragmentShaderSource string) Result[Shader] {

	sh := Shader(gl.CreateProgram())

	if err := compileShaderAndAttach(uint32(sh), vertexShaderSource, gl.VERTEX_SHADER); err != nil {
		return Err[Shader](errors.Join(err, errors.New("failed to compile vertex shader")))
	}

	if err := compileShaderAndAttach(uint32(sh), fragmentShaderSource, gl.FRAGMENT_SHADER); err != nil {
		return Err[Shader](errors.Join(err, errors.New("failed to compile fragment shader")))
	}

	gl.LinkProgram(uint32(sh))
	if err := checkProgramLinkStatus(uint32(sh)); err != nil {
		return Err[Shader](errors.Join(err, errors.New("failed to link shader program")))
	}

	return Ok[Shader](sh)
}

func CreateShaderProgramFromFiles(vertexShaderPath, fragmentShaderPath string) Result[Shader] {

	vertexShaderSourceResult := loadShaderSourceFromFile(vertexShaderPath)
	if vertexShaderSourceResult.IsErr() {
		fmt.Fprintf(os.Stderr, "Default vertex shader will be used\n")
		vertexShaderSourceResult = Ok(DefaultVertexShader)
	}

	fragmentShaderSourceResult := loadShaderSourceFromFile(fragmentShaderPath)
	if fragmentShaderSourceResult.IsErr() {
		fmt.Fprintf(os.Stderr, "Warning :Default fragment shader will be used\n")
		fragmentShaderSourceResult = Ok(DefaultFragmentShader)
	}

	Assert(
		vertexShaderSourceResult.IsOk() && fragmentShaderSourceResult.IsOk(),
		"Failed to load shaders from files",
	)

	vertexShaderSource := vertexShaderSourceResult.Ok
	fragmentShaderSource := fragmentShaderSourceResult.Ok

	return CreateShaderProgram(vertexShaderSource, fragmentShaderSource)
}

func compileShaderAndAttach(program uint32, source string, shaderType uint32) error {
	shader := gl.CreateShader(shaderType)
	defer gl.DeleteShader(shader)

	cSources, free := gl.Strs(source + "\x00")
	defer free()
	gl.ShaderSource(shader, 1, cSources, nil)
	gl.CompileShader(shader)

	if err := checkShaderCompileStatus(shader); err != nil {
		return fmt.Errorf("failed to compile shader (type: %d): %w", shaderType, err)
	}

	gl.AttachShader(program, shader)

	return nil
}

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

func loadShaderSourceFromFile(filePath string) Result[string] {
	source, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error :Failed to read shader file: %v\n", err)
		return Err[string](err)
	}
	return Ok(string(source))
}

func (sh *Shader) SetUniformFloat32(name string, value float32) {
	location := sh.GetUniformLocation(name)
	gl.Uniform1f(location, value)

}

func (sh *Shader) SetUniformBool(name string, value bool) {
	location := sh.GetUniformLocation(name)
	gl.Uniform1i(location, int32(boolToInt(value)))
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func (sh *Shader) SetUniformInt32(name string, value int32) {
	location := sh.GetUniformLocation(name)
	gl.Uniform1i(location, value)
}

func (sh *Shader) SetUniformMatrixFloat32(name string, value *float32) {
	location := sh.GetUniformLocation(name)
	gl.UniformMatrix4fv(location, 1, false, value)
}

func (sh *Shader) GetUniformLocation(name string) int32 {
	sh.Activate()
	nameCStr := gl.Str(name + "\x00")
	location := gl.GetUniformLocation(uint32(*sh), nameCStr)
	if location == -1 {
		fmt.Fprintf(os.Stderr, "Uniform location not found: %s\n", name)
	}
	return location
}

func (sh *Shader) Activate() {
	gl.UseProgram(uint32(*sh))
}

func (sh *Shader) Delete() {
	gl.DeleteProgram(uint32(*sh))
}
