package noor

type Camera interface {
	Projection() *float32
	View() *float32
}

type DefaultCamera struct{}

var IMat = [16]float32{
	1, 0, 0, 0,
	0, 1, 0, 0,
	0, 0, 1, 0,
	0, 0, 0, 1,
}

func (c DefaultCamera) Projection() *float32 {
	return &IMat[0]
}

func (c DefaultCamera) View() *float32 {
	return &IMat[0]
}
