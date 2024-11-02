package noor

import (
	"github.com/ahmedsat/madar"
)

var (
	defaultPosition = madar.Vector3{}
	defaultScale    = madar.Vector3{X: 1, Y: 1, Z: 1}
	defaultRotation = madar.Vector3{}
	defaultMesh     = Mesh{}    // todo: add default mesh
	defaultShader   = Shader(0) // todo: add default shader
	defaultTextures = []Texture{}
)

type ObjectCreateInfo struct {
	Position madar.Vector3
	Scale    madar.Vector3
	Rotation madar.Vector3
	Mesh     Mesh
	Shader   Shader
	Textures []Texture

	// TODO: Add more fields

}

type Object struct {
	position    madar.Vector3
	scale       madar.Vector3
	rotation    madar.Vector3
	mesh        Mesh
	shader      Shader
	textures    []Texture
	modelMatrix madar.Matrix4X4
	isDirty     bool
}

func CreateObject(info ObjectCreateInfo) *Object {

	o := &Object{
		position: If(info.Position.IsZero(), defaultPosition, info.Position),
		scale:    If(info.Scale.IsZero(), defaultScale, info.Scale),
		rotation: If(info.Rotation.IsZero(), defaultRotation, info.Rotation),
		mesh:     info.Mesh,
		shader:   info.Shader,
		textures: info.Textures,

		isDirty: true,
	}
	o.Update()

	return o
}

func (o *Object) Draw(c Camera) {
	o.shader.Activate()
	for i, t := range o.textures {
		t.Activate(o.shader, uint32(i))
	}

	o.shader.SetUniformMatrix4fv("uModel", o.modelMatrix)
	o.shader.SetUniformMatrix4fv("uView", c.GetViewMatrix())
	o.shader.SetUniformMatrix4fv("uProjection", c.GetProjectionMatrix())

	o.mesh.Draw()
}

func (o *Object) Update() {
	if o.isDirty {
		modelMatrix := madar.IdentityMatrix4X4()
		scaleMatrix := madar.ScaleMatrix4X4(o.scale)
		rotationMatrix := madar.RotationMatrix4X4(o.rotation)
		positionMatrix := madar.TranslationMatrix4X4(o.position)
		modelMatrix = modelMatrix.Multiply(scaleMatrix)
		modelMatrix = modelMatrix.Multiply(rotationMatrix)
		modelMatrix = modelMatrix.Multiply(positionMatrix)

		o.modelMatrix = modelMatrix

		o.isDirty = false
	}
}

func If[T any](condition bool, True, False T) T {
	if condition {
		return True
	}
	return False
}
