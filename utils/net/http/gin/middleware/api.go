package middleware

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/spec"
	"github.com/hopeio/cherry/utils/net/http/api/apidoc"
	"github.com/hopeio/cherry/utils/reflect"
)

// Deprecated
func ApiMiddle(ctx *gin.Context) {
	currentRouteName := ctx.Request.RequestURI[len(ctx.Request.Method):]

	var pathItem *spec.PathItem

	doc := apidoc.GetDoc("../swagger.json")

	if doc.Paths != nil && doc.Paths.Paths != nil {
		if path, ok := doc.Paths.Paths[currentRouteName]; ok {
			pathItem = &path
		} else {
			pathItem = new(spec.PathItem)
		}
	} else {
		doc.Paths = &spec.Paths{Paths: map[string]spec.PathItem{}}
		pathItem = new(spec.PathItem)
	}

	parameters := make([]spec.Parameter, len(ctx.Params), len(ctx.Params))

	params := ctx.Params

	for i := range params {
		key := params[i].Key

		//val := params[i].ValueRaw
		parameters[i] = spec.Parameter{
			ParamProps: spec.ParamProps{
				Name:        key,
				In:          "path",
				Description: "Description",
			},
		}
	}

	if stop, _ := ctx.GetQuery("apidoc"); stop == "stop" {
		defer apidoc.WriteToFile("../", "")
	}

	var res spec.Responses
	op := spec.Operation{
		OperationProps: spec.OperationProps{
			Description: "Description",
			Consumes:    []string{"application/x-www-form-urlencoded"},
			Tags:        []string{"Tags"},
			Summary:     "Summary",
			ID:          "currentRouteName" + ctx.Request.Method,
			Parameters:  parameters,
			Responses:   &res,
		},
	}

	switch ctx.Request.Method {
	case http.MethodGet:
		pathItem.Get = &op
	case http.MethodPost:
		pathItem.Post = &op
	case http.MethodPut:
		pathItem.Put = &op
	case http.MethodDelete:
		pathItem.Delete = &op
	case http.MethodOptions:
		pathItem.Options = &op
	case http.MethodPatch:
		pathItem.Patch = &op
	case http.MethodHead:
		pathItem.Head = &op
	}
	doc.Paths.Paths[currentRouteName] = *pathItem
	ctx.Next()
}

func DefinitionsApi(definitions map[string]spec.Schema, v interface{}, exclude []string) {
	schema := spec.Schema{
		SchemaProps: spec.SchemaProps{
			Type:       []string{"object"},
			Properties: make(map[string]spec.Schema),
		},
	}

	body := reflect.TypeOf(v).Elem()
	var typ, subFieldName string
	for i := 0; i < body.NumField(); i++ {
		json := strings.Split(body.Field(i).Tag.Get("json"), ",")[0]
		if json == "" || json == "-" {
			continue
		}
		fieldType := body.Field(i).Type
		switch fieldType.Kind() {
		case reflect.Struct:
			typ = "object"
			v = reflect.ValueOf(v).Elem().Field(i).Addr().Interface()
			subFieldName = fieldType.Name()
		case reflect.Ptr:
			typ = "object"
			v = reflect.New(fieldType.Elem()).Interface()
			subFieldName = fieldType.Elem().Name()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			typ = "integer"
		case reflect.Array, reflect.Slice:
			typ = "array"
			v = reflect.New(reflect.OriginalType(fieldType)).Interface()
			subFieldName = reflect.OriginalType(fieldType).Name()
		case reflect.Float32, reflect.Float64:
			typ = "number"
		case reflect.String:
			typ = "string"
		case reflect.Bool:
			typ = "boolean"

		}
		subSchema := spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{typ},
			},
		}
		if typ == "object" {
			subSchema.Ref = spec.MustCreateRef("#/definitions/" + subFieldName)
			DefinitionsApi(definitions, v, nil)
		}
		if typ == "array" {
			subSchema.Items = new(spec.SchemaOrArray)
			subSchema.Items.Schema = &spec.Schema{}
			subSchema.Items.Schema.Ref = spec.MustCreateRef("#/definitions/" + subFieldName)
			DefinitionsApi(definitions, v, nil)
		}
		schema.Properties[json] = subSchema
	}
	definitions[body.Name()] = schema
}

func genSchema(v interface{}) *spec.Schema {
	return nil
}
