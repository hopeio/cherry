package initconf

import "github.com/hopeio/cherry/initialize/conf_center"

type InitConfig struct {
	// 配置文件路径
	ConfUrl string `flag:"name:config;short:c;usage:配置文件路径,默认./config.toml或./config/config.toml;env:CONFIG"`
	BasicConfig
	EnvConfig
}

// BasicConfig
// zh: 基本配置，包含模块名
type BasicConfig struct {
	// 模块名
	Module string `flag:"name:mod;short:m;default:cherry-app;usage:模块名;env:MODULE"`
	// environment
	Env string `flag:"name:env;short:e;default:dev;usage:环境;env:ENV"`
}

type EnvConfig struct {
	Debug             bool   `flag:"name:debug;short:d;default:debug;usage:是否测试;env:DEBUG"`
	ConfigTemplateDir string `flag:"name:conf_tmpl_dir;short:t;usage:是否生成配置模板;env:CONFIG_TEMPLATE_DIR"`
	// 代理, socks5://localhost:1080
	Proxy    string `flag:"name:proxy;short:p;default:'socks5://localhost:1080';usage:代理;env:HTTP_PROXY" `
	NoInject []string
	// config字段顺序不能变,ConfigCenter 保持在最后
	ConfigCenter conf_center.Config
}
