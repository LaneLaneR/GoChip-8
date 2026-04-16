package iosystem

// type Input interface{}

type Output interface {
	Clear()
	Draw()
	DrawPixel(x int, y int, value bool) bool
}

type IoChip8 struct {
	// Input  Input
	Output
}

func NewIoChip8(output Output) *IoChip8 {
	return &IoChip8{
		Output: output,
	}
}
