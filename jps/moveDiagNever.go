package jps

import "cxlib/geo"

type jpsMoveDiagNever struct {
	finder *Finder
}

func (jps *jpsMoveDiagNever) findNeighbors(pos geo.Vec2[int64]) []geo.Vec2[int64] {
	var neighbors []geo.Vec2[int64]
	parentPos, parentOk := jps.finder.cellMap.GetParent(pos)

	if parentOk {
		dx, dy := dir(pos, parentPos)
		if dx != 0 {
			deltas := [12]int64{
				0, -1, 0, -1,
				0, 1, 0, 1,
				dx, 0, dx, 0,
			}
			neighbors = jps.finder.findNeighbors(pos, deltas[:], nil, true)
		} else {
			deltas := [12]int64{
				-1, 0, -1, 0,
				1, 0, 1, 0,
				0, dy, 0, dy,
			}
			neighbors = jps.finder.findNeighbors(pos, deltas[:], nil, true)
		}
	} else {
		neighbors = jps.finder.findDefaultNeighbors(pos, MOVE_DIAG_NEVER)
	}

	return neighbors
}

func (jps *jpsMoveDiagNever) jump(pos, parent geo.Vec2[int64]) (next geo.Vec2[int64], ok bool) {
	if !jps.finder.canWalk(pos) {
		return
	}

	next = pos

	if pos == jps.finder.endPos {
		ok = true
		return
	}

	dx, dy := dir(pos, parent)
	if dx != 0 {
		if jumpCanWalk(jps.finder, pos, [8]int64{
			0, -1, -dx, -1, 0, 1, -dx, 1,
		}) {
			ok = true
			return
		}
	} else {
		if jumpCanWalk(jps.finder, pos, [8]int64{
			-1, 0, -1, -dy, 1, 0, 1, -dy,
		}) {
			ok = true
			return
		}

		if _, ok = jps.jump(geo.Vec2[int64]{X: pos.X + 1, Y: pos.Y}, pos); ok {
			return
		}

		if _, ok = jps.jump(geo.Vec2[int64]{X: pos.X - 1, Y: pos.Y}, pos); ok {
			return
		}
	}

	return jps.jump(geo.Vec2[int64]{X: pos.X + dx, Y: pos.Y + dy}, pos)
}
