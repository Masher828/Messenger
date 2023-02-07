package controller

import (
	"encoding/json"
	"errors"
	"github.com/Masher828/MessengerBackend/authapp/models"
	"github.com/Masher828/MessengerBackend/common-shared-package/system"
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

func (controller *Controller) SignUp(c *gin.Context, log *zap.SugaredLogger) ([]byte, error) {
	var user models.User

	err := c.BindJSON(&user)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	err = user.Insert(log)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	resp := response{
		Success: true,
	}

	return json.Marshal(resp)

}

func (controller *Controller) ResetPassword(c *gin.Context, log *zap.SugaredLogger) ([]byte, error) {

	data := models.RequestUser{}

	err := c.Bind(&data)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	user := models.User{}
	err = user.SetUserByEmail(log, data.EmailId)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	resetPasswordToken, err := user.ResetPassword(log)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	resp := response{
		Success: true,
		Data:    resetPasswordToken,
	}

	return json.Marshal(resp)
}

func (controller *Controller) UpdatePassword(c *gin.Context, log *zap.SugaredLogger) ([]byte, error) {

	var user models.User
	err := c.BindJSON(&user)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	userContext := system.GetUserContextFromGinContext(c)

	if userContext == nil && len(user.ResetPasswordToken) == 0 {
		err = errors.New("invalid payload to update password")
		log.Errorln(err)
		return nil, err
	}

	err = user.UpdatePassword(log)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	resp := response{
		Success: true,
	}

	return json.Marshal(resp)
}

func (controller *Controller) UpdateStatus(c *gin.Context, log *zap.SugaredLogger) ([]byte, error) {
	userContext := system.GetUserContextFromGinContext(c)
	if userContext == nil {
		err := system.ErrUnauthorizedAccess
		log.Errorln(err)
		return nil, err
	}

	user := models.User{Id: userContext.UserId}

	err := user.UpdateStatus(log, c.Param("status"))
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	resp := response{
		Success: true,
	}

	return json.Marshal(resp)
}

func (controller *Controller) UpdateProfile(c *gin.Context, log *zap.SugaredLogger) ([]byte, error) {

	userContext := system.GetUserContextFromGinContext(c)
	if userContext == nil {
		err := system.ErrUnauthorizedAccess
		log.Errorln(err)
		return nil, err
	}

	updatedUser := models.User{}
	err := c.ShouldBind(&updatedUser)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	existingUser := models.User{}

	err = existingUser.SetUserById(log, userContext.UserId)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	err = existingUser.Update(log, &updatedUser)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	resp := response{Success: true}

	return json.Marshal(resp)
}
