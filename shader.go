package noor

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ahmedsat/bayaan"
	"github.com/go-gl/gl/v4.5-core/gl"
)

type Shader struct {
	program uint32
}

func CreateShaderProgram(vertexShaderSource, fragmentShaderSource string) (sh Shader, err error) {

	bayaan.Info("creating shader program")

	sh.program = gl.CreateProgram()

	err = compileShaderAndAttach(sh.program, vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		err = errors.Join(err, errors.New("failed to create vertex shader: "+vertexShaderSource))
		return
	}
	bayaan.Trace("vertex shader has been compiled")

	err = compileShaderAndAttach(sh.program, fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		err = errors.Join(err, errors.New("failed to create fragment shader: "+fragmentShaderSource))
		return
	}
	bayaan.Trace("fragment shader has been compiled")

	gl.LinkProgram(sh.program)

	err = isShaderSuccess(sh.program, gl.LINK_STATUS)
	if err != nil {
		err = errors.Join(err, errors.New("failed to link shader program"))
		return
	}
	bayaan.Trace("shader program has been linked")

	return
}

func CreateShaderProgramFromFiles(vertexShaderPath, fragmentShaderPath string) (sh Shader, err error) {

	vertexShaderSource, err := os.ReadFile(vertexShaderPath)
	if err != nil {
		err = errors.Join(err, errors.New("failed to create vertex shader: "+vertexShaderPath))
		return
	}
	bayaan.Trace("vertex shader source: %s has been loaded to memory", vertexShaderPath)

	fragmentShaderSource, err := os.ReadFile(fragmentShaderPath)
	if err != nil {
		err = errors.Join(err, errors.New("failed to create fragment shader: "+fragmentShaderPath))
		return
	}
	bayaan.Trace("vertex shader source: %s has been loaded to memory", fragmentShaderPath)

	sh, err = CreateShaderProgram(string(vertexShaderSource), string(fragmentShaderSource))
	if err != nil {
		err = errors.Join(err, errors.New("failed to create shader program: "+vertexShaderPath+" and "+fragmentShaderPath))
		return
	}

	return
}

func compileShaderAndAttach(program uint32, source string, shaderType uint32) error {
	shader := gl.CreateShader(shaderType)

	cSources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, cSources, nil)
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

func isShaderSuccess(shader uint32, pName uint32) error {

	var success int32

	switch pName {
	case gl.COMPILE_STATUS:
		gl.GetShaderiv(shader, pName, &success)
	case gl.LINK_STATUS:
		gl.GetProgramiv(shader, pName, &success)
	default:
		return fmt.Errorf("unknown pName: %v", pName)
	}

	if success == 0 {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return fmt.Errorf("failed to compile %v: %v", log, shader)
	}

	return nil
}
