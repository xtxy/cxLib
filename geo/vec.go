package geo

import (
	"math"
)

type Number interface {
	int32 | int64 | float32 | float64
}

type Vec2[T Number] struct {
	X T
	Y T
}

func (v1 Vec2[T]) Add(v2 Vec2[T]) Vec2[T] {
	return Vec2[T]{X: v1.X + v2.X, Y: v1.Y + v2.Y}
}

func (v1 *Vec2[T]) AddFrom(v2 *Vec2[T]) {
	v1.X += v2.X
	v1.Y += v2.Y
}

func (v1 Vec2[T]) Sub(v2 Vec2[T]) Vec2[T] {
	return Vec2[T]{X: v1.X - v2.X, Y: v1.Y - v2.Y}
}

func (v1 *Vec2[T]) SubFrom(v2 *Vec2[T]) {
	v1.X -= v2.X
	v1.Y -= v2.Y
}

func (v1 Vec2[T]) Cross(v2 Vec2[T]) T {
	return v1.X*v2.Y - v1.Y*v2.X
}

func (v1 Vec2[T]) LenSqr() T {
	return v1.X*v1.X + v1.Y*v1.Y
}

func (v1 Vec2[T]) Dot(v2 Vec2[T]) T {
	return v1.X*v2.X + v1.Y*v2.Y
}

func (v1 Vec2[T]) Amplitude() float64 {
	if v1.X == 0 && v1.Y == 0 {
		return 0
	}
	return math.Sqrt(float64(v1.X*v1.X + v1.Y*v1.Y))
}

func (v1 Vec2[T]) Angle(v2 Vec2[T]) float64 {
	angle, clockwise := v1.AngleDir(v2)
	if !clockwise {
		angle = 360 - angle
	}

	return angle
}

func (v1 Vec2[T]) AngleDir(v2 Vec2[T]) (float64, bool) {
	a := v1.Amplitude()
	b := v2.Amplitude()
	if a == 0 || b == 0 {
		return 0, true
	}

	angle := math.Acos(float64(v1.Dot(v2))/a/b) * 180 / math.Pi
	return angle, v1.Cross(v2) <= 0
}

func (v1 *Vec2[T]) Rotate(angle float64) {
	radian := float64(angle) * math.Pi / 180
	sinRet, cosRet := math.Sincos(radian)
	x := float64(v1.X)*cosRet + float64(v1.Y)*sinRet
	y := -float64(v1.X)*sinRet + float64(v1.Y)*cosRet
	v1.X, v1.Y = T(x), T(y)
}

func (v1 Vec2[T]) IsZero() bool {
	return v1.X == 0 && v1.Y == 0
}

func (v1 Vec2[T]) Len() float64 {
	return math.Hypot(float64(v1.X), float64(v1.Y))
}

func (v1 *Vec2[T]) Normalize() {
	if v1.X == 0 && v1.Y == 0 {
		return
	}

	l := v1.Len()
	v1.X = T(float64(v1.X) / l)
	v1.Y = T(float64(v1.Y) / l)
}
