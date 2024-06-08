package style

import (
	"bufio"
	"fmt"
	"math"
	"math/rand/v2"
	"strings"
)

const rbgcFormat = "\x1b[38;2;%d;%d;%dm%c"
const rbgcFormatWithReset = rbgcFormat + reset
const rbgcBgFormat = "\x1b[48;2;%d;%d;%dm%c"
const rbgcBgFormatWithReset = rbgcBgFormat + reset
const rbgsFormat = "\x1b[38;2;%d;%d;%dm%s"
const rbgsFormatWithReset = rbgsFormat + reset
const rbgsBgFormat = "\x1b[48;2;%d;%d;%dm%s"
const rbgsBgFormatWithReset = rbgsBgFormat + reset

type colorRGB struct {
	r, g, b int16
}

func (c *colorRGB) Format(s string) string {
	return fmt.Sprintf(rbgsFormatWithReset, c.r, c.g, c.b, s)
}

func ColorRGBFormat(s string, r, g, b byte) string {
	return fmt.Sprintf(rbgsFormatWithReset, r, g, b, s)
}

func NewColorRGB(r, g, b byte) colorRGB {
	return colorRGB{r: int16(r), g: int16(g), b: int16(b)}
}

func ColorRGB(s string, r, g, b byte) string {
	return fmt.Sprintf(rbgsFormatWithReset, r, g, b, s)
}

func BgColorRGB(s string, r, g, b byte) string {
	return fmt.Sprintf(rbgsBgFormatWithReset, r, g, b, s)
}

func Gradient(text string, begin, end colorRGB) string {
	var colorText []string
	for i, r := range text {
		var ratio = float64(i) / float64(len(text)-1)
		var red = byte(math.Round(float64(begin.r) + float64(end.r-begin.r)*ratio))
		var green = byte(math.Round(float64(begin.g) + float64(end.g-begin.g)*ratio))
		var blue = byte(math.Round(float64(begin.b) + float64(end.b-begin.b)*ratio))
		colorText = append(colorText, fmt.Sprintf(rbgcFormat, red, green, blue, r))
	}
	colorText = append(colorText, reset)
	return strings.Join(colorText, "")
}

func GradientRandom(text string) string {
	var begin = colorRGB{r: int16(rand.N(byte(255))), g: int16(rand.N(byte(255))), b: int16(rand.N(byte(255)))}
	var end = colorRGB{r: int16(rand.N(byte(255))), g: int16(rand.N(byte(255))), b: int16(rand.N(byte(255)))}
	return Gradient(text, begin, end)
}

func GradientMultiLine(text string, begin, end colorRGB) string {
	scanner := bufio.NewScanner(strings.NewReader(text))
	var colorText []string
	for scanner.Scan() {
		colorText = append(colorText, Gradient(scanner.Text(), begin, end))
	}
	return strings.Join(colorText, "\n")
}

func GradientMultiLineRandom(text string) string {
	scanner := bufio.NewScanner(strings.NewReader(text))
	var colorText []string
	for scanner.Scan() {
		var begin = colorRGB{r: int16(rand.N(byte(255))), g: int16(rand.N(byte(255))), b: int16(rand.N(byte(255)))}
		var end = colorRGB{r: int16(rand.N(byte(255))), g: int16(rand.N(byte(255))), b: int16(rand.N(byte(255)))}
		colorText = append(colorText, Gradient(scanner.Text(), begin, end))
	}
	return strings.Join(colorText, "\n")
}

var (
	RainbowRGB       = [...]colorRGB{RainbowRedRGB, RainbowOrangeRGB, RainbowYellowRGB, RainbowGreenRGB, RainbowCyanRGB, RainbowBlueRGB, RainbowPurpleRGB}
	RainbowRedRGB    = NewColorRGB(255, 0, 0)
	RainbowOrangeRGB = NewColorRGB(255, 165, 0)
	RainbowYellowRGB = NewColorRGB(255, 255, 0)
	RainbowGreenRGB  = NewColorRGB(0, 255, 0)
	RainbowCyanRGB   = NewColorRGB(0, 255, 255)
	RainbowBlueRGB   = NewColorRGB(0, 0, 255)
	RainbowPurpleRGB = NewColorRGB(128, 0, 128)
)

func Rainbow(text string) string {
	var colorText []string
	var n int
	for _, r := range text {
		color := RainbowRGB[n]
		colorText = append(colorText, fmt.Sprintf(rbgcFormat, color.r, color.g, color.b, r))
		n++
		if n == 7 {
			n = 0
		}
	}
	colorText = append(colorText, reset)
	return strings.Join(colorText, "")
}

func RainbowMultiLine(text string) string {
	scanner := bufio.NewScanner(strings.NewReader(text))
	var colorText []string
	var n int
	for scanner.Scan() {
		color := RainbowRGB[n]
		colorText = append(colorText, color.Format(scanner.Text()))
		n++
		if n == 7 {
			n = 0
		}
	}
	colorText = append(colorText, reset)
	return strings.Join(colorText, "\n")
}

// TODO
func rainbowGradient(text string) {

}

func rgbToAnsi256(r, g, b byte) byte {
	// We use the extended greyscale palette here, with the exception of
	// black and white. normal palette only has 4 greyscale shades.
	if r>>4 == g>>4 && g>>4 == b>>4 {
		if r < 8 {
			return 16
		}
		if r > 248 {
			return 231
		}
		return byte(math.Round(float64(((r-8)/247)*24)) + 232)
	}
	var ansi = 16 +
		36*math.Round(float64((r/255)*5)) +
		6*math.Round(float64((g/255)*5)) +
		math.Round(float64((b/255)*5))
	return byte(ansi)
}
