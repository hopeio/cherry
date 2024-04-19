package initialize

import (
	"flag"
	reflecti "github.com/hopeio/cherry/utils/reflect/converter"
	"github.com/spf13/pflag"
	"net/http"
	"net/url"
	"os"
	"reflect"
)

const flagTagName = "flag"

// TODO: 优先级高于其他Config,覆盖环境变量及配置中心的配置
// example
/*type FlagConfig struct {
	// environment
	Env string `flag:"name:env;short:e;default:dev;usage:环境"`
	// 配置文件路径
	ConfUrl string `flag:"name:confdao;short:c;default:config.toml;usage:配置文件路径,默认./config.toml或./config/config.toml"`
}*/

type FlagTagSettings struct {
	Name    string `meta:"name"`
	Short   string `meta:"short"`
	Env     string `meta:"env" explain:"从环境变量读取"`
	Default string `meta:"default"`
	Usage   string `meta:"usage"`
}

func init() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	commandLine := newCommandLine()
	injectFlagConfig(commandLine, reflect.ValueOf(&globalConfig1.BasicConfig).Elem())
	parseFlag(commandLine)

	if globalConfig1.Proxy != "" {
		proxyURL, _ := url.Parse(globalConfig1.Proxy)
		http.DefaultClient.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	} else {
		http.DefaultClient.Transport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		}
	}
}

type anyValue reflect.Value

func (a anyValue) String() string {
	return reflecti.String(reflect.Value(a))
}

func (a anyValue) Type() string {
	return reflect.Value(a).Kind().String()
}

func (a anyValue) Set(v string) error {
	return reflecti.SetFieldByString(v, reflect.Value(a))
}

func injectFlagConfig(commandLine *pflag.FlagSet, fcValue reflect.Value) {
	if !fcValue.IsValid() || fcValue.IsZero() {
		return
	}
	fcTyp := fcValue.Type()

	for i := 0; i < fcTyp.NumField(); i++ {
		fieldType := fcTyp.Field(i)
		if !fieldType.IsExported() {
			continue
		}
		flagTag := fieldType.Tag.Get(flagTagName)
		fieldValue := fcValue.Field(i)
		kind := fieldValue.Kind()
		if kind == reflect.Pointer {
			if !fieldValue.IsValid() || fieldValue.IsNil() {
				fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
			}
			injectFlagConfig(commandLine, fieldValue.Elem())
		}
		if kind == reflect.Struct {
			injectFlagConfig(commandLine, fieldValue)
		}
		if flagTag != "" {
			var flagTagSettings FlagTagSettings
			ParseTagSetting(flagTag, ";", &flagTagSettings)
			// 从环境变量设置
			if flagTagSettings.Env != "" {
				if value, ok := os.LookupEnv(flagTagSettings.Env); ok {
					reflecti.SetFieldByString(value, fcValue.Field(i))
				}
			}
			// flag设置
			flag := commandLine.VarPF(anyValue(fieldValue), flagTagSettings.Name, flagTagSettings.Short, flagTagSettings.Usage)
			if kind == reflect.Bool {
				flag.NoOptDefVal = "true"
			}
		}
	}
}

func (gc *globalConfig) applyFlagConfig() {
	commandLine := newCommandLine()
	fcValue := reflect.ValueOf(&gc.BasicConfig).Elem()
	injectFlagConfig(commandLine, fcValue)
	fcValue = reflect.ValueOf(&gc.EnvConfig).Elem()
	injectFlagConfig(commandLine, fcValue)
	parseFlag(commandLine)
}

func newCommandLine() *pflag.FlagSet {
	commandLine := pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
	commandLine.ParseErrorsWhitelist.UnknownFlags = true
	return commandLine
}

func parseFlag(commandLine *pflag.FlagSet) {
	commandLine.Parse(os.Args[1:])
}
