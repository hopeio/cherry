package path

import (
	"fmt"
	"testing"
)

func TestDir(t *testing.T) {
	dir, _ := Split("F:\\a\\video")
	t.Log(dir)
}

func TestClean(t *testing.T) {
	s := `......`
	fmt.Println(len(s))
	r := FileNameClean(s)
	fmt.Println(r)
	fmt.Println(len(r))
}

func TestRune(t *testing.T) {
	t.Log('，')
	t.Log('、' == '、')
}

func TestGetDirName(t *testing.T) {

	t.Log(GetDirName(""))
}
