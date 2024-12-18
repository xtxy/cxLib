package jps

import (
	"cxlib/geo"
	"cxlib/logs"
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

	var stepX int = 0
	if deltaX > 0 {
		stepX = 1
	} else if deltaX < 0 {
		stepX = -1
	}

	var stepY int = 0
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

func RemovePathTail(path []string, end geo.Vec2Int, distanceSqr int, points PointUtil) ([]string, bool) {
	ok := false
	index := len(path) - 1
	pos := geo.Vec2Int{}
	for ; index >= 0; index-- {
		x, y := points.Key2Point(path[index])
		pos.X = x
		pos.Y = y

		if pos.DistanceSqr(end) > distanceSqr {
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

func dir(key, parent string, points PointUtil) (int, int) {
	cellX, cellY := points.Key2Point(key)
	parentX, parentY := points.Key2Point(parent)

	dx := clamp(int(cellX - parentX))
	dy := clamp(int(cellY - parentY))

	if dx == 0 && dy == 0 {
		logs.Error("dx,dy,both are 0")
	}

	return dx, dy
}

func clamp(a int) int {
	if a > 0 {
		return 1
	} else if a < 0 {
		return -1
	}

	return 0
}
