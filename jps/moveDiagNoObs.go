package jps

type jpsMoveDiagNoObs struct {
	finder    *Finder
	pointUtil PointUtil
}

func (jps *jpsMoveDiagNoObs) findNeighbors(key string) []string {
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
			if jps.finder.canWalkAt(cellX, cellY+dy) && jps.finder.canWalkAt(cellX+dx, cellY) {
				neighbors = append(neighbors, jps.pointUtil.Point2Key(cellX+dx, cellY+dy))
			}
		} else if dx != 0 {
			nextOk := jps.finder.canWalkAt(cellX+dx, cellY)
			topOk := jps.finder.canWalkAt(cellX, cellY+1)
			bottomOk := jps.finder.canWalkAt(cellX, cellY-1)

			if nextOk {
				neighbors = append(neighbors, jps.pointUtil.Point2Key(cellX+dx, cellY))
				if topOk {
					neighbors = append(neighbors, jps.pointUtil.Point2Key(cellX+dx, cellY+1))
				}
				if bottomOk {
					neighbors = append(neighbors, jps.pointUtil.Point2Key(cellX+dx, cellY-1))
				}
			}
			if topOk {
				neighbors = append(neighbors, jps.pointUtil.Point2Key(cellX, cellY+1))
			}
			if bottomOk {
				neighbors = append(neighbors, jps.pointUtil.Point2Key(cellX, cellY-1))
			}
		} else {
			nextOk := jps.finder.canWalkAt(cellX, cellY+dy)
			rightOk := jps.finder.canWalkAt(cellX+1, cellY)
			leftOk := jps.finder.canWalkAt(cellX-1, cellY)

			if nextOk {
				neighbors = append(neighbors, jps.pointUtil.Point2Key(cellX, cellY+dy))
				if rightOk {
					neighbors = append(neighbors, jps.pointUtil.Point2Key(cellX+1, cellY+dy))
				}
				if leftOk {
					neighbors = append(neighbors, jps.pointUtil.Point2Key(cellX-1, cellY+dy))
				}
			}
			if rightOk {
				neighbors = append(neighbors, jps.pointUtil.Point2Key(cellX+1, cellY))
			}
			if leftOk {
				neighbors = append(neighbors, jps.pointUtil.Point2Key(cellX-1, cellY))
			}
		}
	} else {
		neighbors = jps.finder.findDefaultNeighbors(key, MOVE_DIAG_NO_OBS)
	}

	return neighbors
}

func (jps *jpsMoveDiagNoObs) jump(key, parent string) string {
	if !jps.finder.canWalk(key) {
		return ""
	}

	if key == jps.finder.endKey {
		return key
	}

	cellX, cellY := jps.pointUtil.Key2Point(key)
	dx, dy := dir(key, parent, jps.pointUtil)
	if dx != 0 && dy != 0 {
		if jps.jump(jps.pointUtil.Point2Key(cellX+dx, cellY), key) != "" ||
			jps.jump(jps.pointUtil.Point2Key(cellX, cellY+dy), key) != "" {
			return key
		}
	} else {
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
		}
	}

	if jps.finder.canWalkAt(cellX+dx, cellY) && jps.finder.canWalkAt(cellX, cellY+dy) {
		return jps.jump(jps.pointUtil.Point2Key(cellX+dx, cellY+dy), key)
	}

	return ""
}
