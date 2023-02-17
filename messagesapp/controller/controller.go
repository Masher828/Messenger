package controller

import (
	"encoding/json"
	"github.com/Masher828/MessengerBackend/common-shared-package/system"
	"github.com/Masher828/MessengerBackend/messagesapp/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller struct {
}

type response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}

func (controller *Controller) SendMessage(c *gin.Context, log *zap.SugaredLogger) ([]byte, error) {

	userContext := system.GetUserContextFromGinContext(c)
	if userContext == nil {
		err := system.ErrUnauthorizedAccess
		log.Errorln(err)
		return nil, err
	}

	data := models.Message{}
	err := c.Bind(&data)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	data.SenderId = userContext.UserId

	conversation := models.Conversation{}
	conversation.Id = data.ConversationId
	err = conversation.SendMessage(log, &data)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	messages, err := conversation.GetMessages(log, 0, system.MessagesLimit)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	resp := response{Success: true, Data: messages}
	return json.Marshal(resp)
}

func (controller *Controller) GetMessagesForConversation(c *gin.Context, log *zap.SugaredLogger) ([]byte, error) {

	userContext := system.GetUserContextFromGinContext(c)
	if userContext == nil {
		err := system.ErrUnauthorizedAccess
		log.Errorln(err)
		return nil, err
	}

	conversationId := c.Param("conversationId")
	if len(conversationId) == 0 {
		err := system.ErrInvalidConversationId
		log.Errorln(err)
		return nil, err
	}

	offset, limit := system.GetOffsetAndLimitFromContext(c, system.MessagesLimit)

	conversation := models.Conversation{Id: conversationId}
	messages, err := conversation.GetMessages(log, offset, limit)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	resp := response{Success: true, Data: messages}
	return json.Marshal(resp)
}

func (controller *Controller) GetConversations(c *gin.Context, log *zap.SugaredLogger) ([]byte, error) {

	userContext := system.GetUserContextFromGinContext(c)
	if userContext == nil {
		err := system.ErrUnauthorizedAccess
		log.Errorln(err)
		return nil, err
	}

	offset, limit := system.GetOffsetAndLimitFromContext(c, system.ConversationLimit)

	conversationName := c.Query("name")

	conversation := models.Conversation{}
	conversations, err := conversation.SearchConversationByName(log, conversationName, userContext.UserId, offset, limit)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	resp := response{Success: true, Data: conversations}
	return json.Marshal(resp)
}

func (controller *Controller) GetConversationWithFriend(c *gin.Context, log *zap.SugaredLogger) ([]byte, error) {

	userContext := system.GetUserContextFromGinContext(c)
	if userContext == nil {
		err := system.ErrUnauthorizedAccess
		log.Errorln(err)
		return nil, err
	}

	friendId := c.Param("friendId")

	var participants = []string{userContext.UserId, friendId}

	var conv models.Conversation
	conversation, err := conv.GetConversationWithFriend(log, participants)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	resp := response{
		Success: true,
		Data:    conversation,
	}

	return json.Marshal(resp)
}

func (controller *Controller) GetConversationById(c *gin.Context, log *zap.SugaredLogger) ([]byte, error) {
	userContext := system.GetUserContextFromGinContext(c)
	if userContext == nil {
		err := system.ErrUnauthorizedAccess
		log.Errorln(err)
		return nil, err
	}

	conversationId := c.Param("conversationId")

	conversation := models.Conversation{Id: conversationId}

	err := conversation.SetConversationById(log, userContext.UserId)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	resp := response{
		Success: true,
		Data:    conversation,
	}

	return json.Marshal(resp)
}
