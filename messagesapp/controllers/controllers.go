package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/Masher828/MessengerBackend/common-packages/system"
	"github.com/Masher828/MessengerBackend/messagesapp/models"
	"github.com/Masher828/MessengerBackend/messagesapp/services"
	"github.com/sirupsen/logrus"
	"github.com/zenazn/goji/web"
)

type Controller struct {
	system.Controller
}

func (controller *Controller) CreateConversation(c web.C, w http.ResponseWriter, r *http.Request, log *logrus.Entry) ([]byte, error) {
	var conversation models.CreateConversationRequest

	var response map[string]string

	err := json.NewDecoder(r.Body).Decode(&conversation)
	if err != nil {
		log.Errorln(err)
		return []byte{}, err
	}

	err = services.CreateConversation(&conversation.Conversation, conversation.Users, log)
	if err != nil {
		log.Errorln(err)
		return []byte{}, err
	}

	response["success"] = "ok"

	return json.Marshal(response)
}
