package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/Masher828/MessengerBackend/authapp/models"
	"github.com/Masher828/MessengerBackend/authapp/services"
	"github.com/Masher828/MessengerBackend/common-packages/system"
	"github.com/sirupsen/logrus"
	"github.com/zenazn/goji/web"
)

type Controller struct {
	system.Controller
}

func (controller *Controller) UserSignup(c web.C, w http.ResponseWriter, r *http.Request, log *logrus.Entry) ([]byte, error) {

	user := models.UserModel{}
	response := map[string]string{"status": "ok"}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Errorln(err.Error())
		return []byte{}, err
	}
	err = services.UserSignup(&user, log)
	if err != nil {
		log.Errorln(err)
		return []byte{}, err
	}

	return json.Marshal(response)
}

func (controller *Controller) UserSignin(c web.C, w http.ResponseWriter, r *http.Request, log *logrus.Entry) ([]byte, error) {
	user := models.UserLoginModel{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	userContext, err := services.UserSignIn(&user, log)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	response := map[string]interface{}{"status": "ok", "userContext": userContext}

	return json.Marshal(response)

}
