package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/Masher828/MessengerBackend/authapp/models"
	"github.com/Masher828/MessengerBackend/common-packages/system"
	"github.com/sirupsen/logrus"
	"github.com/zenazn/goji/web"
)

type Controller struct {
	system.Controller
}

func (controler *Controller) CreateUser(c web.C, w http.ResponseWriter, r *http.Request, log *logrus.Entry) ([]byte, error) {

	user := models.UserModel{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Errorln(err.Error())
	}
	return []byte("helo"), nil
}
