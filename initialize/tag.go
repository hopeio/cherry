package initialize

import (
	"fmt"
	"github.com/hopeio/cherry/utils/log"
	"github.com/hopeio/cherry/utils/reflect/structtag"
)

// example:
/*
type Dao struct {
	DB mysql.DB `init:"config:MysqlTest"`
}

type Config struct {
	Env string
}

*/

const (
	initTagName = "init"
	exprTagName = "expr"

	metaTagName = "meta"
)

type InitTagSettings struct {
	ConfigName   string `meta:"config"`
	DefaultValue string `meta:"default"`
}

func ParseInitTagSettings(str string) *InitTagSettings {
	if str == "" {
		return &InitTagSettings{}
	}
	var settings InitTagSettings
	ParseTagSetting(str, ";", &settings)
	return &settings
}

// ParseTagSetting default sep ;
func ParseTagSetting(str string, sep string, settings any) {
	err := structtag.ParseTagSettingIntoStruct(str, sep, settings, metaTagName)
	if err != nil {
		log.Fatal(err)
	}
}

func genEncodingTag(name string) string {
	return fmt.Sprintf(`json:"%s" toml:"%s" yaml:"%s"`, name, name, name)
}
