package slices

import (
	"fmt"
	"testing"
)

func TestContains(t *testing.T) {
	val1 := []string{"a", "b", "c"}
	val2 := "a"
	val3 := "d"
	fmt.Println(Contains(val1, val2))
	fmt.Println(Contains(val1, val3))
}

func TestForEachByIdx(t *testing.T) {
	val1 := []string{"a", "b", "c"}
	ForEachIndex(val1, func(i int) {
		fmt.Println(val1[i])
	})
}
