package local

import (
	"errors"
	"fmt"
	"github.com/hopeio/cherry/utils/configor/local"
	"os"
)

var ConfigCenter = &Local{}

type Local struct {
	local.Config
	ConfigPath string
}

func (cc *Local) Type() string {
	return "local"
}

// 本地配置
func (cc *Local) HandleConfig(handle func([]byte)) error {
	if cc.ConfigPath == "" {
		return errors.New("empty local config path")
	}
	_, err := os.Stat(cc.ConfigPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("找不到配置: %v", err)
	}

	err = local.New(&cc.Config).Handle(handle, cc.ConfigPath)
	if err != nil {
		return fmt.Errorf("配置错误: %v", err)
	}

	return nil
}
