package jps

type jpsMoveNone struct {
	finder *Finder
}

func (jps *jpsMoveNone) findNeighbors(key string) []string {
	return jps.finder.findDefaultNeighbors(key, MOVE_DIAG_ALWAYS)
}

func (jps *jpsMoveNone) jump(key, parent string) string {
	return key
}
