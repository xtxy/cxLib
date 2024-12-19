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

type MapCell interface {
	SetParent(geo.Vec2[int64])
	GetParent() (geo.Vec2[int64], bool)
	SetG(uint32)
	GetG() uint32
	SetState(uint8)
	GetState() uint8
}

type CellMap interface {
	Reset()
	CalcG(srcCell, dstCell MapCell, srcPos, dstPos geo.Vec2[int64]) uint32
	GetCell(geo.Vec2[int64]) MapCell
}

type jpsMove interface {
	findNeighbors(geo.Vec2[int64]) []geo.Vec2[int64]
	jump(pos geo.Vec2[int64], parent geo.Vec2[int64]) (geo.Vec2[int64], bool)
}

type Finder struct {
	cellMap CellMap
	endPos  geo.Vec2[int64]
	move    jpsMove

	nearest bool
	blocks  map[string]struct{}
}

func NewFinder(cellMap CellMap, move int) *Finder {
	finder := new(Finder)
	finder.cellMap = cellMap

	switch move {
	case MOVE_DIAG_ALWAYS:
		moveInstance := new(jpsMoveDiag)
		moveInstance.finder = finder
		finder.move = moveInstance

	case MOVE_DIAG_MOST_ONE:
		moveInstance := new(jpsMoveDiagOne)
		moveInstance.finder = finder
		finder.move = moveInstance

	case MOVE_DIAG_NO_OBS:
		moveInstance := new(jpsMoveDiagNoObs)
		moveInstance.finder = finder
		finder.move = moveInstance

	case MOVE_DIAG_NEVER:
		moveInstance := new(jpsMoveDiagNever)
		moveInstance.finder = finder
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

func (finder *Finder) Find(start, end geo.Vec2[int64], options ...FindOption) []geo.Vec2[int64] {
	if cell := finder.cellMap.GetCell(start); cell == nil {
		logs.Warning("start.point.in.block:", start)
	}

	finder.nearest = false
	finder.blocks = nil

	for _, v := range options {
		v(finder)
	}

	finder.cellMap.Reset()

	if endCell := finder.cellMap.GetCell(end); (endCell == nil || endCell.GetState() == CELL_STATE_BLOCK) && !finder.nearest {
		logs.Error("end.point.in.block:", end)
	}

	finder.endPos = end
	found := false
	foundNearest := false
	nearestPos := geo.Vec2[int64]{}
	openList := []geo.Vec2[int64]{start}
	var nearestDistance int64 = 0

	for len(openList) > 0 && !found {
		pos := openList[0]
		openList = openList[1:]

		finder.cellMap.GetCell(pos).SetState(CELL_STATE_CLOSE)
		if pos == finder.endPos {
			found = true
			break
		}

		if finder.nearest {
			distanceSqr := pos.Sub(finder.endPos).LenSqr()
			if nearestDistance == 0 || distanceSqr < nearestDistance {
				nearestDistance = distanceSqr
				foundNearest = true
				nearestPos = pos
			}
		}

		openList = finder.identifySuccessors(openList, pos)
	}

	if !found {
		if finder.nearest && foundNearest {
			end = nearestPos
		} else {
			return nil
		}
	}

	list := make([]geo.Vec2[int64], 0)
	for ; end != start; end, _ = finder.cellMap.GetCell(end).GetParent() {
		list = append(list, end)
	}

	return list
}

func (finder *Finder) identifySuccessors(openList []geo.Vec2[int64], pos geo.Vec2[int64]) []geo.Vec2[int64] {
	src := finder.cellMap.GetCell(pos)
	srcG := src.GetG()
	neighbors := finder.move.findNeighbors(pos)
	for _, v := range neighbors {
		jumpPos, ok := finder.move.jump(v, pos)
		if !ok {
			continue
		}

		next := finder.cellMap.GetCell(jumpPos)
		if next.GetState() == CELL_STATE_CLOSE {
			continue
		}

		g := finder.cellMap.CalcG(src, next, pos, jumpPos)
		if next.GetState() != CELL_STATE_OPEN {
			next.SetState(CELL_STATE_OPEN)

			next.SetG(srcG + g)
			next.SetParent(pos)

			openList = append(openList, jumpPos)
		} else if (srcG + g) < next.GetG() {
			next.SetG(srcG + g)
			next.SetParent(pos)
		}
	}

	return openList
}

func (finder *Finder) canWalk(pos geo.Vec2[int64]) bool {
	if cell := finder.cellMap.GetCell(pos); cell != nil && cell.GetState() != CELL_STATE_BLOCK {
		return true
	}

	return false
}

func (finder *Finder) findDefaultNeighbors(pos geo.Vec2[int64], moveType int) []geo.Vec2[int64] {
	sFlags := [4]bool{}
	dFlags := [4]bool{}

	// up, right, down, left
	deltas := [16]int64{
		0, -1, 0, -1,
		1, 0, 1, 0,
		0, 1, 0, 1,
		-1, 0, -1, 0,
	}
	neighbors := finder.findNeighbors(pos, deltas[:], sFlags[:], true)

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
	allDeltas := [8]int64{
		-1, -1, 1, -1, 1, 1, -1, 1,
	}
	dDeltas := make([]int64, 0)
	for k, v := range dFlags {
		if !v {
			continue
		}

		dDeltas = append(dDeltas, allDeltas[k*2], allDeltas[k*2+1], allDeltas[k*2], allDeltas[k*2+1])
	}

	dNeighbors := finder.findNeighbors(pos, dDeltas, nil, true)
	if len(dNeighbors) > 0 {
		neighbors = append(neighbors, dNeighbors...)
	}

	return neighbors
}

func (finder *Finder) findNeighbors(pos geo.Vec2[int64], deltas []int64, flags []bool, canWalk bool) []geo.Vec2[int64] {
	nPos := geo.Vec2[int64]{}
	neighbors := make([]geo.Vec2[int64], 0)

	for i := 0; i < len(deltas); i += 4 {
		nPos.X = pos.X + deltas[i]
		nPos.Y = pos.Y + deltas[i+1]
		if finder.canWalk(nPos) == canWalk {
			nPos.X = pos.X + deltas[i+2]
			nPos.Y = pos.Y + deltas[i+3]

			neighbors = append(neighbors, nPos)
			if i < len(flags) {
				flags[i] = true
			}
		}
	}

	return neighbors
}
