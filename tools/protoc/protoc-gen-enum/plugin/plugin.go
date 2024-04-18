package plugin

import (
	stringsi "github.com/hopeio/cherry/utils/strings"
	"google.golang.org/protobuf/compiler/protogen"
	"strconv"
	"strings"
)

type pathType int

const (
	pathTypeImport pathType = iota
	pathTypeSourceRelative
)

type Builder struct {
	plugin *protogen.Plugin
}

func NewBuilder(gen *protogen.Plugin) *Builder {
	return &Builder{
		plugin: gen,
	}
}

func parseParameter(param string) map[string]string {
	paramMap := make(map[string]string)

	for _, p := range strings.Split(param, ",") {
		if i := strings.Index(p, "="); i < 0 {
			paramMap[p] = ""
		} else {
			paramMap[p[0:i]] = p[i+1:]
		}
	}

	return paramMap
}

func (b *Builder) Generate() error {

	genFileMap := make(map[string]*protogen.GeneratedFile)

	for _, protoFile := range b.plugin.Files {
		if !protoFile.Generate {
			continue
		}

		if TurnOffExtGenAll(protoFile) || len(protoFile.Enums) == 0 {
			continue
		}

		fileName := protoFile.GeneratedFilenamePrefix
		g := b.plugin.NewGeneratedFile(fileName+".enumext.pb.go", ".")
		genFileMap[fileName] = g
		// third traverse: build associations
		for _, enum := range protoFile.Enums {
			if TurnOffExtGen(enum) {
				genFileMap[fileName] = g
				break
			}
		}

	}

	for _, protoFile := range b.plugin.Files {
		fileName := protoFile.GeneratedFilenamePrefix
		g, ok := genFileMap[fileName]
		if !ok || len(protoFile.Enums) == 0 {
			continue
		}

		g.P("package ", protoFile.GoPackageName)

		for _, enum := range protoFile.Enums {
			b.generate(protoFile, enum, g)
		}

	}

	return nil
}

func (b *Builder) generate(f *protogen.File, e *protogen.Enum, g *protogen.GeneratedFile) {

	if EnabledEnumStringer(e) {
		b.generateString(f, e, g)
	}
	if EnabledEnumJsonMarshal(f, e) {
		b.generateJsonMarshal(f, e, g)
	}
	if EnabledEnumErrorCode(e) {
		b.generateErrorCode(f, e, g)
	}
	if EnabledEnumGqlGen(f, e) {
		b.generateGQLMarshal(f, e, g)
	}
}

func (b *Builder) generateString(f *protogen.File, e *protogen.Enum, g *protogen.GeneratedFile) {
	ccTypeName := stringsi.CamelCase(e.Desc.Name())

	g.P("func (x ", ccTypeName, ") String() string {")
	g.P()
	if len(e.Values) > 64 {
		g.P("return ", ccTypeName, "_name[x]")
	} else {
		g.P("switch x {")
		for _, ev := range e.Values {
			name := stringsi.CamelCase(ev.Desc.Name())
			//PrintComments(e.Comments, g)
			value := name
			g.P("case ", name, " :")
			if cn := GetEnumValueCN(ev); cn != "" {
				value = cn
			}
			g.P("return ", strconv.Quote(value))
		}
	}
	g.P("}")
	g.P("return ", strconv.Quote(""))
	g.P("}")
	g.P()
}

func (b *Builder) generateGQLMarshal(f *protogen.File, e *protogen.Enum, g *protogen.GeneratedFile) {
	ccTypeName := stringsi.CamelCase(e.Desc.Name())

	typ := "uint32"
	if typ1 := GetEnumType(e); typ1 != "" {
		typ = typ1
	}
	g.P("func (x ", ccTypeName, ") MarshalGQL(w ", generateImport("Writer", "io", g), ") {")
	g.P(`w.Write(`, generateImport("QuoteToBytes", "github.com/hopeio/cherry/utils/strings", g), `(x.String()))`)
	g.P("}")
	g.P()
	g.P("func (x *", ccTypeName, ") UnmarshalGQL(v interface{}) error {")
	g.P(`if i, ok := v.(`, typ, "); ok {")
	g.P(`*x = `, ccTypeName, `(i)`)
	g.P("return nil")
	g.P("}")
	g.P(`return `, generateImport("New", "errors", g), `("enum need integer type")`)
	g.P("}")
	g.P()
}

func (b *Builder) generateJsonMarshal(f *protogen.File, e *protogen.Enum, g *protogen.GeneratedFile) {
	ccTypeName := stringsi.CamelCase(e.Desc.Name())

	g.P("func (x ", ccTypeName, ") MarshalJSON() ([]byte, error) {")
	g.P("return ", generateImport("QuoteToBytes", "github.com/hopeio/cherry/utils/strings", g), "(x.String())", ", nil")
	g.P("}")
	g.P()
	g.P("func (x *", ccTypeName, ") UnmarshalJSON(data []byte) error {")

	g.P("value, ok := ", ccTypeName, `_value[string(data)]`)
	g.P("if ok {")

	g.P("*x = ", ccTypeName, "(value)")
	g.P("return nil")

	g.P("}")
	g.P(`return `, generateImport("New", "errors", g), `("invalid`, ccTypeName, `")`)

	g.P("}")
	g.P()
}

func (b *Builder) generateErrorCode(f *protogen.File, e *protogen.Enum, g *protogen.GeneratedFile) {
	ccTypeName := stringsi.CamelCase(e.Desc.Name())

	g.P("func (x ", ccTypeName, ") Error() string {")

	g.P(`return x.String()`)

	g.P("}")
	g.P()
	g.P("func (x ", ccTypeName, ") ErrRep() *", generateImport("ErrRep", "github.com/hopeio/cherry/protobuf/errorcode", g), " {")

	g.P(`return &errorcode.ErrRep{Code: errorcode.ErrCode(x), Message: x.String()}`)

	g.P("}")
	g.P()
	g.P("func (x ", ccTypeName, ") Message(msg string) error {")

	g.P(`return &errorcode.ErrRep{Code: errorcode.ErrCode(x), Message: msg}`)

	g.P("}")
	g.P()
	g.P("func (x ", ccTypeName, ") ErrorLog(err error) error {")

	g.P(generateImport("Error", "github.com/hopeio/cherry/utils/log", g), `(err)`)
	g.P(`return &errorcode.ErrRep{Code: errorcode.ErrCode(x), Message: x.String()}`)

	g.P("}")
	g.P()
	g.P("func (x ", ccTypeName, ") GRPCStatus() *", generateImport("Status", "google.golang.org/grpc/status", g), " {")

	g.P(`return `, `status.New(`, generateImport("Code", "google.golang.org/grpc/codes", g), `(x), x.String())`)

	g.P("}")
	g.P()
}

func generateImport(name string, importPath string, g *protogen.GeneratedFile) string {
	return g.QualifiedGoIdent(protogen.GoIdent{
		GoName:       name,
		GoImportPath: protogen.GoImportPath(importPath),
	})
}
