package binding

import (
	"github.com/gofiber/fiber/v3"
	"github.com/hopeio/cherry/utils/net/http/binding"
	stringsi "github.com/hopeio/cherry/utils/strings"
	"net/http"
)

type Binding interface {
	Name() string

	Bind(fiber.Ctx, interface{}) error
}

type BindingBody interface {
	Binding
	BindBody([]byte, interface{}) error
}

var (
	JSON          = jsonBinding{}
	XML           = xmlBinding{}
	Form          = formBinding{}
	Query         = queryBinding{}
	FormPost      = formPostBinding{}
	FormMultipart = formMultipartBinding{}
	Uri           = uriBinding{}
	ProtoBuf      = protobufBinding{}
	MsgPack       = msgpackBinding{}
	YAML          = yamlBinding{}
	Header        = headerBinding{}
)

func Default(method string, contentType []byte) Binding {
	if method == http.MethodGet {
		return Query
	}

	switch stringsi.BytesToString(contentType) {
	case binding.MIMEJSON:
		return JSON
	case binding.MIMEXML, binding.MIMEXML2:
		return XML
	case binding.MIMEPROTOBUF:
		return ProtoBuf
	case binding.MIMEMSGPACK, binding.MIMEMSGPACK2:
		return MsgPack
	case binding.MIMEYAML:
		return YAML
	case binding.MIMEMultipartPOSTForm:
		return FormMultipart
	default: // case MIMEPOSTForm:
		return Form
	}
}
