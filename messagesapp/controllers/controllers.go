package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Masher828/MessengerBackend/common-packages/constants"
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

	response := make(map[string]string)

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

func (controller *Controller) GetConversation(c web.C, w http.ResponseWriter, r *http.Request, log *logrus.Entry) ([]byte, error) {

	id, err := strconv.ParseInt(c.URLParams["userId"], 10, 64)

	offsetArr := r.URL.Query()["offset"]

	offset := constants.DefaultConversationOffset

	if len(offsetArr) > 0 {
		offset, err = strconv.ParseInt(offsetArr[0], 10, 64)
		if err != nil {
			log.Errorln(err)
			offset = constants.DefaultConversationOffset
		}
	}

	limit := constants.DefaultConversationLimit

	limitArr := r.URL.Query()["limit"]

	if len(limitArr) > 0 {
		limit, err = strconv.ParseInt(limitArr[0], 10, 64)
		if err != nil {
			log.Errorln(err)
			limit = constants.DefaultConversationLimit
		}
	}

	conversations, err := services.GetuserConversation(id, offset, limit, log)
	if err != nil {
		log.Errorln(err)
		return []byte{}, err
	}

	return json.Marshal(conversations)

}
