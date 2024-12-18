package jps

import (
	"encoding/binary"
	"strconv"
	"strings"
)

type StrPoint struct {
}

func (sp StrPoint) Point2Key(x, y int) string {
	return strconv.Itoa(x) + "_" + strconv.Itoa(y)
}

func (sp StrPoint) Key2Point(key string) (int, int) {
	arr := strings.Split(key, "_")
	x, _ := strconv.Atoi(arr[0])
	y, _ := strconv.Atoi(arr[1])

	return x, y
}

type U32Point struct {
}

func (u32p U32Point) Point2Key(x, y int) string {
	slice := make([]byte, 8)
	binary.LittleEndian.PutUint32(slice[:4], uint32(x))
	binary.LittleEndian.PutUint32(slice[4:], uint32(y))

	return string(slice)
}

func (u32p U32Point) Key2Point(key string) (int, int) {
	slice := []byte(key)
	x := binary.LittleEndian.Uint32(slice[:4])
	y := binary.LittleEndian.Uint32(slice[4:])

	return int(x), int(y)
}
