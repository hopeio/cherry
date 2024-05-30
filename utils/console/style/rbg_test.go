package style

import (
	"fmt"
	"testing"
)

func TestGradient(t *testing.T) {
	fmt.Println(Gradient("test:这是一个长长的渐变色的字符串", color{r: 0x3B, g: 0xD1, b: 0x91}, color{r: 0x2B, g: 0x4C, b: 0xEE}))
}

func TestRandomGradient(t *testing.T) {
	fmt.Println(RandomGradient("test:这是一个长长的随机的渐变色的字符串\n还添加了换行\n还添加了unicode符号:☔☕♈♉♊\nunicode后还有字符，会影响渐变吗"))
}

func TestMultiLineGradient(t *testing.T) {
	fmt.Println(MultiLineGradient("test:这是一个长长的随机的渐变色的字符串\n还添加了换行\n还添加了unicode符号:☔☕♈♉♊\nunicode后还有字符，会影响渐变吗", color{r: 0x3B, g: 0xD1, b: 0x91}, color{r: 0x2B, g: 0x4C, b: 0xEE}))
}

func TestMultiLineRandomGradient(t *testing.T) {
	fmt.Println(MultiLineRandomGradient("test:这是一个长长的随机的渐变色的字符串\n还添加了换行\n还添加了unicode符号:☔☕♈♉♊\nunicode后还有字符，会影响渐变吗"))
}

func TestRainbow(t *testing.T) {
	fmt.Println(Rainbow("test:这是一个长长的随机的渐变色的字符串\n还添加了换行\n还添加了unicode符号:☔☕♈♉♊\nunicode后还有字符，会影响渐变吗"))
}

func TestMultiLineRainbow(t *testing.T) {
	fmt.Println(MultiLineRainbow("test:这是一个长长的随机的渐变色的字符串\n还添加了换行\n还添加了unicode符号:☔☕♈♉♊\nunicode后还有字符，会影响渐变吗"))
}
