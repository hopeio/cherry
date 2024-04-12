package style

import (
	"fmt"
	"github.com/hopeio/cherry/utils/slices"
	"strconv"
)

var colorFormat = "\x1b[%dm%s\x1b[0m"
var styleFormat = "\x1b[%s;%dm%s\x1b[0m"

func colorize(colorCode Color, s string) string {
	return fmt.Sprintf(colorFormat, colorCode, s)
}

func Style(colorCode Color, s string, wordFormats ...Decoration) string {
	if len(wordFormats) == 0 {
		return colorize(colorCode, s)
	}
	return fmt.Sprintf(styleFormat, slices.Join(wordFormats, ";"), colorCode, s)
}

type Decoration int

func (d Decoration) String() string {
	return strconv.Itoa(int(d))
}

const (
	Normal Decoration = iota
	Bold
	Italic
	Underline
)

type Color int

func (c Color) String() string {
	return strconv.Itoa(int(c))
}

const (
	ColorWhite         Color = 0
	ColorRed                 = 31
	ColorGreen               = 32
	ColorYellow              = 33
	ColorBlue                = 34
	ColorPurple              = 35
	ColorLightGreen          = 36
	ColorGray                = 37
	ColorRedBackground       = 41
)

func Blue(s string) string {
	return colorize(ColorBlue, s)
}

func LightGreen(s string) string {
	return colorize(ColorLightGreen, s)
}

func Purple(s string) string {
	return colorize(ColorPurple, s)
}

func White(s string) string {
	return colorize(ColorWhite, s)
}

func Gray(s string) string {
	return colorize(ColorGray, s)
}

func Red(s string) string {
	return colorize(ColorRed, s)
}

func RedBackground(s string) string {
	return colorize(ColorRedBackground, s)
}

func Green(s string) string {
	return colorize(ColorGreen, s)
}

func Yellow(s string) string {
	return colorize(ColorYellow, s)
}
