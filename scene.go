package noor

type Scene struct {
	Objects []*Object
	Camera  Camera
}

func NewScene() *Scene {
	return &Scene{}
}

func (s *Scene) AddObject(obj *Object) {
	s.Objects = append(s.Objects, obj)
}

func (s *Scene) Draw() {
	for _, obj := range s.Objects {
		obj.Draw(&s.Camera)
	}
}
