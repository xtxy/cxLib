package jps

import "github.com/xtxy/cxlib/geo"

type jpsMoveDiagNoObs struct {
	finder *Finder
}

func (jps *jpsMoveDiagNoObs) findNeighbors(pos geo.Vec2[int64]) []geo.Vec2[int64] {
	var neighbors []geo.Vec2[int64]
	parentPos, parentOk := jps.finder.cellMap.GetParent(pos)

	if parentOk {
		dx, dy := dir(pos, parentPos)
		if dx != 0 && dy != 0 {
			deltas := [12]int64{
				0, dy, 0, dy,
				dx, 0, dx, 0,
			}
			flags := [2]bool{}

			neighbors = jps.finder.findNeighbors(pos, deltas[:], flags[:], true)
			if flags[0] && flags[1] {
				neighbors = append(neighbors, geo.Vec2[int64]{X: pos.X + dx, Y: pos.Y + dy})
			}
		} else if dx != 0 {
			deltas := [12]int64{
				dx, 0, dx, 0,
				0, 1, 0, 1,
				0, -1, 0, -1,
			}
			flags := [3]bool{}

			neighbors = jps.finder.findNeighbors(pos, deltas[:], flags[:], true)

			if flags[0] {
				if flags[1] {
					neighbors = append(neighbors, geo.Vec2[int64]{X: pos.X + dx, Y: pos.Y + 1})
				}
				if flags[2] {
					neighbors = append(neighbors, geo.Vec2[int64]{X: pos.X + dx, Y: pos.Y - 1})
				}
			}
		} else {
			deltas := [12]int64{
				0, dy, 0, dy,
				1, 0, 1, 0,
				-1, 0, -1, 0,
			}
			flags := [3]bool{}

			neighbors = jps.finder.findNeighbors(pos, deltas[:], flags[:], true)

			if flags[0] {
				if flags[1] {
					neighbors = append(neighbors, geo.Vec2[int64]{X: pos.X + 1, Y: pos.Y + dy})
				}
				if flags[2] {
					neighbors = append(neighbors, geo.Vec2[int64]{X: pos.X - 1, Y: pos.Y + dy})
				}
			}
		}
	} else {
		neighbors = jps.finder.findDefaultNeighbors(pos, MOVE_DIAG_NO_OBS)
	}

	return neighbors
}

func (jps *jpsMoveDiagNoObs) jump(pos, parent geo.Vec2[int64]) (next geo.Vec2[int64], ok bool) {
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
		if _, ok = jps.jump(geo.Vec2[int64]{X: pos.X + dx, Y: pos.Y}, pos); ok {
			return
		}

		if _, ok = jps.jump(geo.Vec2[int64]{X: pos.X, Y: pos.Y + dy}, pos); ok {
			return
		}
	} else {
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
		}
	}

	if jps.finder.canWalk(geo.Vec2[int64]{X: pos.X + dx, Y: pos.Y}) &&
		jps.finder.canWalk(geo.Vec2[int64]{X: pos.X, Y: pos.Y + dy}) {
		return jps.jump(geo.Vec2[int64]{X: pos.X + dx, Y: pos.Y + dy}, pos)
	}

	return
}
