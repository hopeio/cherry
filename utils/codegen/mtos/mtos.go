package mtos

import (
	"fmt"
	stringsi "github.com/hopeio/cherry/utils/strings"
	"reflect"
	"strings"
)

func Gen(m map[string]any, tag string) {
	var structFields strings.Builder
	for key, value := range m {
		dataType := reflect.TypeOf(value)
		fieldCode := generateFieldCode(key, dataType, tag)
		if fieldCode != "" {
			structFields.WriteString(fieldCode + "\n")
		}
	}
}

func generateFieldCode(k string, v any, tag string) string {
	fieldNameCapitalized := stringsi.SnakeToCamel(k)
	switch fieldType.Kind() {
	case reflect.String:
		return fmt.Sprintf("\t%s string `%s:\"%s\"`\n", tag, fieldNameCapitalized, fieldName)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("\t%s int `%s:\"%s\"`\n", fieldNameCapitalized, "json", fieldName)
	case reflect.Bool:
		return fmt.Sprintf("\t%s bool `json:\"%s\"`\n", fieldNameCapitalized, fieldName)
	case reflect.Map:
		innerStructName := fieldNameCapitalized + "Details"
		innerStructCode := mapToStruct(fieldType.Interface().(map[string]interface{}), innerStructName)
		return fmt.Sprintf("\t%s %s `json:\"%s\"`\n", fieldNameCapitalized, innerStructName, fieldName)
	case reflect.Slice:
		// 这里仅示例处理，实际情况需要更详细的检查和处理
		elementType := fieldType.Elem()
		if elementType.Kind() == reflect.String || elementType.Kind() == reflect.Int || elementType.Kind() == reflect.Bool {
			sliceType := "[]" + elementType.Kind().String()
			return fmt.Sprintf("\t%s %s `json:\"%s\"`\n", fieldNameCapitalized, sliceType, fieldName)
		} else {
			log.Printf("Unsupported slice element type for field %s: %v", fieldName, elementType)
			return ""
		}
	default:
		log.Printf("Unsupported type for field %s: %v", fieldName, fieldType)
		return ""
	}
}
