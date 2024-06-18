package local

import (
	"errors"
	"fmt"
	"github.com/hopeio/cherry/utils/configor/local"
	"os"
)

var ConfigCenter = &Local{}

type Local struct {
	Conf Config
}
type Config struct {
	local.Config
	ConfigPath string
}

func (cc *Local) Type() string {
	return "local"
}

func (cc *Local) Config() any {
	return &cc.Conf
}

// 本地配置
func (cc *Local) Handle(handle func([]byte)) error {
	if cc.Conf.ConfigPath == "" {
		return errors.New("empty local config path")
	}
	_, err := os.Stat(cc.Conf.ConfigPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("找不到配置: %v", err)
	}

	err = local.New(&cc.Conf.Config).Handle(handle, cc.Conf.ConfigPath)
	if err != nil {
		return fmt.Errorf("配置错误: %v", err)
	}

	return nil
}

func (cc *Local) Close() error {
	return nil
}
