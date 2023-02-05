package controller

import (
	"encoding/json"
	"github.com/Masher828/MessengerBackend/common-shared-package/system"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller struct {
}

func (controller *Controller) SignIn(c *gin.Context, log *zap.Logger) ([]byte, error) {
	return json.Marshal(system.Response{
		Success: true,
	})
}
