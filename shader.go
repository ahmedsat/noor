package noor

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type Shader struct {
	program uint32
}

func CreateShaderProgram(vertexShaderSource, fragmentShaderSource string) (sh Shader, err error) {

	sh.program = gl.CreateProgram()

	err = compileShaderAndAttach(sh.program, vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return
	}

	err = compileShaderAndAttach(sh.program, fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return
	}

	gl.LinkProgram(sh.program)
	return
}

func compileShaderAndAttach(program uint32, source string, shaderType uint32) error {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	err := isShaderSuccess(shader, gl.COMPILE_STATUS)
	if err != nil {
		return err
	}

	gl.AttachShader(program, shader)
	gl.DeleteShader(shader)

	return nil
}

func isShaderSuccess(shader uint32, pname uint32) error {

	var success int32
	gl.GetShaderiv(shader, pname, &success)
	if success == 0 {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return fmt.Errorf("failed to compile %v: %v", log, shader)
	}

	return nil
}
