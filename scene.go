package noor

type Scene struct {
	Objects []*Object
	Camera  Camera
}

func NewScene() *Scene {
	return &Scene{
		Objects: []*Object{},
		Camera:  DefaultCamera{},
	}
}

func (s *Scene) AddObject(obj *Object) {
	s.Objects = append(s.Objects, obj)
}

func (s *Scene) RemoveObject(obj Object) {
	for i, o := range s.Objects {
		if o.Name == obj.Name {
			s.Objects = append(s.Objects[:i], s.Objects[i+1:]...)
			break
		}
	}
}

func (s *Scene) Render() {
	for _, obj := range s.Objects {
		obj.Render(s.Camera)
	}
}
