package controller

import (
	"encoding/json"
	"github.com/Masher828/MessengerBackend/authapp/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller struct {
}

type response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}

func (controller *Controller) SignIn(c *gin.Context, log *zap.SugaredLogger) ([]byte, error) {

	data := models.RequestUser{}

	err := c.Bind(&data)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	user := models.User{}
	userDetails, err := user.SignIn(log, data.EmailId, data.Password)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	resp := response{Success: true, Data: userDetails}
	return json.Marshal(resp)
}
