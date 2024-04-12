package server

import "github.com/hopeio/cherry/server"

// 全局变量,只一个实例,只提供config
type Config server.Config

func (c *Config) InitBeforeInject() {
	(*server.Config)(c).InitBeforeInject()
}
func (c *Config) InitAfterInject() {
	(*server.Config)(c).InitAfterInject()
}

// TODO: 是否会随着配置而更新
func (c *Config) Update() bool {
	return false
}

func (c *Config) Origin() *server.Config {
	return (*server.Config)(c)
}
