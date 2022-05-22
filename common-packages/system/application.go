package system

import (
	"net/http"
	"reflect"

	"github.com/Masher828/MessengerBackend/common-packages/log"
	"github.com/sirupsen/logrus"
	"github.com/zenazn/goji/web"
)

type Controller struct {
}
type Application struct {
}

func (application *Application) Route(controller interface{}, route string) interface{} {
	fn := func(c web.C, w http.ResponseWriter, r *http.Request) {
		methodInterface := reflect.ValueOf(controller).MethodByName(route).Interface()
		logger := log.GetDefaultLogger(0, r.RequestURI, r.Method)
		method := methodInterface.(func(c web.C, w http.ResponseWriter, r *http.Request, logger *logrus.Entry) ([]byte, error))
		response, err := method(c, w, r, logger)
		if err != nil {
			if err.Error() == "EOF" {
				err = InvalidPayloadData
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		} else {
			w.Write([]byte(response))
		}
	}
	return fn
}
