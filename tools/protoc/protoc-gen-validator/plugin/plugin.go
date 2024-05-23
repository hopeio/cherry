package plugin

import (
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/hopeio/cherry/protobuf/utils/validator"
	"google.golang.org/protobuf/proto"
	descriptor "google.golang.org/protobuf/types/descriptorpb"
	"os"
	"reflect"
	"strconv"
)

const uuidPattern = "^([a-fA-F0-9]{8}-" +
	"[a-fA-F0-9]{4}-" +
	"[%s][a-fA-F0-9]{3}-" +
	"[8|9|aA|bB][a-fA-F0-9]{3}-" +
	"[a-fA-F0-9]{12})?$"

type plugin struct {
	*protogen.Plugin
	regexPkg     protogen.GoImportPath
	fmtPkg       protogen.GoImportPath
	errorsPkg    protogen.GoImportPath
	validatorPkg protogen.GoImportPath
}

func New(p *protogen.Plugin) *plugin {
	return &plugin{
		Plugin:       p,
		regexPkg:     "regexp",
		fmtPkg:       "fmt",
		errorsPkg:    "errors",
		validatorPkg: "github.com/hopeio/cherry/utils/validation/validator",
	}
}

func (p *plugin) Name() string {
	return "validator"
}

func (p *plugin) Generate() error {
	genFileMap := make(map[string]*protogen.GeneratedFile)

	for _, protoFile := range p.Files {
		if !protoFile.Generate {
			continue
		}

		if len(protoFile.Messages) == 0 {
			continue
		}

		fileName := protoFile.GeneratedFilenamePrefix
		g := p.NewGeneratedFile(fileName+".validator.pb.go", protoFile.GoImportPath)

		genFileMap[fileName] = g
		g.P("package ", protoFile.GoPackageName)

		for _, msg := range protoFile.Messages {
			for _, nmsg := range msg.Messages {
				p.gen(protoFile, nmsg, g)
			}
			p.gen(protoFile, msg, g)
		}
	}
	return nil
}

func (p *plugin) gen(file *protogen.File, message *protogen.Message, g *protogen.GeneratedFile) {
	if message.Desc.Options().(*descriptor.MessageOptions).GetMapEntry() {
		return
	}
	p.generateRegexVars(message, g)
	if file.Proto.GetSyntax() == "proto3" {
		p.generateProto3Message(message, g)
	} else {
		p.generateProto2Message(message, g)
	}
}

func getFieldValidatorIfAny(field *protogen.Field) *validator.FieldValidator {
	if field.Desc.Options() != nil {
		v := proto.GetExtension(field.Desc.Options(), validator.E_Field)
		if v.(*validator.FieldValidator) != nil {
			return v.(*validator.FieldValidator)
		}
	}
	return nil
}

func getOneofValidatorIfAny(oneof *protogen.Oneof) *validator.OneofValidator {
	if oneof.Desc.Options() != nil {
		v := proto.GetExtension(oneof.Desc.Options(), validator.E_Oneof)
		if v.(*validator.OneofValidator) != nil {
			return (v.(*validator.OneofValidator))
		}
	}
	return nil
}

func (p *plugin) isSupportedInt(field *protogen.Field) bool {
	switch field.Desc.Kind() {
	case protoreflect.Int32Kind, protoreflect.Int64Kind:
		return true
	case protoreflect.Uint32Kind, protoreflect.Uint64Kind:
		return true
	case protoreflect.Sint32Kind, protoreflect.Sint64Kind:
		return true
	}
	return false
}

func (p *plugin) isSupportedFloat(field *protogen.Field) bool {
	switch field.Desc.Kind() {
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		return true
	case protoreflect.Fixed32Kind, protoreflect.Fixed64Kind:
		return true
	case protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind:
		return true
	}
	return false
}

func (p *plugin) generateRegexVars(message *protogen.Message, g *protogen.GeneratedFile) {
	ccTypeName := message.GoIdent
	for _, field := range message.Fields {
		validator := getFieldValidatorIfAny(field)
		if validator != nil {
			fieldName := field.GoName
			if validator.Regex != nil && validator.UuidVer != nil {
				fmt.Fprintf(os.Stderr, "WARNING: regex and uuid validator is set for field %v.%v, only one of them can be set. Regex and UUID validator is ignored for this field.", ccTypeName, fieldName)
			} else if validator.UuidVer != nil {
				uuid, err := getUUIDRegex(validator.UuidVer)
				if err != nil {
					fmt.Fprintf(os.Stderr, "WARNING: field %v.%v error %s.\n", ccTypeName, fieldName, err)
				} else {
					validator.Regex = &uuid
					g.P(`var `, p.regexName(ccTypeName.GoName, fieldName), ` = `, p.regexPkg.Ident("MustCompile"), `(`, "`", *validator.Regex, "`", `)`)
				}
			} else if validator.Regex != nil {
				g.P(`var `, p.regexName(ccTypeName.GoName, fieldName), ` = `, p.regexPkg.Ident("MustCompile"), `(`, "`", *validator.Regex, "`", `)`)
			}
		}
	}
}

func (p *plugin) generateProto2Message(message *protogen.Message, g *protogen.GeneratedFile) {
	ccTypeName := message.GoIdent

	g.P(`func (x *`, ccTypeName, `) Validate() error {`)
	for _, field := range message.Fields {
		fieldName := field.GoName
		fieldValidator := getFieldValidatorIfAny(field)
		if fieldValidator == nil && field.Desc.Kind() != protoreflect.MessageKind {
			continue
		}
		if p.validatorWithMessageExists(fieldValidator) {
			fmt.Fprintf(os.Stderr, "WARNING: field %v.%v is a proto2 message, validator.msg_exists has no effect\n", ccTypeName, fieldName)
		}
		variableName := "x." + fieldName
		repeated := field.Desc.IsList()

		nullable := field.Parent != nil
		// For proto2 syntax, only Gogo generates non-pointer fields

		if repeated {
			p.generateRepeatedCountValidator(variableName, fieldName, fieldValidator, g)
			if field.Desc.Kind() == protoreflect.MessageKind || p.validatorWithNonRepeatedConstraint(fieldValidator) {
				g.P(`for _, item := range `, variableName, `{`)
				variableName = "item"
			}
		} else if nullable {
			g.P(`if `, variableName, ` != nil {`)
			if field.Desc.Kind() != protoreflect.BytesKind {
				variableName = "*(" + variableName + ")"
			}
		} else if field.Desc.Kind() != protoreflect.MessageKind {
			variableName = `x.Get` + fieldName + `()`
		}

		if !repeated && fieldValidator != nil {
			if fieldValidator.RepeatedCountMin != nil {
				fmt.Fprintf(os.Stderr, "WARNING: field %v.%v is not repeated, validator.min_elts has no effects\n", ccTypeName, fieldName)
			}
			if fieldValidator.RepeatedCountMax != nil {
				fmt.Fprintf(os.Stderr, "WARNING: field %v.%v is not repeated, validator.max_elts has no effects\n", ccTypeName, fieldName)
			}
		}
		if field.Desc.Kind() == protoreflect.StringKind {
			p.generateStringValidator(variableName, ccTypeName.GoName, fieldName, fieldValidator, g)
		} else if p.isSupportedInt(field) {
			p.generateIntValidator(variableName, ccTypeName.GoName, fieldName, fieldValidator, g)
		} else if field.Desc.Kind() == protoreflect.EnumKind {
			p.generateEnumValidator(field, variableName, ccTypeName.GoName, fieldName, fieldValidator, g)
		} else if p.isSupportedFloat(field) {
			p.generateFloatValidator(variableName, ccTypeName.GoName, fieldName, fieldValidator, g)
		} else if field.Desc.Kind() == protoreflect.BytesKind {
			p.generateLengthValidator(variableName, fieldName, fieldValidator, g)
		} else if field.Desc.Kind() == protoreflect.MessageKind {
			if repeated && nullable {
				variableName = "*(item)"
			}
			g.P(`if err := `, p.validatorPkg.Ident("CallValidatorIfExists"), `(&(`, variableName, `)); err != nil {`)
			g.P(`return `, p.validatorPkg.Ident("FieldError"), `("`, fieldName, `", err)`)
			g.P(`}`)
		}
		if repeated {
			// end the repeated loop
			if field.Desc.Kind() == protoreflect.MessageKind || p.validatorWithNonRepeatedConstraint(fieldValidator) {
				// This internal 'if' cannot be refactored as it would change semantics with respect to the corresponding prelude 'if's
				g.P(`}`)
			}
		} else if nullable {
			g.P(`}`)
		}
	}
	g.P(`return nil`)
	g.P(`}`)
}

func (p *plugin) generateProto3Message(message *protogen.Message, g *protogen.GeneratedFile) {
	ccTypeName := message.GoIdent
	g.P(`func (x *`, ccTypeName, `) Validate() error {`)

	for _, oneof := range message.Oneofs {
		oneofValidator := getOneofValidatorIfAny(oneof)
		if oneofValidator == nil {
			continue
		}
		if oneofValidator.GetRequired() {
			oneOfName := oneof.GoName
			g.P(`if x.Get` + oneOfName + `() == nil {`)
			g.P(`return `, p.validatorPkg.Ident("FieldError"), `("`, oneOfName, `",`, p.fmtPkg.Ident("Errorf"), `("one of the fields must be set"))`)
			g.P(`}`)
		}
	}
	for _, field := range message.Fields {
		fieldValidator := getFieldValidatorIfAny(field)
		if fieldValidator == nil && field.Desc.Kind() != protoreflect.MessageKind {
			continue
		}
		isOneOf := field.Oneof != nil
		fieldName := field.GoName
		variableName := "x." + fieldName
		repeated := field.Desc.IsList()
		// Golang's proto3 has no concept of unset primitive fields
		nullable := field.Desc.Kind() == protoreflect.MessageKind && field.Parent != nil
		if field.Desc.IsMap() {
			g.P(`// Validation of proto3 map<> fields is unsupported.`)
			continue
		}
		if isOneOf {
			// if x, ok := m.GetType().(*OneOfMessage3_OneInt); ok {
			g.P("if x, ok := x.Get", field.Oneof.GoName, "().(*", field.GoIdent, "); ok {")
			variableName = "x." + field.GoName
		}
		if repeated {
			p.generateRepeatedCountValidator(variableName, fieldName, fieldValidator, g)
			if field.Desc.Kind() == protoreflect.MessageKind || p.validatorWithNonRepeatedConstraint(fieldValidator) {
				g.P(`for _, item := range `, variableName, `{`)
				variableName = "item"
			}
		} else if fieldValidator != nil {
			if fieldValidator.RepeatedCountMin != nil {
				fmt.Fprintf(os.Stderr, "WARNING: field %v.%v is not repeated, validator.min_elts has no effects\n", ccTypeName, fieldName)
			}
			if fieldValidator.RepeatedCountMax != nil {
				fmt.Fprintf(os.Stderr, "WARNING: field %v.%v is not repeated, validator.max_elts has no effects\n", ccTypeName, fieldName)
			}
		}
		if field.Desc.Kind() == protoreflect.StringKind {
			p.generateStringValidator(variableName, ccTypeName.GoName, fieldName, fieldValidator, g)
		} else if p.isSupportedInt(field) {
			p.generateIntValidator(variableName, ccTypeName.GoName, fieldName, fieldValidator, g)
		} else if field.Desc.Kind() == protoreflect.EnumKind {
			p.generateEnumValidator(field, variableName, ccTypeName.GoName, fieldName, fieldValidator, g)
		} else if p.isSupportedFloat(field) {
			p.generateFloatValidator(variableName, ccTypeName.GoName, fieldName, fieldValidator, g)
		} else if field.Desc.Kind() == protoreflect.BytesKind {
			p.generateLengthValidator(variableName, fieldName, fieldValidator, g)
		} else if field.Desc.Kind() == protoreflect.MessageKind {
			if p.validatorWithMessageExists(fieldValidator) {
				if nullable && !repeated {
					g.P(`if nil == `, variableName, `{`)
					g.P(`return `, p.validatorPkg.Ident("FieldError"), `("`, fieldName, `",`, p.fmtPkg.Ident("Errorf"), `("message must exist"))`)
					g.P(`}`)
				} else if repeated {
					fmt.Fprintf(os.Stderr, "WARNING: field %v.%v is repeated, validator.msg_exists has no effect\n", ccTypeName, fieldName)
				} else if !nullable {
					fmt.Fprintf(os.Stderr, "WARNING: field %v.%v is a nullable=false, validator.msg_exists has no effect\n", ccTypeName, fieldName)
				}
			}

			if nullable {
				g.P(`if `, variableName, ` != nil {`)
			} else {
				variableName = "&(" + variableName + ")"
			}
			g.P(`if err := `, p.validatorPkg.Ident("CallValidatorIfExists"), `(`, variableName, `); err != nil {`)
			g.P(`return `, p.validatorPkg.Ident("FieldError"), `("`, fieldName, `", err)`)
			g.P(`}`)
			if nullable {
				g.P(`}`)
			}
		}
		if repeated && (field.Desc.Kind() == protoreflect.MessageKind || p.validatorWithNonRepeatedConstraint(fieldValidator)) {
			// end the repeated loop
			g.P(`}`)
		}
		if isOneOf {
			// end the oneof if statement
			g.P(`}`)
		}
	}
	g.P(`return nil`)
	g.P(`}`)
}

func (p *plugin) generateIntValidator(variableName string, ccTypeName string, fieldName string, fv *validator.FieldValidator, g *protogen.GeneratedFile) {
	if fv.IntGt != nil {
		g.P(`if !(`, variableName, ` > `, fv.GetIntGt(), `) {`)
		errorStr := fmt.Sprintf(`be greater than '%d'`, fv.GetIntGt())
		p.generateErrorString(variableName, fieldName, errorStr, fv, g)
		g.P(`}`)
	}
	if fv.IntLt != nil {
		g.P(`if !(`, variableName, ` < `, fv.GetIntLt(), `) {`)
		errorStr := fmt.Sprintf(`be less than '%d'`, fv.GetIntLt())
		p.generateErrorString(variableName, fieldName, errorStr, fv, g)
		g.P(`}`)
	}
}

func (p *plugin) generateEnumValidator(
	field *protogen.Field,
	variableName, ccTypeName, fieldName string,
	fv *validator.FieldValidator, g *protogen.GeneratedFile) {
	if fv.GetIsInEnum() {
		enum := field.Enum
		g.P(`if _, ok := `, enum.GoIdent.GoName, "_name[int32(", variableName, ")]; !ok {")
		p.generateErrorString(variableName, fieldName, fmt.Sprintf("be a valid %s field", enum.GoIdent.GoName), fv, g)
		g.P(`}`)
	}
}

func (p *plugin) generateLengthValidator(variableName string, fieldName string, fv *validator.FieldValidator, g *protogen.GeneratedFile) {
	if fv.LengthGt != nil {
		g.P(`if !( len(`, variableName, `) > `, fv.GetLengthGt(), `) {`)
		errorStr := fmt.Sprintf(`have a length greater than '%d'`, fv.GetLengthGt())
		p.generateErrorString(variableName, fieldName, errorStr, fv, g)
		g.P(`}`)
	}

	if fv.LengthLt != nil {
		g.P(`if !( len(`, variableName, `) < `, fv.GetLengthLt(), `) {`)
		errorStr := fmt.Sprintf(`have a length smaller than '%d'`, fv.GetLengthLt())
		p.generateErrorString(variableName, fieldName, errorStr, fv, g)
		g.P(`}`)
	}

	if fv.LengthEq != nil {
		g.P(`if !( len(`, variableName, `) == `, fv.GetLengthEq(), `) {`)
		errorStr := fmt.Sprintf(`have a length equal to '%d'`, fv.GetLengthEq())
		p.generateErrorString(variableName, fieldName, errorStr, fv, g)
		g.P(`}`)
	}
}

func (p *plugin) generateFloatValidator(variableName string, ccTypeName string, fieldName string, fv *validator.FieldValidator, g *protogen.GeneratedFile) {
	upperIsStrict := true
	lowerIsStrict := true

	// First check for incompatible constraints (i.e flt_lt & flt_lte both defined, etc) and determine the real limits.
	if fv.FloatEpsilon != nil && fv.FloatLt == nil && fv.FloatGt == nil {
		fmt.Fprintf(os.Stderr, "WARNING: field %v.%v has no 'float_lt' or 'float_gt' field so setting 'float_epsilon' has no effect.", ccTypeName, fieldName)
	}
	if fv.FloatLt != nil && fv.FloatLte != nil {
		fmt.Fprintf(os.Stderr, "WARNING: field %v.%v has both 'float_lt' and 'float_lte' constraints, only the strictest will be used.", ccTypeName, fieldName)
		strictLimit := fv.GetFloatLt()
		if fv.FloatEpsilon != nil {
			strictLimit += fv.GetFloatEpsilon()
		}

		if fv.GetFloatLte() < strictLimit {
			upperIsStrict = false
		}
	} else if fv.FloatLte != nil {
		upperIsStrict = false
	}

	if fv.FloatGt != nil && fv.FloatGte != nil {
		fmt.Fprintf(os.Stderr, "WARNING: field %v.%v has both 'float_gt' and 'float_gte' constraints, only the strictest will be used.", ccTypeName, fieldName)
		strictLimit := fv.GetFloatGt()
		if fv.FloatEpsilon != nil {
			strictLimit -= fv.GetFloatEpsilon()
		}

		if fv.GetFloatGte() > strictLimit {
			lowerIsStrict = false
		}
	} else if fv.FloatGte != nil {
		lowerIsStrict = false
	}

	// Generate the constraint checking code.
	errorStr := ""
	compareStr := ""
	if fv.FloatGt != nil || fv.FloatGte != nil {
		compareStr = fmt.Sprint(`if !(`, variableName)
		if lowerIsStrict {
			errorStr = fmt.Sprintf(`be strictly greater than '%g'`, fv.GetFloatGt())
			if fv.FloatEpsilon != nil {
				errorStr += fmt.Sprintf(` with a tolerance of '%g'`, fv.GetFloatEpsilon())
				compareStr += fmt.Sprint(` + `, fv.GetFloatEpsilon())
			}
			compareStr += fmt.Sprint(` > `, fv.GetFloatGt(), `) {`)
		} else {
			errorStr = fmt.Sprintf(`be greater than or equal to '%g'`, fv.GetFloatGte())
			compareStr += fmt.Sprint(` >= `, fv.GetFloatGte(), `) {`)
		}
		g.P(compareStr)
		p.generateErrorString(variableName, fieldName, errorStr, fv, g)
		g.P(`}`)
	}

	if fv.FloatLt != nil || fv.FloatLte != nil {
		compareStr = fmt.Sprint(`if !(`, variableName)
		if upperIsStrict {
			errorStr = fmt.Sprintf(`be strictly lower than '%g'`, fv.GetFloatLt())
			if fv.FloatEpsilon != nil {
				errorStr += fmt.Sprintf(` with a tolerance of '%g'`, fv.GetFloatEpsilon())
				compareStr += fmt.Sprint(` - `, fv.GetFloatEpsilon())
			}
			compareStr += fmt.Sprint(` < `, fv.GetFloatLt(), `) {`)
		} else {
			errorStr = fmt.Sprintf(`be lower than or equal to '%g'`, fv.GetFloatLte())
			compareStr += fmt.Sprint(` <= `, fv.GetFloatLte(), `) {`)
		}
		g.P(compareStr)
		p.generateErrorString(variableName, fieldName, errorStr, fv, g)
		g.P(`}`)
	}
}

// getUUIDRegex returns a regex to validate that a string is in UUID
// format. The version parameter specified the UUID version. If version is 0,
// the returned regex is valid for any UUID version
func getUUIDRegex(version *int32) (string, error) {
	if version == nil {
		return "", nil
	} else if *version < 0 || *version > 5 {
		return "", fmt.Errorf("UUID version should be between 0-5, Got %d", *version)
	} else if *version == 0 {
		return fmt.Sprintf(uuidPattern, "1-5"), nil
	} else {
		return fmt.Sprintf(uuidPattern, strconv.Itoa(int(*version))), nil
	}
}

func (p *plugin) generateStringValidator(variableName string, ccTypeName string, fieldName string, fv *validator.FieldValidator, g *protogen.GeneratedFile) {
	if fv.Regex != nil || fv.UuidVer != nil {
		if fv.UuidVer != nil {
			uuid, err := getUUIDRegex(fv.UuidVer)
			if err != nil {
				fmt.Fprintf(os.Stderr, "WARNING: field %v.%v error %s.\n", ccTypeName, fieldName, err)
			} else {
				fv.Regex = &uuid
			}
		}

		g.P(`if !`, p.regexName(ccTypeName, fieldName), `.MatchString(`, variableName, `) {`)
		g.P()
		errorStr := "be a string conforming to regex " + strconv.Quote(fv.GetRegex())
		p.generateErrorString(variableName, fieldName, errorStr, fv, g)
		g.P()
		g.P(`}`)
	}
	if fv.StringNotEmpty != nil && fv.GetStringNotEmpty() {
		g.P(`if `, variableName, ` == "" {`)
		g.P()
		errorStr := "not be an empty string"
		p.generateErrorString(variableName, fieldName, errorStr, fv, g)
		g.P()
		g.P(`}`)
	}
	p.generateLengthValidator(variableName, fieldName, fv, g)
}

func (p *plugin) generateRepeatedCountValidator(variableName string, fieldName string, fv *validator.FieldValidator, g *protogen.GeneratedFile) {
	if fv == nil {
		return
	}
	if fv.RepeatedCountMin != nil {
		compareStr := fmt.Sprint(`if len(`, variableName, `) < `, fv.GetRepeatedCountMin(), ` {`)
		g.P(compareStr)
		g.P()
		errorStr := fmt.Sprint(`contain at least `, fv.GetRepeatedCountMin(), ` elements`)
		p.generateErrorString(variableName, fieldName, errorStr, fv, g)
		g.P()
		g.P(`}`)
	}
	if fv.RepeatedCountMax != nil {
		compareStr := fmt.Sprint(`if len(`, variableName, `) > `, fv.GetRepeatedCountMax(), ` {`)
		g.P(compareStr)
		g.P()
		errorStr := fmt.Sprint(`contain at most `, fv.GetRepeatedCountMax(), ` elements`)
		p.generateErrorString(variableName, fieldName, errorStr, fv, g)
		g.P()
		g.P(`}`)
	}
}

func (p *plugin) generateErrorString(variableName string, fieldName string, specificError string, fv *validator.FieldValidator, g *protogen.GeneratedFile) {
	if fv.GetCustomError() == "" {
		g.P(`return `, p.validatorPkg.Ident("FieldError"), `("`, fieldName, `",`, p.fmtPkg.Ident("Errorf"), "(`value '%v' must ", specificError, "`", `, `, variableName, `))`)
	} else {
		g.P(`return `, p.errorsPkg.Ident("New"), "(`", fv.GetCustomError(), "`)")
	}
}

func (p *plugin) validatorWithMessageExists(fv *validator.FieldValidator) bool {
	return fv != nil && fv.MsgExists != nil && *(fv.MsgExists)
}

func (p *plugin) validatorWithNonRepeatedConstraint(fv *validator.FieldValidator) bool {
	if fv == nil {
		return false
	}

	// Need to use reflection in order to be future-proof for new types of constraints.
	v := reflect.ValueOf(*fv)
	for i := 0; i < v.NumField(); i++ {
		fieldName := v.Type().Field(i).Name

		// All known validators will have a pointer type and we should skip any fields
		// that are not pointers (i.e unknown fields, etc) as well as 'nil' pointers that
		// don't lead to anything.
		if v.Type().Field(i).Type.Kind() != reflect.Ptr || v.Field(i).IsNil() {
			continue
		}

		// Identify non-repeated constraints based on their name.
		if fieldName != "RepeatedCountMin" && fieldName != "RepeatedCountMax" {
			return true
		}
	}
	return false
}

func (p *plugin) regexName(ccTypeName string, fieldName string) string {
	return "regex" + ccTypeName + fieldName
}
