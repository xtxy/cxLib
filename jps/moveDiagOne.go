package jps

import "github.com/xtxy/cxlib/geo"

type jpsMoveDiagOne struct {
	finder *Finder
}

func (jps *jpsMoveDiagOne) findNeighbors(pos geo.Vec2[int64]) []geo.Vec2[int64] {
	var neighbors []geo.Vec2[int64]
	parentPos, parentOk := jps.finder.cellMap.GetParent(pos)

	if parentOk {
		dx, dy := dir(pos, parentPos)
		if dx != 0 && dy != 0 {
			deltas := [8]int64{
				0, dy, 0, dy,
				dx, 0, dx, 0,
			}
			flags := [2]bool{}

			neighbors = jps.finder.findNeighbors(pos, deltas[:], flags[:], true)
			if flags[0] || flags[1] {
				neighbors = append(neighbors, geo.Vec2[int64]{X: pos.X + dx, Y: pos.Y + dy})
			}

			if !jps.finder.canWalk(geo.Vec2[int64]{X: pos.X - dx, Y: pos.Y}) && flags[0] {
				neighbors = append(neighbors, geo.Vec2[int64]{X: pos.X - dx, Y: pos.Y + dy})
			}

			if jps.finder.canWalk(geo.Vec2[int64]{X: pos.X, Y: pos.Y - dy}) && flags[1] {
				neighbors = append(neighbors, geo.Vec2[int64]{X: pos.X + dx, Y: pos.Y - dy})
			}
		} else if dx == 0 {
			deltas := [4]int64{
				0, dy, 0, dy,
			}
			flags := [1]bool{}

			neighbors = jps.finder.findNeighbors(pos, deltas[:], flags[:], true)
			if flags[0] {
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
			deltas := [4]int64{
				dx, 0, dx, 0,
			}
			flags := [1]bool{}

			neighbors = jps.finder.findNeighbors(pos, deltas[:], flags[:], true)
			if flags[0] {
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
		neighbors = jps.finder.findDefaultNeighbors(pos, MOVE_DIAG_MOST_ONE)
	}

	return neighbors
}

func (jps *jpsMoveDiagOne) jump(pos, parent geo.Vec2[int64]) (next geo.Vec2[int64], ok bool) {
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
		}

		if _, ok = jps.jump(geo.Vec2[int64]{X: pos.X + dx, Y: pos.Y}, pos); ok {
			return
		}

		if _, ok = jps.jump(geo.Vec2[int64]{X: pos.X, Y: pos.Y + dy}, pos); ok {
			return
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

	if jps.finder.canWalk(geo.Vec2[int64]{X: pos.X + dx, Y: pos.Y}) ||
		jps.finder.canWalk(geo.Vec2[int64]{X: pos.X, Y: pos.Y + dy}) {
		return jps.jump(geo.Vec2[int64]{X: pos.X + dx, Y: pos.Y + dy}, pos)
	}

	return
}
