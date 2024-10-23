package noor

import "github.com/ahmedsat/madar"

type Scene struct {
	Objects         []*Object
	Camera          *Camera
	LightPos        madar.Vector3
	AmbientColor    madar.Vector3
	AmbientStrength float32
	CameraPos       madar.Vector3
	LightColor      madar.Vector3
}

func NewScene() *Scene {
	return &Scene{}
}

func (s *Scene) AddObject(obj *Object) {
	s.Objects = append(s.Objects, obj)
}

func (s *Scene) Draw() {
	for _, obj := range s.Objects {
		obj.Shader.SetUniform3f("uAmbientColor", s.AmbientColor.X, s.AmbientColor.Y, s.AmbientColor.Z)
		obj.Shader.SetUniform1f("uAmbientStrength", s.AmbientStrength)

		obj.Shader.SetUniform3f("uLightPos", s.LightPos.X, s.LightPos.Y, s.LightPos.Z)
		obj.Shader.SetUniform3f("uLightColor", s.LightColor.X, s.LightColor.Y, s.LightColor.Z)

		obj.Shader.SetUniform3f("uCameraPos", s.CameraPos.X, s.CameraPos.Y, s.CameraPos.Z)

		obj.Draw(s.Camera)
	}
}
