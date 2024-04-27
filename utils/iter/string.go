package iter

import "unicode/utf8"

type String string

func (a String) Iterator() Iterator[rune] {
	return &stringIter{string(a)}
}

func (a String) Count() int {
	return len(a)
}

type stringIter struct {
	str string
}

// String returns an Iterator yielding runes from the supplied string.
func StringOf(input string) Iterator[rune] {
	return &stringIter{
		str: input,
	}
}

func (it *stringIter) Next() (rune, bool) {
	if len(it.str) == 0 {
		return 0, false
	}
	value, width := utf8.DecodeRuneInString(it.str)
	it.str = it.str[width:]
	return value, true
}
