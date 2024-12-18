package jps

import (
	"encoding/binary"
	"strconv"
	"strings"
)

type StrPoint struct {
}

func (sp StrPoint) Point2Key(x, y int64) string {
	return strconv.Itoa(int(x)) + "_" + strconv.Itoa(int(y))
}

func (sp StrPoint) Key2Point(key string) (int64, int64) {
	arr := strings.Split(key, "_")
	x, _ := strconv.Atoi(arr[0])
	y, _ := strconv.Atoi(arr[1])

	return int64(x), int64(y)
}

type U32Point struct {
}

func (u32p U32Point) Point2Key(x, y int) string {
	slice := make([]byte, 8)
	binary.LittleEndian.PutUint32(slice[:4], uint32(x))
	binary.LittleEndian.PutUint32(slice[4:], uint32(y))

	return string(slice)
}

func (u32p U32Point) Key2Point(key string) (int64, int64) {
	slice := []byte(key)
	x := binary.LittleEndian.Uint32(slice[:4])
	y := binary.LittleEndian.Uint32(slice[4:])

	return int64(x), int64(y)
}
