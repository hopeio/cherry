package style

import (
	"bufio"
	"fmt"
	"math"
	"math/rand/v2"
	"strings"
)

const rbgcFormat = "\x1b[38;2;%d;%d;%dm%c\x1b[0m"
const rbgcBgFormat = "\x1b[48;2;%d;%d;%dm%c\x1b[0m"
const rbgsFormat = "\x1b[38;2;%d;%d;%dm%s\x1b[0m"
const rbgsBgFormat = "\x1b[48;2;%d;%d;%dm%s\x1b[0m"

type color struct {
	r, g, b int16
}

func (c *color) Format(text string) string {
	return fmt.Sprintf(rbgsFormat, c.r, c.g, c.b, text)
}

func RGBFormat(text string, r, g, b byte) string {
	return fmt.Sprintf(rbgsFormat, r, g, b, text)
}

func NewColor(r, g, b byte) color {
	return color{r: int16(r), g: int16(g), b: int16(b)}
}

func Gradient(text string, begin, end color) string {
	var colorText []string
	for i, r := range text {
		var ratio = float64(i) / float64(len(text)-1)
		var red = byte(math.Round(float64(begin.r) + float64(end.r-begin.r)*ratio))
		var green = byte(math.Round(float64(begin.g) + float64(end.g-begin.g)*ratio))
		var blue = byte(math.Round(float64(begin.b) + float64(end.b-begin.b)*ratio))
		colorText = append(colorText, fmt.Sprintf(rbgcFormat, red, green, blue, r))
	}
	return strings.Join(colorText, "")
}

func RandomGradient(text string) string {
	var begin = color{r: int16(rand.N(byte(255))), g: int16(rand.N(byte(255))), b: int16(rand.N(byte(255)))}
	var end = color{r: int16(rand.N(byte(255))), g: int16(rand.N(byte(255))), b: int16(rand.N(byte(255)))}
	return Gradient(text, begin, end)
}

func MultiLineGradient(text string, begin, end color) string {
	scanner := bufio.NewScanner(strings.NewReader(text))
	var colorText []string
	for scanner.Scan() {
		colorText = append(colorText, Gradient(scanner.Text(), begin, end))
	}
	return strings.Join(colorText, "\n")
}

func MultiLineRandomGradient(text string) string {
	scanner := bufio.NewScanner(strings.NewReader(text))
	var colorText []string
	for scanner.Scan() {
		var begin = color{r: int16(rand.N(byte(255))), g: int16(rand.N(byte(255))), b: int16(rand.N(byte(255)))}
		var end = color{r: int16(rand.N(byte(255))), g: int16(rand.N(byte(255))), b: int16(rand.N(byte(255)))}
		colorText = append(colorText, Gradient(scanner.Text(), begin, end))
	}
	return strings.Join(colorText, "\n")
}

var (
	RainbowColor  = [...]color{RainbowRed, RainbowOrange, RainbowYellow, RainbowGreen, RainbowCyan, RainbowBlue, RainbowPurple}
	RainbowRed    = NewColor(255, 0, 0)
	RainbowOrange = NewColor(255, 165, 0)
	RainbowYellow = NewColor(255, 255, 0)
	RainbowGreen  = NewColor(0, 255, 0)
	RainbowCyan   = NewColor(0, 255, 255)
	RainbowBlue   = NewColor(0, 0, 255)
	RainbowPurple = NewColor(128, 0, 128)
)

func Rainbow(text string) string {
	var colorText string
	var n int
	for _, r := range text {
		color := RainbowColor[n]
		colorText += fmt.Sprintf(rbgcFormat, color.r, color.g, color.b, r)
		n++
		if n == 7 {
			n = 0
		}
	}
	return colorText
}

func MultiLineRainbow(text string) string {
	scanner := bufio.NewScanner(strings.NewReader(text))
	var colorText []string
	var n int
	for scanner.Scan() {
		color := RainbowColor[n]
		colorText = append(colorText, color.Format(scanner.Text()))
		n++
		if n == 7 {
			n = 0
		}
	}
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
