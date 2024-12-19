package jps

import (
	"cxlib/geo"
)

func dir(pos, parentPos geo.Vec2[int64]) (int64, int64) {
	dx := clamp(pos.X - parentPos.X)
	dy := clamp(pos.Y - parentPos.Y)

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

func jumpCanWalk(finder *Finder, pos geo.Vec2[int64], deltas [8]int64) bool {
	if (finder.canWalk(geo.Vec2[int64]{X: pos.X + deltas[0], Y: pos.Y + deltas[1]}) &&
		!finder.canWalk(geo.Vec2[int64]{X: pos.X + deltas[2], Y: pos.Y + deltas[3]})) ||
		(finder.canWalk(geo.Vec2[int64]{X: pos.X + deltas[4], Y: pos.Y + deltas[5]}) &&
			!finder.canWalk(geo.Vec2[int64]{X: pos.X + deltas[6], Y: pos.Y + deltas[7]})) {
		return true
	}
	return false
}
