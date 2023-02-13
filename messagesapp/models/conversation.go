package models

import (
	mongo_common_repo "github.com/Masher828/MessengerBackend/common-shared-package/mongo-common-repo"
	"github.com/Masher828/MessengerBackend/common-shared-package/system"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ConversationUnreadCount struct {
	UserId         string `json:"userId" bson:"userId"`
	ConversationId string `json:"conversationId" bson:"conversationId"`
	Count          int64  `json:"count" bson:"count"`
}

type Conversation struct {
	Id               string   `json:"id" bson:"_id"`
	Name             string   `json:"name" bson:"name"`
	ParticipantsName []string `json:"participantsName" bson:"participantsName"` //only for system.ConversationTypeOne2One
	Participants     []string `json:"participants" bson:"participants"`
	Type             string   `json:"conversationType" bson:"conversationType"`
	CreatedBy        string   `json:"createdBy" bson:"createdBy"`
	CreatedOn        int64    `json:"createdOn" bson:"createdOn"`
	UpdatedOn        int64    `json:"updatedOn" bson:"updatedOn"`
}

func (conversation *Conversation) Validate(log *zap.SugaredLogger) error {
	if len(conversation.Participants) <= 1 && conversation.Type == system.ConversationTypeGroup {
		err := system.ErrGroupConversationMinimumOneUser
		log.Errorln(err)
		return err
	}

	if len(conversation.ParticipantsName) == 0 && conversation.Type == system.ConversationTypeOne2One {
		err := system.ErrOne2OneConversationNoName
		log.Errorln(err)
		return err
	}
	return nil
}

func (conversation *Conversation) GetConversations(log *zap.SugaredLogger, userId string, offset, limit int64) ([]*Conversation, error) {

	filter := map[string]interface{}{"participants": userId}
	conversations, err := conversation.GetConversationsWithFilter(log, filter, offset, limit)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	return conversations, nil

}

func (conversation *Conversation) CreateIndividualChat(log *zap.SugaredLogger) error {

	user := system.UserProfile{}
	filter := map[string]interface{}{"_id": map[string]interface{}{"$in": conversation.Participants}}
	users, err := user.GetUserByFilter(log, filter)
	if err != nil {
		log.Errorln(err)
		return err
	}

	for _, user := range users {
		name := user.FirstName
		if len(user.LastName) > 0 {
			name += " " + user.LastName
		}

		conversation.ParticipantsName = append(conversation.ParticipantsName, name)
	}

	return nil
}

func (conversation *Conversation) CreateGroup(log *zap.SugaredLogger) error {

	return nil
}

func (conversation *Conversation) IsParticipant(log *zap.SugaredLogger, participantId string) bool {
	return system.ContainsString(conversation.Participants, participantId)
}

func (conversation *Conversation) DeleteMessage(log *zap.SugaredLogger, messageId, userId string) error {
	if !conversation.IsParticipant(log, userId) {
		err := system.ErrNotMemberOfConversation
		log.Errorln(err)
		return err
	}

	message := Message{Id: messageId, SenderId: userId}
	err := message.Delete(log)
	if err != nil {
		log.Errorln(err)
		return err
	}

	return nil
}

func (conversation *Conversation) SendMessage(log *zap.SugaredLogger, message *Message) error {

	if len(conversation.Id) == 0 {
		conversation.CreatedBy = message.SenderId

		err := conversation.Create(log)
		if err != nil {
			log.Errorln(err)
			return err
		}
	} else if !conversation.IsParticipant(log, message.SenderId) {
		err := system.ErrNotMemberOfConversation
		log.Errorln(err)
		return err
	}

	err := message.Send(log)
	if err != nil {
		log.Errorln(err)
		return err
	}

	err = conversation.SendNotificationToParticipants(log, message.SenderId)
	if err != nil {
		log.Errorln(err)
		return err
	}
	return nil
}

func (conversation *Conversation) GetMessages(log *zap.SugaredLogger, offset, limit int64) ([]*Message, error) {
	message := Message{ConversationId: conversation.Id}
	messages, err := message.Get(log, offset, limit)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	return messages, nil
}

func (conversation *Conversation) SendNotificationToParticipants(log *zap.SugaredLogger, senderId string) error {

	//TODO : implement pub/sub and remove the sender and move this to messages
	filter := map[string]interface{}{"conversationId": conversation.Id, "userId": map[string]interface{}{"$in": conversation.Participants}}
	updateQuery := map[string]interface{}{"$inc": map[string]interface{}{"count": 1}}
	err := mongo_common_repo.UpsertDocumentCustomQuery(log, system.CollectionNameConversationUnread, filter, updateQuery)
	if err != nil {
		log.Errorln(err)
		return err
	}

	return nil
}

func (conversation *Conversation) Create(log *zap.SugaredLogger) error {
	var err error = nil

	if !system.ContainsString(conversation.Participants, conversation.CreatedBy) {
		conversation.Participants = append(conversation.Participants, conversation.CreatedBy)
	}

	conversation.Id = uuid.NewString()
	conversation.CreatedOn = system.NowInUTCMicro()
	conversation.UpdatedOn = conversation.CreatedOn

	if conversation.Type == system.ConversationTypeOne2One {
		err = conversation.CreateIndividualChat(log)
	} else if conversation.Type == system.ConversationTypeGroup {
		err = conversation.CreateGroup(log)
	} else {
		err = system.ErrInvalidConversationType
	}
	if err != nil {
		log.Errorln(err)
		return err
	}

	err = mongo_common_repo.InsertDocument(log, system.CollectionNameConversation, conversation)
	if err != nil {
		log.Errorln(err)
		return err
	}

	go func() {
		err = conversation.SendNotificationToParticipants(log, conversation.CreatedBy)
		if err != nil {
			log.Errorln(err)
		}
	}()

	return nil
}

func (conversation *Conversation) GetConversationsWithFilter(log *zap.SugaredLogger, filter map[string]interface{}, offset, limit int64) ([]*Conversation, error) {
	var conversations []*Conversation

	err := mongo_common_repo.GetDocumentsWithFilter(log, system.CollectionNameConversation, filter, offset, limit, conversations)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	return conversations, nil
}

func (conversation *Conversation) RemoveUser(log *zap.SugaredLogger) {

}

func (conversation *Conversation) SearchConversationByName(log *zap.SugaredLogger, searchQuery string, userId string, offset, limit int64) ([]*Conversation, error) {
	filter := map[string]interface{}{"participants": userId}
	filter["$or"] = []map[string]interface{}{{"participantsName": searchQuery}, {"name": searchQuery}}

	conversations, err := conversation.GetConversationsWithFilter(log, filter, offset, limit)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	return conversations, nil
}
