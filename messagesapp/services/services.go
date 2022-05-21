package services

import (
	"github.com/Masher828/MessengerBackend/messagesapp/models"
	"github.com/sirupsen/logrus"
)

func CreateConversation(conversation *models.Conversation, log *logrus.Entry) error {

	if _, err := conversation.IsValid(); err != nil {
		log.Errorln(err)
		return err
	}

}
