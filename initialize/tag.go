package initialize

import (
	"fmt"
	"github.com/hopeio/cherry/utils/log"
	"github.com/hopeio/cherry/utils/reflect/structtag"
	stringsi "github.com/hopeio/cherry/utils/strings"
	"reflect"
	"strings"
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

// get field name, return filed config name and skip flag
func getFieldConfigName(v reflect.StructField) (string, tagOptions, bool) {
	tag := v.Tag.Get("mapstructure")
	if tag == "" {
		return v.Name, "", true
	}
	if tag == "-" {
		return "", "", false
	}
	name, opts := parseTag(tag)
	if name == "" {
		return v.Name, opts, true
	}
	return stringsi.UpperCaseFirst(name), opts, true
}

// tagOptions is the string following a comma in a struct field's "json"
// tag, or the empty string. It does not include the leading comma.
type tagOptions string

// parseTag splits a struct field's json tag into its name and
// comma-separated options.
func parseTag(tag string) (string, tagOptions) {
	tag, opt, _ := strings.Cut(tag, ",")
	return tag, tagOptions(opt)
}

// Contains reports whether a comma-separated list of options
// contains a particular substr flag. substr must be surrounded by a
// string boundary or commas.
func (o tagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var name string
		name, s, _ = strings.Cut(s, ",")
		if name == optionName {
			return true
		}
	}
	return false
}

// 整合viper,提供单个注入,viper实现也挺简单的,原来的方案就能实现啊
/*func Inject(v Config, path string) {

}
*/
