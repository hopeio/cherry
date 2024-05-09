package strings

import (
	"bytes"
	"github.com/hopeio/cherry/utils/slices"
	"math/rand"
	"strings"
	"unicode"
)

func FormatLen(s string, length int) string {
	if len(s) < length {
		return s + strings.Repeat(" ", length-len(s))
	}
	return s[:length]
}

func QuoteToBytes(s string) []byte {
	b := make([]byte, 0, len(s)+2)
	b = append(b, '"')
	b = append(b, []byte(s)...)
	b = append(b, '"')
	return b
}

func CamelToSnake(name string) string {
	var ret bytes.Buffer

	multipleUpper := false
	var lastUpper rune
	var beforeUpper rune

	for _, c := range name {
		// Non-lowercase character after uppercase is considered to be uppercase too.
		isUpper := unicode.IsUpper(c) || (lastUpper != 0 && !unicode.IsLower(c))

		if lastUpper != 0 {
			// Output a delimiter if last character was either the first uppercase character
			// in a row, or the last one in a row (e.g. 'S' in "HTTPServer").
			// Do not output a delimiter at the beginning of the name.

			firstInRow := !multipleUpper
			lastInRow := !isUpper

			if ret.Len() > 0 && (firstInRow || lastInRow) && beforeUpper != '_' {
				ret.WriteByte('_')
			}
			ret.WriteRune(unicode.ToLower(lastUpper))
		}

		// Buffer uppercase char, do not output it yet as a delimiter may be required if the
		// next character is lowercase.
		if isUpper {
			multipleUpper = lastUpper != 0
			lastUpper = c
			continue
		}

		ret.WriteRune(c)
		lastUpper = 0
		beforeUpper = c
		multipleUpper = false
	}

	if lastUpper != 0 {
		ret.WriteRune(unicode.ToLower(lastUpper))
	}
	return string(ret.Bytes())
}

func ConvertToCamelCase(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'X')
		i++
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && isASCIILower(s[i+1]) {
			continue // Skip the underscore in s.
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if isASCIILower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		for i+1 < len(s) && isASCIILower(s[i+1]) {
			i++
			t = append(t, s[i])
		}
	}
	return string(t)
}

// Is c an ASCII lower-case letter?
func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// Is c an ASCII digit?
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// 仅首位小写（更符合接口的规范）
func LowerCaseFirst(t string) string {
	if t == "" {
		return ""
	}
	b := []byte(t)
	b[0] = LowerCase(b[0])
	return BytesToString(b)
	//return string(LowerCase(t[0])) + t[1:]
}

func LowerCase(c byte) byte {
	if 'A' <= c && c <= 'Z' {
		c += 'a' - 'A'
	}
	return c
}

func UpperCaseFirst(t string) string {
	if t == "" {
		return ""
	}
	b := []byte(t)
	b[0] = UpperCase(b[0])
	return BytesToString(b)
	//return string(UpperCase(t[0])) + t[1:]
}

func UpperCase(c byte) byte {
	if 'a' <= c && c <= 'z' {
		c -= 'a' - 'A'
	}
	return c
}

// TODO
func ReplaceRunes(s string, olds []rune, new rune) string {
	if len(olds) == 0 || (len(olds) == 1 && olds[0] == new) {
		return s // avoid allocation
	}

	panic("TODO")
}

func ReplaceRunesEmpty(s string, old ...rune) string {
	if len(old) == 0 {
		return s // avoid allocation
	}

	// Apply replacements to buffer.
	t := make([]byte, len(s))
	w := 0
	start := 0
	needCopy := false
	last := false
	for i, r := range s {
		if slices.In(r, old) {
			if needCopy {
				w += copy(t[w:], s[start:i])
				needCopy = false
			}
			last = true
			continue
		}
		needCopy = true
		if last {
			start = i
			last = false
		}
	}
	if needCopy {
		w += copy(t[w:], s[start:])
	}
	return string(t[0:w])
}

func Camel(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'X')
		i++
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && IsASCIILower(s[i+1]) {
			continue // Skip the underscore in s.
		}
		if IsASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if IsASCIILower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		for i+1 < len(s) && IsASCIILower(s[i+1]) {
			i++
			t = append(t, s[i])
		}
	}
	return string(t)
}

// And now lots of helper functions.

func CamelCase[T ~string](s T) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'X')
		i++
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && isASCIILower(s[i+1]) {
			continue // Caller the underscore in s.
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if isASCIILower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		for i+1 < len(s) && isASCIILower(s[i+1]) {
			i++
			t = append(t, s[i])
		}
	}
	return string(t)
}

func CamelCaseSlice(elem []string) string { return CamelCase(strings.Join(elem, "_")) }

type NumLetterSlice[T any] ['z' - '0' + 1]T

// 原来数组支持这样用
func (n *NumLetterSlice[T]) Set(b byte, v T) {
	n[b-'0'] = v
}

func ReplaceBytes(s string, olds []byte, new byte) string {
	if len(olds) == 0 || (len(olds) == 1 && olds[0] == new) {
		return s // avoid allocation
	}
	tmpl := make([]bool, 255)

	for _, b := range olds {
		tmpl[b] = true
	}

	// Apply replacements to buffer.
	t := make([]byte, len(s))
	copy(t, s)

	for i, r := range s {
		if r < 256 && tmpl[r] {
			t[i] = new
		}

	}

	return string(t)
}

// 将字符串中指定的ascii字符替换为空
func ReplaceBytesEmpty(s string, old ...byte) string {
	if len(old) == 0 {
		return s // avoid allocation
	}
	tmpl := make([]bool, 255)

	for _, b := range old {
		tmpl[b] = true
	}

	// Apply replacements to buffer.
	t := make([]byte, len(s))
	w := 0
	start := 0
	needCopy := false
	last := false
	for i, r := range s {
		if r < 256 && tmpl[r] {
			if needCopy {
				w += copy(t[w:], s[start:i])
				needCopy = false
			}
			last = true
			continue
		}
		needCopy = true
		if last {
			start = i
			last = false
		}
	}
	if needCopy {
		w += copy(t[w:], s[start:])
	}
	return string(t[0:w])
}

func Rand(length int) string {
	randId := make([]byte, length)
	for i := range randId {
		n := rand.Intn(62)
		if n > 9 {
			if n > 35 {
				randId[i] = byte(n - 36 + 'a')
			} else {
				randId[i] = byte(n - 10 + 'a')
			}

		} else {
			randId[i] = byte(n + '0')
		}
	}
	return BytesToString(randId)
}

/*
从字符串尾开始,返回指定字符截断后的字符串
ReverseCutPart("https://video.weibo.com/media/play?livephoto=https%3A%2F%2Flivephoto.us.sinaimg.cn%2F002OnXdGgx07YpcajtkH0f0f0100gv8Q0k01.mov", "%2F")
002OnXdGgx07YpcajtkH0f0f0100gv8Q0k01.mov
*/
func ReverseCutPart(s, key string) string {
	keyLen := len(key)
	sEndIndex := len(s) - 1
	if sEndIndex < keyLen {
		return s
	}
	for end := sEndIndex; end > 0; end-- {
		begin := end - keyLen
		if begin >= 0 && s[begin:end] == key {
			return s[end:]
		}
	}
	return s
}

/*
指定字符截断，返回阶段前的字符串
CutPart("https://wx1.sinaimg.cn/orj360/6ebedee6ly1h566bbzyc6j20n00cuabd.jpg", "wx1")
https://
*/
func CutPart(s, key string) string {
	keyLen := len(key)
	sEndIndex := len(s) - 1
	for begin := 0; begin < sEndIndex; begin++ {
		end := begin + keyLen
		if begin <= sEndIndex && s[begin:end] == key {
			return s[:begin]
		}
	}
	return s
}

/*
指定字符截断，返回阶段前加指定字符的字符串
CutPartContain("https://f.video.weibocdn.com/o0/F9Nmm1ZJlx080UxqxlJK010412004rJS0E010.mp4?label=mp4_hd&template=540x960.24.0&ori=0&ps=1CwnkDw1GXwCQx&Expires=1670569613&ssig=fAQcBh4HGt&KID=unistore,video", "mp4")
https://f.video.weibocdn.com/o0/F9Nmm1ZJlx080UxqxlJK010412004rJS0E010.mp4
*/
func CutPartContain(s, key string) string {
	keyLen := len(key)
	sEndIndex := len(s) - 1
	for begin := 0; begin < sEndIndex; begin++ {
		end := begin + keyLen
		if begin <= sEndIndex && s[begin:end] == key {
			return s[:begin] + key
		}
	}
	return s
}

func Cut(s, sep string) (string, string, bool) {
	if i := strings.Index(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, "", false
}

func ReverseCut(s, sep string) (string, string, bool) {
	if i := strings.LastIndex(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, "", false
}

// 寻找括号区间
// BracketsIntervals 在给定字符串中寻找由特定开始和结束符号包围的区间。
// 它会返回第一个找到的由tokenBegin和tokenEnd界定的字符串区间，
// 如果找到了则返回该区间和true，否则返回空字符串和false。
//
// 参数:
// s - 待搜索的字符串。
// tokenBegin - 搜索的开始符号。
// tokenEnd - 搜索的结束符号。
//
// 返回值:
// 第一个找到的由tokenBegin和tokenEnd界定的字符串区间，
// 如果找到了则返回该区间和true，否则返回空字符串和false。
func BracketsIntervals(s string, tokenBegin, tokenEnd rune) (string, bool) {
	var level int // 当前嵌套层级
	begin := -1   // 记录开始符号的索引
	for i, r := range s {
		if r == tokenBegin {
			if begin == -1 {
				begin = i // 首次遇到开始符号，记录其索引
			}
			level++ // 嵌套层级加一
		} else if r == tokenEnd {
			level-- // 遇到结束符号，嵌套层级减一
			if level == 0 {
				// 当层级归零时，表示找到了匹配的区间，返回该区间
				return s[begin : i+1], true
			}
		}
	}
	// 若遍历结束仍未找到匹配的区间，返回空字符串和false
	return "", false
}
