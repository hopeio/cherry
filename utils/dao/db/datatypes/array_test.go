package datatypes

import "testing"

func TestArray(t *testing.T) {
	data := `{{{1,2},{1,2}},{{1,2},{1,2}}}`
	arr := Array[Array[Array[int]]]{}
	err := arr.Scan(data)
	if err != nil {
		t.Error(err)
	}
	t.Log(arr)
}
