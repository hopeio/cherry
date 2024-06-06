package pickrouter

import (
	"github.com/hopeio/cherry/context/httpctx"
	"github.com/hopeio/pick"
	"net/http"
	"reflect"
)

var (
	Svcs            = make([]pick.Service[http.HandlerFunc], 0)
	HttpContextType = reflect.TypeOf((*httpctx.Context)(nil))
)

func RegisterService(svc ...pick.Service[http.HandlerFunc]) {
	Svcs = append(Svcs, svc...)
}
