package path

import (
	"fmt"
	"testing"
)

func TestDir(t *testing.T) {
	t.Log(Split("F:\\a\\video"))
}

func TestClean(t *testing.T) {
	s := `......`
	fmt.Println(len(s))
	r := FileClean(s)
	fmt.Println(r)
	fmt.Println(len(r))
}

func TestRune(t *testing.T) {
	t.Log('，')
	t.Log('、' == '、')
}

func TestGetDirName(t *testing.T) {

	t.Log(CleanDir(""))
}
