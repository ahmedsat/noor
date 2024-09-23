package noor

type Scene struct {
	Objects []*Object
}

func NewScene() *Scene {
	return &Scene{}
}

func (s *Scene) AddObject(o *Object) {
	s.Objects = append(s.Objects, o)
}

func (s *Scene) Draw() {
	for _, o := range s.Objects {
		o.Draw()
	}
}

func (s *Scene) Clear() {
	s.Objects = []*Object{}
}

func (s *Scene) DeleteObject(o *Object) {
	for i, object := range s.Objects {
		if object == o {
			s.Objects = append(s.Objects[:i], s.Objects[i+1:]...)
		}
	}
}
