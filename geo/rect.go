package geo

import (
	"math/rand"
	"time"
)

type Rect[T Number] struct {
	X      T
	Y      T
	Width  T
	Height T
}

func (rect *Rect[T]) Center() Vec2[T] {
	return Vec2[T]{
		X: rect.X + rect.Width/2,
		Y: rect.Y + rect.Height/2,
	}
}

func (rect *Rect[T]) CenterInt() Vec2[int64] {
	return Vec2[int64]{
		X: int64(float64(rect.X) + 0.5 + float64(rect.Width)/2),
		Y: int64(float64(rect.Y) + 0.5 + float64(rect.Height)/2),
	}
}

func (rect *Rect[T]) RandPos() Vec2[T] {
	r := rand.New(rand.NewSource(time.Now().UnixMilli()))
	return Vec2[T]{
		X: rect.X + 1 + T(r.Intn(int(rect.Width)-1)),
		Y: rect.Y + T(1+r.Intn(int(rect.Height)-1)),
	}
}

func (rect *Rect[T]) RandPosInt() Vec2[int64] {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return Vec2[int64]{
		X: int64(rect.X) + 1 + int64(r.Intn(int(rect.Width)-1)),
		Y: int64(rect.Y) + 1 + int64(r.Intn(int(rect.Height)-1)),
	}
}

func (rect *Rect[T]) Contain(pos Vec2[T]) bool {
	return rect.X <= pos.X && rect.X+rect.Width >= pos.X &&
		rect.Y <= pos.Y && rect.Y+rect.Height >= pos.Y
}
