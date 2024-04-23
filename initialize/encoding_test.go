package initialize

import (
	"github.com/hopeio/cherry/utils/encoding"
	"github.com/hopeio/cherry/utils/encoding/common"
	"testing"
)

func TestYaml(t *testing.T) {
	basicConfig := &EnvConfig{}
	data, err := common.Marshal(encoding.Yaml, basicConfig)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(data))
}
