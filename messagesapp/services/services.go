package services

import (
	"github.com/Masher828/MessengerBackend/common-packages/system"
	"github.com/Masher828/MessengerBackend/messagesapp/models"
	"github.com/Masher828/MessengerBackend/messagesapp/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func CreateConversation(conversation *models.Conversation, users []int64, log *logrus.Entry) error {

	if _, err := conversation.IsValid(); err != nil {
		log.Errorln(err)
		return err
	}

	conversation.Id = uuid.New().String()

	now := system.GetUTCTime().Unix()
	conversation.UpdatedOn = now
	conversation.UpdatedOn = now
	err := repository.CreateConversation(conversation, log)
	if err != nil {
		log.Errorln(err)
		return err
	}

	var userConversations []interface{}

	for _, userId := range users {
		userConversation := models.UserConversation{
			Id:             string(uuid.New().String()),
			ConversationId: conversation.Id,
			UserId:         userId,
			IsArchived:     false,
			IsMuted:        false,
			UpdatedOn:      now,
			CreatedOn:      now,
		}

		userConversations = append(userConversations, &userConversation)
	}

	err = repository.AddUserToConversation(userConversations, log)
	if err != nil {
		log.Errorln(err)
	}

	return err
}

func GetuserConversation(id int64, offset, limit int64, log *logrus.Entry) ([]models.ResponseUserConversation, error) {

	return repository.GetuserConversation(id, offset, limit, log)
}
