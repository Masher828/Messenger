package repository

import (
	"context"

	"github.com/Masher828/MessengerBackend/common-packages/constants"
	"github.com/Masher828/MessengerBackend/common-packages/system"
	"github.com/Masher828/MessengerBackend/messagesapp/models"
	"github.com/sirupsen/logrus"
)

func CreateConversation(conversation *models.Conversation, log *logrus.Entry) error {

	client := system.SocialContext.MongoClient
	db := client.Database(constants.DatabaseSocialDB).Collection(constants.ConversationCollection)

	_, err := db.InsertOne(context.TODO(), conversation)
	if err != nil {
		log.Errorln(err)
	}
	return err
}

func AddUserToConversation(userConversations []interface{}, log *logrus.Entry) error {

	client := system.SocialContext.MongoClient
	db := client.Database(constants.DatabaseSocialDB).Collection(constants.UserConversationCollection)

	_, err := db.InsertMany(context.TODO(), userConversations)
	if err != nil {
		log.Errorln(err)
	}
	return err
}
