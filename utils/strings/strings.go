package strings

import (
	"bytes"
	"github.com/hopeio/cherry/utils/slices"
	"math/rand"
	"strings"
	"unicode"
	"unicode/utf8"
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
		return c ^ ' '
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
		return c ^ ' '
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

// And now lots of helper functions.

func SnakeToCamel[T ~string](s T) string {
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
			continue // Caller the underscore in s.
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

func CamelCaseSlice(elem []string) string { return SnakeToCamel(strings.Join(elem, "_")) }

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
func CutPart(s, sep string) string {
	sepLen := len(sep)
	sEndIndex := len(s) - 1
	for begin := 0; begin < sEndIndex; begin++ {
		end := begin + sepLen
		if begin <= sEndIndex && s[begin:end] == sep {
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
func CutPartContain(s, sep string) string {
	sepLen := len(sep)
	sEndIndex := len(s) - 1
	for begin := 0; begin < sEndIndex; begin++ {
		end := begin + sepLen
		if begin <= sEndIndex && s[begin:end] == sep {
			return s[:begin] + sep
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

// Split splits the camelcase word and returns a list of words. It also
// supports digits. Both lower camel case and upper camel case are supported.
// For more info please check: http://en.wikipedia.org/wiki/CamelCase
//
// Examples
//
//	"" =>                     [""]
//	"lowercase" =>            ["lowercase"]
//	"Class" =>                ["Class"]
//	"MyClass" =>              ["My", "Class"]
//	"MyC" =>                  ["My", "C"]
//	"HTML" =>                 ["HTML"]
//	"PDFLoader" =>            ["PDF", "Loader"]
//	"AString" =>              ["A", "String"]
//	"SimpleXMLParser" =>      ["Simple", "XML", "Parser"]
//	"vimRPCPlugin" =>         ["vim", "RPC", "Plugin"]
//	"GL11Version" =>          ["GL", "11", "Version"]
//	"99Bottles" =>            ["99", "Bottles"]
//	"May5" =>                 ["May", "5"]
//	"BFG9000" =>              ["BFG", "9000"]
//	"BöseÜberraschung" =>     ["Böse", "Überraschung"]
//	"Two  spaces" =>          ["Two", "  ", "spaces"]
//	"BadUTF8\xe2\xe2\xa1" =>  ["BadUTF8\xe2\xe2\xa1"]
//
// Splitting rules
//
//  1. If string is not valid UTF-8, return it without splitting as
//     single item array.
//  2. Assign all unicode characters into one of 4 sets: lower case
//     letters, upper case letters, numbers, and all other characters.
//  3. Iterate through characters of string, introducing splits
//     between adjacent characters that belong to different sets.
//  4. Iterate through array of split strings, and if a given string
//     is upper case:
//     if subsequent string is lower case:
//     move last character of upper case string to beginning of
//     lower case string
func SplitCamelCase(src string) (entries []string) {
	// don't split invalid utf8
	if !utf8.ValidString(src) {
		return []string{src}
	}
	entries = []string{}
	var runes [][]rune
	lastClass := 0
	class := 0
	// split into fields based on class of unicode character
	for _, r := range src {
		switch true {
		case unicode.IsLower(r):
			class = 1
		case unicode.IsUpper(r):
			class = 2
		case unicode.IsDigit(r):
			class = 3
		default:
			class = 4
		}
		if class == lastClass {
			runes[len(runes)-1] = append(runes[len(runes)-1], r)
		} else {
			runes = append(runes, []rune{r})
		}
		lastClass = class
	}
	// handle upper case -> lower case sequences, e.g.
	// "PDFL", "oader" -> "PDF", "Loader"
	for i := 0; i < len(runes)-1; i++ {
		if unicode.IsUpper(runes[i][0]) && unicode.IsLower(runes[i+1][0]) {
			runes[i+1] = append([]rune{runes[i][len(runes[i])-1]}, runes[i+1]...)
			runes[i] = runes[i][:len(runes[i])-1]
		}
	}
	// construct []string from results
	for _, s := range runes {
		if len(s) > 0 {
			entries = append(entries, string(s))
		}
	}
	return
}

// CamelCase camel-cases a protobuf name for use as a Go identifier.
//
// If there is an interior underscore followed by a lower case letter,
// drop the underscore and convert the letter to upper case.
func CamelCase(s string) string {
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	var b []byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '.' && i+1 < len(s) && IsASCIILower(s[i+1]):
			// Skip over '.' in ".{{lowercase}}".
		case c == '.':
			b = append(b, '_') // convert '.' to '_'
		case c == '_' && (i == 0 || s[i-1] == '.'):
			// Convert initial '_' to ensure we start with a capital letter.
			// Do the same for '_' after '.' to match historic behavior.
			b = append(b, 'X') // convert '_' to 'X'
		case c == '_' && i+1 < len(s) && IsASCIILower(s[i+1]):
			// Skip over '_' in "_{{lowercase}}".
		case IsASCIIDigit(c):
			b = append(b, c)
		default:
			// Assume we have a letter now - if not, it's a bogus identifier.
			// The next word is a sequence of characters that must start upper case.
			if IsASCIILower(c) {
				c ^= ' ' // convert lowercase to uppercase
			}
			b = append(b, c)

			// Accept lower case sequence that follows.
			for ; i+1 < len(s) && IsASCIILower(s[i+1]); i++ {
				b = append(b, s[i+1])
			}
		}
	}
	return string(b)
}
