package noor

import "github.com/ahmedsat/madar"

type Object struct {
	Name string
	Mesh *Mesh

	Position madar.Vector3
	Rotation madar.Vector3
	Scale    madar.Vector3

	Shader
	Textures []*Texture
}

func NewObject(name string, mesh *Mesh) *Object {

	defaultShader := CreateShaderProgram(
		DefaultVertexShader,
		DefaultFragmentShader,
	).UnwrapOrPanic()

	return &Object{
		Name:     name,
		Mesh:     mesh,
		Position: madar.Vector3{X: 0, Y: 0, Z: 0},
		Rotation: madar.Vector3{X: 0, Y: 0, Z: 0},
		Scale:    madar.Vector3{X: 1, Y: 1, Z: 1},

		Shader:   defaultShader,
		Textures: make([]*Texture, 0),
	}
}

func (o *Object) Render(camera Camera) {
	o.Shader.Activate()

	for i, tex := range o.Textures {
		tex.Activate(o.Shader, uint32(i), tex.Name)
	}

	o.Shader.SetUniformMatrixFloat32("uView", camera.View())
	o.Shader.SetUniformMatrixFloat32("uProjection", camera.Projection())
	o.Shader.SetUniformMatrixFloat32("uModel", o.ModelMatrix())
	o.Mesh.Draw()
}

func (o *Object) ModelMatrix() *float32 {
	var mat madar.Matrix = madar.TranslationMatrix(o.Position.X, o.Position.Y, o.Position.Z)
	mat = mat.Multiply(madar.RotationMatrix(o.Rotation.X, o.Rotation.Y, o.Rotation.Z))
	mat = mat.Multiply(madar.ScalingMatrix(o.Scale.X, o.Scale.Y, o.Scale.Z))
	return mat.Ptr()
}

func (o *Object) Translate(x, y, z float32) {
	o.Position.X += x
	o.Position.Y += y
	o.Position.Z += z
}

func (o *Object) Rotate(x, y, z float32) {
	o.Rotation.X += x
	o.Rotation.Y += y
	o.Rotation.Z += z
}

func (o *Object) ScaleBy(x, y, z float32) {
	o.Scale.X *= x
	o.Scale.Y *= y
	o.Scale.Z *= z
}

func (o *Object) AddTexture(tex *Texture) {
	o.Textures = append(o.Textures, tex)
}

func (o *Object) RemoveTexture(tex Texture) {
	for i, t := range o.Textures {
		if t.Name == tex.Name {
			o.Textures = append(o.Textures[:i], o.Textures[i+1:]...)
			break
		}
	}
}

func (o *Object) SetShader(shader Shader) {
	o.Shader.Delete()
	o.Shader = shader
}

func (o *Object) Delete() {
	o.Shader.Delete()
	o.Mesh.Delete()
}
