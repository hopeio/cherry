package style

import (
	"strconv"
)

var cursorFormat = "\x1b[%sm%s\x1b[0m"

type direction string

const (
	DirectionUp    direction = "A"
	DirectionDown  direction = "B"
	DirectionRight direction = "C"
	DirectionLeft  direction = "D"
)

type operation string

const (
	OperationSave          operation = "s"
	OperationRestore       operation = "u"
	OperationClear         operation = "2J"
	OperationClearLineTail operation = "K"
	OperationHide          operation = "?25l"
	OperationDisplay       operation = "?25h"
)

func Move(n int, direction direction) string {
	return strconv.Itoa(n) + string(direction)
}

func SetPosition(x, y int) string {
	return strconv.Itoa(x) + ";" + strconv.Itoa(y) + "H"
}
