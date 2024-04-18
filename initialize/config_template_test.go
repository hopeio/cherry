package initialize

import (
	"github.com/hopeio/cherry/utils/encoding"
	"testing"
)

func TestGenConfigTemplate(t *testing.T) {
	type args struct {
		format encoding.Format
		config Config
		dao    Dao
	}
}
