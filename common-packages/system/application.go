package system

import (
	"net/http"
	"reflect"

	"github.com/zenazn/goji/web"
)

type Controller struct {
}
type Application struct {
}

func (application *Application) Route(controller interface{}, route string) interface{} {
	fn := func(c web.C, w http.ResponseWriter, r *http.Request) {
		methodInterface := reflect.ValueOf(controller).MethodByName(route).Interface()
		method := methodInterface.(func(c web.C, w http.ResponseWriter, r *http.Request) ([]byte, error))
		response, err := method(c, w, r)
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write([]byte(response))
		}
	}
	return fn
}
