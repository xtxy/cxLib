package jps

import (
	"cxlib/geo"
	"cxlib/logs"
)

const (
	CELL_STATE_NORMAL = iota
	CELL_STATE_CLOSE
	CELL_STATE_OPEN
	CELL_STATE_BLOCK
)

const (
	MOVE_DIAG_NEVER = iota
	MOVE_DIAG_NO_OBS
	MOVE_DIAG_MOST_ONE
	MOVE_DIAG_ALWAYS
	MOVE_ASTAR
)

type PointUtil interface {
	Point2Key(x, y int64) string
	Key2Point(string) (int64, int64)
}

type Cell struct {
	parent string
	gState uint32 // highest 8 bits is State, other 24 bits is G
}

func (cell *Cell) init() {
	cell.gState = 0
	cell.parent = ""
}

func (cell *Cell) SetG(value uint32) {
	cell.gState = (cell.gState & 0xff000000) | (value & 0xffffff)
}

func (cell *Cell) GetG() uint32 {
	return cell.gState & 0xffffff
}

func (cell *Cell) SetState(value uint8) {
	cell.gState = (cell.gState & 0xffffff) | (uint32(value) << 24)
}

func (cell *Cell) GetState() uint8 {
	return uint8(cell.gState >> 24)
}

func NewCell() *Cell {
	cell := new(Cell)
	cell.init()

	return cell
}

type jpsMove interface {
	findNeighbors(key string) []string
	jump(key, parent string) string
}

type Finder struct {
	cellMap   map[string]*Cell
	endKey    string
	move      jpsMove
	pointUtil PointUtil

	nearest bool
	blocks  map[string]struct{}
}

func NewFinder(cellMap map[string]*Cell, pointUtil PointUtil, move int) *Finder {
	finder := new(Finder)
	finder.cellMap = cellMap
	finder.pointUtil = pointUtil

	switch move {
	case MOVE_DIAG_ALWAYS:
		moveInstance := new(jpsMoveDiag)
		moveInstance.finder = finder
		moveInstance.pointUtil = pointUtil
		finder.move = moveInstance

	case MOVE_DIAG_MOST_ONE:
		moveInstance := new(jpsMoveDiagOne)
		moveInstance.finder = finder
		moveInstance.pointUtil = pointUtil
		finder.move = moveInstance

	case MOVE_DIAG_NO_OBS:
		moveInstance := new(jpsMoveDiagNoObs)
		moveInstance.finder = finder
		moveInstance.pointUtil = pointUtil
		finder.move = moveInstance

	case MOVE_DIAG_NEVER:
		moveInstance := new(jpsMoveDiagNever)
		moveInstance.finder = finder
		moveInstance.pointUtil = pointUtil
		finder.move = moveInstance

	case MOVE_ASTAR:
		moveInstance := new(jpsMoveNone)
		moveInstance.finder = finder
		finder.move = moveInstance

	default:
		logs.Error("unkown.move.type:", move)
		return nil
	}

	return finder
}

type FindOption func(finder *Finder)

func FindOptNearest(nearest bool) FindOption {
	return func(finder *Finder) {
		finder.nearest = nearest
	}
}

func FindOptBlocks(blocks map[string]struct{}) FindOption {
	return func(finder *Finder) {
		finder.blocks = blocks
	}
}

func (finder *Finder) Find(startKey, endKey string, options ...FindOption) []string {
	if _, ok := finder.cellMap[startKey]; !ok {
		logs.Warning("start.point.in.block:", startKey, len(finder.cellMap))

		finder.cellMap[startKey] = NewCell()
	}

	finder.nearest = false
	finder.blocks = nil

	for _, v := range options {
		v(finder)
	}

	finder.reset()

	if endCell, ok := finder.cellMap[endKey]; (!ok || endCell.GetState() == CELL_STATE_BLOCK) && !finder.nearest {
		logs.Error("end.point.in.block:", endKey)
	}

	finder.endKey = endKey
	found := false
	openList := []string{startKey}
	var nearestDistance int64 = 0
	var nearestKey string = ""
	endPos := finder.key2Point(endKey)

	for len(openList) > 0 && !found {
		key := openList[0]
		openList = openList[1:]

		finder.cellMap[key].SetState(CELL_STATE_CLOSE)
		if key == endKey {
			found = true
			break
		}

		if finder.nearest {
			cellPos := finder.key2Point(key)
			distanceSqr := distanceSqr(cellPos, endPos)
			if nearestDistance == 0 || distanceSqr < nearestDistance {
				nearestDistance = distanceSqr
				nearestKey = key
			}
		}

		openList = finder.identifySuccessors(openList, key)
	}

	if !found {
		if finder.nearest && nearestKey != "" {
			endKey = nearestKey
		} else {
			return nil
		}
	}

	list := make([]string, 0)
	for ; endKey != startKey; endKey = finder.cellMap[endKey].parent {
		list = append(list, endKey)
	}

	return list
}

func (finder *Finder) reset() {
	blockLen := len(finder.blocks)

	for k, v := range finder.cellMap {
		v.init()

		if blockLen > 0 {
			if _, ok := finder.blocks[k]; ok {
				v.SetState(CELL_STATE_BLOCK)
			}
		}
	}
}

func (finder *Finder) identifySuccessors(openList []string, key string) []string {
	src := finder.cellMap[key]
	srcPos := finder.key2Point(key)
	srcG := src.GetG()
	neighbors := finder.move.findNeighbors(key)
	for _, v := range neighbors {
		jumpPoint := finder.move.jump(v, key)
		if jumpPoint == "" {
			continue
		}

		next := finder.cellMap[jumpPoint]
		if next.GetState() == CELL_STATE_CLOSE {
			continue
		}

		nextPos := finder.key2Point(jumpPoint)
		g := calcG(srcPos, nextPos)

		if next.GetState() != CELL_STATE_OPEN {
			next.SetState(CELL_STATE_OPEN)

			next.SetG(srcG + g)
			next.parent = key

			openList = append(openList, jumpPoint)
		} else if (srcG + g) < next.GetG() {
			next.SetG(srcG + g)
			next.parent = key
		}
	}

	return openList
}

func (finder *Finder) canWalk(key string) bool {
	if cell, ok := finder.cellMap[key]; ok && cell.GetState() != CELL_STATE_BLOCK {
		return true
	}

	return false
}

func (finder *Finder) canWalkAt(x, y int64) bool {
	key := finder.pointUtil.Point2Key(x, y)
	return finder.canWalk(key)
}

func (finder *Finder) findDefaultNeighbors(key string, moveType int) []string {
	sFlags := make([]bool, 4)
	dFlags := make([]bool, 4)

	cellPos := finder.key2Point(key)
	// up, right, down, left
	pos := []int64{cellPos.X, cellPos.Y - 1, cellPos.X + 1, cellPos.Y, cellPos.X, cellPos.Y + 1,
		cellPos.X - 1, cellPos.Y}
	neighbors := make([]string, 0)

	for k := range sFlags {
		nKey := finder.pointUtil.Point2Key(pos[k*2], pos[k*2+1])
		if finder.canWalk(nKey) {
			neighbors = append(neighbors, nKey)
			sFlags[k] = true
		} else {
			sFlags[k] = false
		}
	}

	switch moveType {
	case MOVE_DIAG_NEVER:
		return neighbors

	case MOVE_DIAG_NO_OBS:
		dFlags[0] = sFlags[3] && sFlags[0]
		dFlags[1] = sFlags[0] && sFlags[1]
		dFlags[2] = sFlags[1] && sFlags[2]
		dFlags[3] = sFlags[2] && sFlags[3]

	case MOVE_DIAG_MOST_ONE:
		dFlags[0] = sFlags[3] || sFlags[0]
		dFlags[1] = sFlags[0] || sFlags[1]
		dFlags[2] = sFlags[1] || sFlags[2]
		dFlags[3] = sFlags[2] || sFlags[3]

	default:
		dFlags[0], dFlags[1], dFlags[2], dFlags[3] = true, true, true, true
	}

	// leftup, rightup, rightdown, leftdown
	pos = []int64{cellPos.X - 1, cellPos.Y - 1, cellPos.X + 1, cellPos.Y - 1, cellPos.X + 1, cellPos.Y + 1,
		cellPos.X - 1, cellPos.Y + 1}

	for k, v := range dFlags {
		if !v {
			continue
		}

		nKey := finder.pointUtil.Point2Key(pos[k*2], pos[k*2+1])
		if finder.canWalk(nKey) {
			neighbors = append(neighbors, nKey)
		}
	}

	return neighbors
}

func (finder *Finder) key2Point(key string) geo.Vec2[int64] {
	x, y := finder.pointUtil.Key2Point(key)
	return geo.Vec2[int64]{X: int64(x), Y: int64(y)}
}

func calcG(cell, next geo.Vec2[int64]) uint32 {
	delta := cell.Sub(next)
	return uint32(delta.LenSqr())
}
