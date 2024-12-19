package jps

import "cxlib/geo"

type jpsMoveDiag struct {
	finder *Finder
}

func (jps *jpsMoveDiag) findNeighbors(pos geo.Vec2[int64]) []geo.Vec2[int64] {
	var neighbors []geo.Vec2[int64]
	parentPos, parentOk := jps.finder.cellMap.GetParent(pos)

	if parentOk {
		dx, dy := dir(pos, parentPos)
		if dx != 0 && dy != 0 {
			deltas := [20]int64{
				0, dy, 0, dy,
				dx, 0, dx, 0,
				dx, dy, dx, dy,
				-dx, 0, -dx, dy,
				0, -dy, dx, -dy,
			}

			neighbors = jps.finder.findNeighbors(pos, deltas[:], nil, true)

		} else if dx == 0 {
			nPos := geo.Vec2[int64]{X: pos.X, Y: pos.Y + dy}
			if jps.finder.canWalk(nPos) {
				neighbors = append(neighbors, nPos)

				deltas := [8]int64{
					1, 0, 1, dy,
					-1, 0, -1, dy,
				}
				newNeighbors := jps.finder.findNeighbors(pos, deltas[:], nil, false)
				if len(newNeighbors) > 0 {
					neighbors = append(neighbors, newNeighbors...)
				}
			}
		} else {
			nPos := geo.Vec2[int64]{X: pos.X + dx, Y: pos.Y}
			if jps.finder.canWalk(nPos) {
				neighbors = append(neighbors, nPos)

				deltas := [8]int64{
					0, 1, dx, 1,
					0, -1, dx, -1,
				}
				newNeighbors := jps.finder.findNeighbors(pos, deltas[:], nil, false)
				if len(newNeighbors) > 0 {
					neighbors = append(neighbors, newNeighbors...)
				}
			}
		}
	} else {
		neighbors = jps.finder.findDefaultNeighbors(pos, MOVE_DIAG_ALWAYS)
	}

	return neighbors
}

func (jps *jpsMoveDiag) jump(pos, parent geo.Vec2[int64]) (next geo.Vec2[int64], ok bool) {
	if !jps.finder.canWalk(pos) {
		return
	}

	next = pos

	if pos == jps.finder.endPos {
		ok = true
		return
	}

	dx, dy := dir(pos, parent)
	if dx != 0 && dy != 0 {
		if jumpCanWalk(jps.finder, pos, [8]int64{
			-dx, dy, -dx, 0, dx, -dy, 0, -dy,
		}) {
			ok = true
			return
		} else {
			if _, ok = jps.jump(geo.Vec2[int64]{X: pos.X + dx, Y: pos.Y}, pos); ok {
				return
			}

			if _, ok = jps.jump(geo.Vec2[int64]{X: pos.X, Y: pos.Y + dy}, pos); ok {
				return
			}
		}
	} else {
		if dx != 0 {
			if jumpCanWalk(jps.finder, pos, [8]int64{
				dx, 1, 0, 1, dx, -1, 0, -1,
			}) {
				ok = true
				return
			}
		} else {
			if jumpCanWalk(jps.finder, pos, [8]int64{
				1, dy, 1, 0, -1, dy, -1, 0,
			}) {
				ok = true
				return
			}
		}
	}

	return jps.jump(geo.Vec2[int64]{X: pos.X + dx, Y: pos.Y + dy}, pos)
}
