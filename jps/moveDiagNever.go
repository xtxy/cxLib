package jps

type jpsMoveDiagNever struct {
	finder    *Finder
	pointUtil PointUtil
}

func (jps *jpsMoveDiagNever) findNeighbors(key string) []string {
	cell := jps.finder.cellMap[key]
	cellX, cellY := jps.pointUtil.Key2Point(key)
	parent := cell.parent

	neighbors := make([]string, 0)

	if parent != "" {
		dx, dy := dir(key, parent, jps.pointUtil)
		if dx != 0 {
			if nKey := jps.pointUtil.Point2Key(cellX, cellY-1); jps.finder.canWalk(nKey) {
				neighbors = append(neighbors, nKey)
			}
			if nKey := jps.pointUtil.Point2Key(cellX, cellY+1); jps.finder.canWalk(nKey) {
				neighbors = append(neighbors, nKey)
			}
			if nKey := jps.pointUtil.Point2Key(cellX+dx, cellY); jps.finder.canWalk(nKey) {
				neighbors = append(neighbors, nKey)
			}
		} else {
			if nKey := jps.pointUtil.Point2Key(cellX-1, cellY); jps.finder.canWalk(nKey) {
				neighbors = append(neighbors, nKey)
			}
			if nKey := jps.pointUtil.Point2Key(cellX+1, cellY); jps.finder.canWalk(nKey) {
				neighbors = append(neighbors, nKey)
			}
			if nKey := jps.pointUtil.Point2Key(cellX, cellY+dy); jps.finder.canWalk(nKey) {
				neighbors = append(neighbors, nKey)
			}
		}
	} else {
		neighbors = jps.finder.findDefaultNeighbors(key, MOVE_DIAG_NEVER)
	}

	return neighbors
}

func (jps *jpsMoveDiagNever) jump(key, parent string) string {
	if !jps.finder.canWalk(key) {
		return ""
	}

	if key == jps.finder.endKey {
		return key
	}

	cellX, cellY := jps.pointUtil.Key2Point(key)
	dx, dy := dir(key, parent, jps.pointUtil)
	if dx != 0 {
		if (jps.finder.canWalkAt(cellX, cellY-1) && !jps.finder.canWalkAt(cellX-dx, cellY-1)) ||
			(jps.finder.canWalkAt(cellX, cellY+1) && !jps.finder.canWalkAt(cellX-dx, cellY+1)) {
			return key
		}
	} else {
		if (jps.finder.canWalkAt(cellX-1, cellY) && !jps.finder.canWalkAt(cellX-1, cellY-dy)) ||
			(jps.finder.canWalkAt(cellX+1, cellY) && !jps.finder.canWalkAt(cellX+1, cellY-dy)) {
			return key
		}

		if jps.jump(jps.pointUtil.Point2Key(cellX+1, cellY), key) != "" ||
			jps.jump(jps.pointUtil.Point2Key(cellX-1, cellY), key) != "" {
			return key
		}
	}

	return jps.jump(jps.pointUtil.Point2Key(cellX+dx, cellY+dy), key)
}
