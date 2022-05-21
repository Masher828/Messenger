package repository

import (
	"github.com/Masher828/MessengerBackend/common-packages/system"
	"github.com/Masher828/MessengerBackend/messagesapp/models"
	"github.com/sirupsen/logrus"
)

func InsertConversation(conversation *models.Conversation, log *logrus.Entry) error {

	db := system.SocialContext.MongoClient

	db.
}
