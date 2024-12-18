package jps

import (
	"cxlib/geo"
)

func KeyPath2FullPath(keyPath []string, points PointUtil) []string {
	if len(keyPath) < 2 {
		return keyPath
	}

	fullPath := make([]string, 0)
	for k := 0; k < len(keyPath)-1; k++ {
		fullPath = append(fullPath, keyPath[k])
		points := FillPathPoint(keyPath[k], keyPath[k+1], points)
		if len(points) > 0 {
			fullPath = append(fullPath, points...)
		}
	}

	fullPath = append(fullPath, keyPath[len(keyPath)-1])
	return fullPath
}

func FillPathPoint(startKey, endKey string, points PointUtil) []string {
	startX, startY := points.Key2Point(startKey)
	endX, endY := points.Key2Point(endKey)

	deltaX := endX - startX
	deltaY := endY - startY

	var stepX int64 = 0
	if deltaX > 0 {
		stepX = 1
	} else if deltaX < 0 {
		stepX = -1
	}

	var stepY int64 = 0
	if deltaY > 0 {
		stepY = 1
	} else if deltaY < 0 {
		stepY = -1
	}

	path := make([]string, 0)
	for x, y := startX+stepX, startY+stepY; x != endX || y != endY; x, y = x+stepX, y+stepY {
		path = append(path, points.Point2Key(x, y))
	}

	return path
}

func RemovePathTail(path []string, end geo.Vec2[int64], lenSqr int64, points PointUtil) ([]string, bool) {
	ok := false
	index := len(path) - 1
	pos := geo.Vec2[int64]{}
	for ; index >= 0; index-- {
		x, y := points.Key2Point(path[index])
		pos.X = x
		pos.Y = y

		if distanceSqr(pos, end) > lenSqr {
			break
		}

		ok = true
	}

	index += 2
	if index < len(path) {
		return path[:index], ok
	}

	return path, ok
}

func distanceSqr(v1, v2 geo.Vec2[int64]) int64 {
	delta := v1.Sub(v2)
	return delta.LenSqr()
}

func dir(key, parent string, points PointUtil) (int64, int64) {
	cellX, cellY := points.Key2Point(key)
	parentX, parentY := points.Key2Point(parent)

	dx := clamp(cellX - parentX)
	dy := clamp(cellY - parentY)

	return dx, dy
}

func clamp(a int64) int64 {
	if a > 0 {
		return 1
	} else if a < 0 {
		return -1
	}

	return 0
}
