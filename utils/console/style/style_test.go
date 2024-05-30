package style

import (
	"fmt"
	"testing"
)

func TestCustom(t *testing.T) {
	fmt.Println(Custom(31, 39)("红色"))
	fmt.Println(Custom(47, 0)("红色"))
}
