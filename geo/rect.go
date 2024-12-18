package geo

import (
	"math/rand"
	"time"
)

type Rect struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

func (rect *Rect) Center() Vec2 {
	return Vec2{
		X: rect.X + rect.Width/2,
		Y: rect.Y + rect.Height/2,
	}
}

func (rect *Rect) CenterInt() Vec2Int {
	return Vec2Int{
		X: int(rect.X + 0.5 + rect.Width/2),
		Y: int(rect.Y + 0.5 + rect.Height/2),
	}
}

func (rect *Rect) RandPos() Vec2 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return Vec2{
		X: rect.X + float64(1+r.Intn(int(rect.Width)-1)),
		Y: rect.Y + float64(1+r.Intn(int(rect.Height)-1)),
	}
}

func (rect *Rect) RandPosInt() Vec2Int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return Vec2Int{
		X: int(rect.X) + 1 + r.Intn(int(rect.Width)-1),
		Y: int(rect.Y) + 1 + r.Intn(int(rect.Height)-1),
	}
}

func (rect *Rect) Contain(pos Vec2) bool {
	return rect.X <= pos.X && rect.X+rect.Width >= pos.X &&
		rect.Y <= pos.Y && rect.Y+rect.Height >= pos.Y
}

type RectInt struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

func (rect *RectInt) Center() Vec2Int {
	return Vec2Int{
		X: int(rect.X + rect.Width/2),
		Y: int(rect.Y + rect.Height/2),
	}
}

func (rect *RectInt) Contain(pos Vec2Int) bool {
	return rect.X <= pos.X && rect.X+rect.Width >= pos.X &&
		rect.Y <= pos.Y && rect.Y+rect.Height >= pos.Y
}

func (rect *RectInt) RandPosInt() Vec2Int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return Vec2Int{
		X: rect.X + 1 + r.Intn(rect.Width-1),
		Y: rect.Y + 1 + r.Intn(rect.Height-1),
	}
}
