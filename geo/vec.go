package geo

import (
	"math"
	"strconv"
)

type Vec2 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (v1 Vec2) Sub(v2 Vec2) Vec2 {
	return Vec2{X: v1.X - v2.X, Y: v1.Y - v2.Y}
}

func (v1 Vec2) Cross(v2 Vec2) float64 {
	return v1.X*v2.Y - v1.Y*v2.X
}

func (v1 Vec2) LenSqr() float64 {
	return v1.X*v1.X + v1.Y*v1.Y
}

func (v1 Vec2) Dot(v2 Vec2) float64 {
	return v1.X*v2.X + v1.Y*v2.Y
}

func (v1 Vec2) Amplitude() float64 {
	if v1.X == 0 && v1.Y == 0 {
		return 0
	}
	return math.Sqrt(v1.X*v1.X + v1.Y*v1.Y)
}

func (v1 Vec2) Angle(v2 Vec2) float64 {
	angle, clockwise := v1.AngleDir(v2)
	if !clockwise {
		angle = 360 - angle
	}

	return angle
}

func (v1 Vec2) AngleDir(v2 Vec2) (float64, bool) {
	a := v1.Amplitude()
	b := v2.Amplitude()
	if a == 0 || b == 0 {
		return 0, true
	}

	angle := math.Acos(v1.Dot(v2)/a/b) * 180 / math.Pi
	return angle, v1.Cross(v2) <= 0
}

func (v1 *Vec2) Rotate(angle float64) {
	radian := float64(angle) * math.Pi / 180
	sinRet, cosRet := math.Sincos(radian)
	x := v1.X*cosRet + v1.Y*sinRet
	y := -v1.X*sinRet + v1.Y*cosRet
	v1.X, v1.Y = x, y
}

func (v1 Vec2) IsZero() bool {
	return v1.X == 0 && v1.Y == 0
}

func (v1 Vec2) Len() float64 {
	return math.Hypot(v1.X, v1.Y)
}

func (v1 *Vec2) Normalize() {
	if v1.X == 0 && v1.Y == 0 {
		return
	}

	l := v1.Len()
	v1.X /= l
	v1.Y /= l
}

type Vec2Int struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (v1 Vec2Int) Sub(v2 Vec2Int) Vec2Int {
	return Vec2Int{X: v1.X - v2.X, Y: v1.Y - v2.Y}
}

func (v1 Vec2Int) Eq(v2 Vec2Int) bool {
	return v1.X == v2.X && v1.Y == v2.Y
}

func (v1 Vec2Int) DistanceSqr(v2 Vec2Int) int {
	dx := v1.X - v2.X
	dy := v1.Y - v2.Y
	return dx*dx + dy*dy
}

func (v1 Vec2Int) Distance(v2 Vec2Int) float64 {
	return math.Sqrt(float64(v1.DistanceSqr(v2)))
}

func (v1 Vec2Int) String() string {
	return strconv.Itoa(v1.X) + "_" + strconv.Itoa(v1.Y)
}

func (v1 Vec2Int) IsZero() bool {
	return v1.X == 0 && v1.Y == 0
}

type Vec3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type Vec3Int struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}
