package jps

import "cxlib/geo"

type jpsMoveNone struct {
	finder *Finder
}

func (jps *jpsMoveNone) findNeighbors(pos geo.Vec2[int64]) []geo.Vec2[int64] {
	return jps.finder.findDefaultNeighbors(pos, MOVE_DIAG_ALWAYS)
}

func (jps *jpsMoveNone) jump(pos, parent geo.Vec2[int64]) (next geo.Vec2[int64], ok bool) {
	next = pos
	ok = true
	return
}
