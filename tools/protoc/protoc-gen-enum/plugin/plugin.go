package plugin

import (
	"github.com/hopeio/cherry/tools/protoc/protoc-gen-enum/options"
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
	plugin          *protogen.Plugin
	importStatus    protogen.GoImportPath
	importCodes     protogen.GoImportPath
	importLog       protogen.GoImportPath
	importStrings   protogen.GoImportPath
	importStrconv   protogen.GoImportPath
	importErrorcode protogen.GoImportPath
	importErrors    protogen.GoImportPath
	importIo        protogen.GoImportPath
}

func NewBuilder(gen *protogen.Plugin) *Builder {
	return &Builder{
		plugin:          gen,
		importStatus:    "google.golang.org/grpc/status",
		importCodes:     "google.golang.org/grpc/codes",
		importLog:       "github.com/hopeio/cherry/utils/log",
		importStrings:   "github.com/hopeio/cherry/utils/strings",
		importStrconv:   "strconv",
		importErrorcode: "github.com/hopeio/cherry/protobuf/errorcode",
		importErrors:    "errors",
		importIo:        "io",
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
		g := b.plugin.NewGeneratedFile(fileName+".enumext.pb.go", protoFile.GoImportPath)
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
		b.generateJsonMarshal(e, g)
	}
	if EnabledEnumErrorCode(e) {
		b.generateErrorCode(e, g)
	}
	if EnabledEnumGqlGen(f, e) {
		b.generateGQLMarshal(e, g)
	}
}

func (b *Builder) generateString(f *protogen.File, e *protogen.Enum, g *protogen.GeneratedFile) {
	noEnumPrefix := options.FileOptions(f).GetNoEnumPrefix()
	ccTypeName := e.GoIdent

	g.P("func (x ", ccTypeName, ") String() string {")
	g.P()
	if len(e.Values) > 64 {
		g.P("return ", ccTypeName, "_name[x]")
	} else {
		g.P("switch x {")
		for _, ev := range e.Values {
			opts := options.ValueOptions(ev)
			name := opts.GetName()
			if name == "" {
				name = ev.GoIdent.GoName
				if noEnumPrefix {
					name = replacePrefix(ev.GoIdent.GoName, e.GoIdent.GoName+"_", "")
				}
			}

			//PrintComments(e.Comments, g)

			g.P("case ", name, " :")
			if cn := GetEnumValueCN(ev); cn != "" {
				g.P("return ", strconv.Quote(cn))
			} else {
				g.P("return ", strconv.Quote(name))
			}

		}
	}
	g.P("}")
	g.P("return ", strconv.Quote(""))
	g.P("}")
	g.P()
}

func (b *Builder) generateGQLMarshal(e *protogen.Enum, g *protogen.GeneratedFile) {
	ccTypeName := e.GoIdent

	typ := "uint32"
	if typ1 := GetEnumType(e); typ1 != "" {
		typ = typ1
	}
	g.P("func (x ", ccTypeName, ") MarshalGQL(w ", b.importIo.Ident("Writer"), ") {")
	g.P(`w.Write(`, b.importStrings.Ident("QuoteToBytes"), `(x.String()))`)
	g.P("}")
	g.P()
	g.P("func (x *", ccTypeName, ") UnmarshalGQL(v interface{}) error {")
	g.P(`if i, ok := v.(`, typ, "); ok {")
	g.P(`*x = `, ccTypeName, `(i)`)
	g.P("return nil")
	g.P("}")
	g.P(`return `, b.importErrors.Ident("New"), `("enum need integer type")`)
	g.P("}")
	g.P()
}

func (b *Builder) generateJsonMarshal(e *protogen.Enum, g *protogen.GeneratedFile) {
	ccTypeName := e.GoIdent

	g.P("func (x ", ccTypeName, ") MarshalJSON() ([]byte, error) {")
	g.P("return ", b.importStrings.Ident("QuoteToBytes"), "(x.String())", ", nil")
	g.P("}")
	g.P()
	g.P("func (x *", ccTypeName, ") UnmarshalJSON(data []byte) error {")

	g.P("if len(data) > 0 && data[0] == '\"' {")

	g.P("value, ok := ", ccTypeName, `_value[string(data[1:len(data)-1])]`)
	g.P("if ok {")

	g.P("*x = ", ccTypeName, "(value)")
	g.P("return nil")
	g.P("}")
	g.P("} else {")
	g.P("value, err := ", b.importStrconv.Ident("ParseInt"), `(string(data), 10, 32)`)
	g.P("if err == nil {")
	g.P("_, ok := ", ccTypeName, `_name[int32(value)]`)
	g.P("if ok {")
	g.P("*x = ", ccTypeName, "(value)")
	g.P("return nil")
	g.P("}")
	g.P("}")
	g.P("}")
	g.P(`return `, b.importErrors.Ident("New"), `("invalid enum value: `, ccTypeName, `")`)

	g.P("}")
	g.P()
}

func (b *Builder) generateErrorCode(e *protogen.Enum, g *protogen.GeneratedFile) {
	ccTypeName := e.GoIdent

	g.P("func (x ", ccTypeName, ") Error() string {")

	g.P(`return x.String()`)

	g.P("}")
	g.P()
	g.P("func (x ", ccTypeName, ") ErrRep() *", b.importErrorcode.Ident("ErrRep"), " {")

	g.P(`return &errorcode.ErrRep{Code: errorcode.ErrCode(x), Message: x.String()}`)

	g.P("}")
	g.P()
	g.P("func (x ", ccTypeName, ") Message(msg string) error {")

	g.P(`return &errorcode.ErrRep{Code: errorcode.ErrCode(x), Message: msg}`)

	g.P("}")
	g.P()
	g.P("func (x ", ccTypeName, ") ErrorLog(err error) error {")

	g.P(b.importLog.Ident("Error"), `(err)`)
	g.P(`return &errorcode.ErrRep{Code: errorcode.ErrCode(x), Message: x.String()}`)

	g.P("}")
	g.P()
	g.P("func (x ", ccTypeName, ") GrpcStatus() *", b.importStatus.Ident("Status"), " {")

	g.P(`return `, `status.New(`, b.importCodes.Ident("Code"), `(x), x.String())`)

	g.P("}")
	g.P()
}

func replacePrefix(s, prefix, with string) string {
	return with + strings.TrimPrefix(s, prefix)
}
