package style

import (
	"fmt"
	"github.com/hopeio/cherry/utils/slices"
	"strconv"
)

var styleFormat = "\x1b[%sm%s\x1b[0m"

func colorize(colorCode Style, s string) string {
	return fmt.Sprintf(styleFormat, colorCode.String(), s)
}

func StyleFormat(s string, styles ...Style) string {
	if len(styles) == 0 {
		return s
	}
	return fmt.Sprintf(styleFormat, slices.Join(styles, ";"), s)
}

type Style int

func (d Style) String() string {
	return strconv.Itoa(int(d))
}

const (
	DecorationNormal Style = iota
	DecorationBold
	DecorationItalic
	DecorationUnderline
)

const (
	ColorBlack Style = 30 + iota
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorPurple
	ColorLightGreen
	ColorGray
)

const (
	BackGroundBlack Style = 40 + iota
	BackgroundRed
	BackgroundGreen
	BackgroundYellow
	BackgroundBlue
	BackgroundPurple
	BackgroundLightGreen
	BackgroundGray
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

func Gray(s string) string {
	return colorize(ColorGray, s)
}

func Red(s string) string {
	return colorize(ColorRed, s)
}

func RedBackground(s string) string {
	return colorize(BackgroundRed, s)
}

func Green(s string) string {
	return colorize(ColorGreen, s)
}

func Yellow(s string) string {
	return colorize(ColorYellow, s)
}

func Custom(begin, end any) func(text string) string {
	return func(text string) string {
		return fmt.Sprintf("\x1b[%vm%s\x1b[%vm", begin, text, end)
	}
}
