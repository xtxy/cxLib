package jps

type jpsMoveDiagOne struct {
	finder    *Finder
	pointUtil PointUtil
}

func (jps *jpsMoveDiagOne) findNeighbors(key string) []string {
	cell := jps.finder.cellMap[key]
	cellX, cellY := jps.pointUtil.Key2Point(key)
	parent := cell.parent

	neighbors := make([]string, 0)

	if parent != "" {
		dx, dy := dir(key, parent, jps.pointUtil)
		if dx != 0 && dy != 0 {
			if nKey := jps.pointUtil.Point2Key(cellX, cellY+dy); jps.finder.canWalk(nKey) {
				neighbors = append(neighbors, nKey)
			}
			if nKey := jps.pointUtil.Point2Key(cellX+dx, cellY); jps.finder.canWalk(nKey) {
				neighbors = append(neighbors, nKey)
			}
			if jps.finder.canWalkAt(cellX, cellY+dy) || jps.finder.canWalkAt(cellX+dx, cellY) {
				neighbors = append(neighbors, jps.pointUtil.Point2Key(cellX+dx, cellY+dy))
			}
			if !jps.finder.canWalkAt(cellX-dx, cellY) && jps.finder.canWalkAt(cellX, cellY+dy) {
				neighbors = append(neighbors, jps.pointUtil.Point2Key(cellX-dx, cellY+dy))
			}
			if !jps.finder.canWalkAt(cellX, cellY-dy) && jps.finder.canWalkAt(cellX+dx, cellY) {
				neighbors = append(neighbors, jps.pointUtil.Point2Key(cellX+dx, cellY-dy))
			}
		} else if dx == 0 {
			if nKey := jps.pointUtil.Point2Key(cellX, cellY+dy); jps.finder.canWalk(nKey) {
				neighbors = append(neighbors, nKey)
				if !jps.finder.canWalkAt(cellX+1, cellY) {
					neighbors = append(neighbors, jps.pointUtil.Point2Key(cellX+1, cellY+dy))
				}
				if !jps.finder.canWalkAt(cellX-1, cellY) {
					neighbors = append(neighbors, jps.pointUtil.Point2Key(cellX-1, cellY+dy))
				}
			}
		} else {
			if nKey := jps.pointUtil.Point2Key(cellX+dx, cellY); jps.finder.canWalk(nKey) {
				neighbors = append(neighbors, nKey)
				if !jps.finder.canWalkAt(cellX, cellY+1) {
					neighbors = append(neighbors, jps.pointUtil.Point2Key(cellX+dx, cellY+1))
				}
				if !jps.finder.canWalkAt(cellX, cellY-1) {
					neighbors = append(neighbors, jps.pointUtil.Point2Key(cellX+dx, cellY-1))
				}
			}
		}
	} else {
		neighbors = jps.finder.findDefaultNeighbors(key, MOVE_DIAG_MOST_ONE)
	}

	return neighbors
}

func (jps *jpsMoveDiagOne) jump(key, parent string) string {
	if !jps.finder.canWalk(key) {
		return ""
	}

	if key == jps.finder.endKey {
		return key
	}

	cellX, cellY := jps.pointUtil.Key2Point(key)
	dx, dy := dir(key, parent, jps.pointUtil)
	if dx != 0 && dy != 0 {
		if (jps.finder.canWalkAt(cellX-dx, cellY+dy) && !jps.finder.canWalkAt(cellX-dx, cellY)) ||
			(jps.finder.canWalkAt(cellX+dx, cellY-dy) && !jps.finder.canWalkAt(cellX, cellY-dy)) {
			return key
		} else if jps.jump(jps.pointUtil.Point2Key(cellX+dx, cellY), key) != "" ||
			jps.jump(jps.pointUtil.Point2Key(cellX, cellY+dy), key) != "" {
			return key
		}
	} else {
		if dx != 0 {
			if (jps.finder.canWalkAt(cellX+dx, cellY+1) && !jps.finder.canWalkAt(cellX, cellY+1)) ||
				(jps.finder.canWalkAt(cellX+dx, cellY-1) && !jps.finder.canWalkAt(cellX, cellY-1)) {
				return key
			}
		} else {
			if (jps.finder.canWalkAt(cellX+1, cellY+dy) && !jps.finder.canWalkAt(cellX+1, cellY)) ||
				(jps.finder.canWalkAt(cellX-1, cellY+dy) && !jps.finder.canWalkAt(cellX-1, cellY)) {
				return key
			}
		}
	}

	if jps.finder.canWalkAt(cellX+dx, cellY) || jps.finder.canWalkAt(cellX, cellY+dy) {
		return jps.jump(jps.pointUtil.Point2Key(cellX+dx, cellY+dy), key)
	}

	return ""
}
