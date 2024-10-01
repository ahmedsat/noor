package noor

type Scene struct {
	Objects []*Object
	Camera  *Camera
}

func NewScene(camera *Camera, objects ...*Object) *Scene {
	s := &Scene{
		Objects: objects,
		Camera:  camera,
	}
	return s
}

func (s *Scene) AddObject(obj *Object) {
	s.Objects = append(s.Objects, obj)
}

func (s *Scene) SetCamera(cam *Camera) {
	s.Camera = cam
}

func (s *Scene) UpdateCamera() {
	for i := range s.Objects {
		s.Objects[i].UpdateCamera(*s.Camera)
	}
}

// func (s *Scene) AddObject(o *Object) {
// 	o.UpdateCamera(*s.Camera)
// 	s.Objects = append(s.Objects, o)
// }

func (s *Scene) Draw() {
	for _, o := range s.Objects {
		o.Draw(*s.Camera)
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
