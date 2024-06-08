package style

import (
	"fmt"
	"github.com/hopeio/cherry/utils/slices"
	"strconv"
)

const styleFormat = "\x1b[%sm%s"
const styleFormatWithReset = styleFormat + reset
const reset = "\x1b[0m"
const color256 = "\x1b[38;5;%dm%s"
const color256WithReset = color256 + reset
const bgColor256 = "\x1b[48;5;%dm%s"
const bgColor256WithReset = bgColor256 + reset

func colorize(colorCode Style, s string) string {
	return fmt.Sprintf(styleFormat+reset, colorCode.String(), s)
}

func StylesFormat(s string, styles ...Style) string {
	if len(styles) == 0 {
		return s
	}
	return fmt.Sprintf(styleFormatWithReset, slices.Join(styles, ";"), s)
}

type Style int

func (d Style) String() string {
	return strconv.Itoa(int(d))
}

const (
	DecorationNormal Style = iota
	DecorationBold
	_
	DecorationItalic
	DecorationUnderline
	DecorationFlashing
	_
	DecorationReverse
)
const (
	DecorationResetBold = 22 + iota
	_
	DecorationResetUnderline
	DecorationResetFlashing
	_
	DecorationResetReverse
)

const (
	ColorBlack Style = 30 + iota
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorGray
)

const (
	HighLightColorBlack Style = 90 + iota
	HighLightColorRed
	HighLightColorGreen
	HighLightColorYellow
	HighLightColorBlue
	HighLightColorMagenta
	HighLightColorCyan
	HighLightColorGray
)

const (
	BackGroundBlack Style = 40 + iota
	BackgroundRed
	BackgroundGreen
	BackgroundYellow
	BackgroundBlue
	BackgroundMagenta
	BackgroundCyan
	BackgroundGray
)

const (
	HighLightBackGroundBlack Style = 100 + iota
	HighLightBackGroundRed
	HighLightBackGroundGreen
	HighLightBackGroundYellow
	HighLightBackGroundBlue
	HighLightBackGroundMagenta
	HighLightBackGroundCyan
	HighLightBackGroundGray
)

func Blue(s string) string {
	return colorize(ColorBlue, s)
}

func Cyan(s string) string {
	return colorize(ColorCyan, s)
}

func Magenta(s string) string {
	return colorize(ColorMagenta, s)
}

func Gray(s string) string {
	return colorize(ColorGray, s)
}

func Red(s string) string {
	return colorize(ColorRed, s)
}

func BgRed(s string) string {
	return colorize(BackgroundRed, s)
}

func Green(s string) string {
	return colorize(ColorGreen, s)
}

func Yellow(s string) string {
	return colorize(ColorYellow, s)
}

func Custom(text string, begin, end any) string {
	return fmt.Sprintf("\x1b[%vm%s\x1b[%vm", begin, text, end)
}

func Color256(s string, c byte) string {
	return fmt.Sprintf(color256WithReset, c, s)
}

func BgColor256(s string, c byte) string {
	return fmt.Sprintf(bgColor256WithReset, c, s)
}
