package controllers

import (
	"net/http"

	"github.com/Masher828/MessengerBackend/common-packages/system"
	"github.com/zenazn/goji/web"
)

type Controller struct {
	system.Controller
}

func (controller *Controller) CreateUser(c web.C, w http.ResponseWriter, r *http.Request) ([]byte, error) {
	return []byte("Hello"), nil
}