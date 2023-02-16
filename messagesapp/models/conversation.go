package models

import (
	"errors"
	mongo_common_repo "github.com/Masher828/MessengerBackend/common-shared-package/mongo-common-repo"
	"github.com/Masher828/MessengerBackend/common-shared-package/system"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
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
	RecentMessage    string   `json:"recentMessage,omitempty" bson:"recentMessage,omitempty"`
	ParticipantsName []string `json:"participantsName" bson:"participantsName"` //only for system.ConversationTypeOne2One
	Participants     []string `json:"participants" bson:"participants"`
	Type             string   `json:"conversationType" bson:"conversationType"`
	CreatedBy        string   `json:"createdBy" bson:"createdBy"`
	CreatedOn        int64    `json:"createdOn" bson:"createdOn"`
	UpdatedOn        int64    `json:"updatedOn" bson:"updatedOn"`
}

func (conversation *Conversation) ValidateIndividualConversation(log *zap.SugaredLogger) error {

	if len(conversation.ParticipantsName) == 0 && conversation.Type == system.ConversationTypeOne2One {
		err := system.ErrOne2OneConversationNoName
		log.Errorln(err)
		return err
	}

	var conv Conversation
	filter := map[string]interface{}{"participants": map[string]interface{}{"$all": conversation.Participants}, "conversationType": system.ConversationTypeOne2One}
	err := mongo_common_repo.GetSingleDocumentByFilter(log, system.CollectionNameConversation, filter, &conv)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Errorln(err)
		return err
	}

	if len(conv.Id) != 0 {
		conversation.Id = conv.Id
		err = system.ErrConversationAlreadyExist
		return err
	}

	return nil
}

func (conversation *Conversation) ValidateGroupConversation(log *zap.SugaredLogger) error {
	if len(conversation.Participants) <= 1 && conversation.Type == system.ConversationTypeGroup {
		err := system.ErrGroupConversationMinimumOneUser
		log.Errorln(err)
		return err
	}

	return errors.New("fake")
}

func (conversation *Conversation) SetConversationById(log *zap.SugaredLogger) error {
	err := mongo_common_repo.GetSingleDocumentByFilter(log, system.CollectionNameConversation, map[string]interface{}{"_id": conversation.Id}, &conversation)
	if err != nil {
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

	var users []*system.UserContext

	filter := map[string]interface{}{"_id": map[string]interface{}{"$in": conversation.Participants}}

	err := mongo_common_repo.GetDocumentsWithFilter(log, system.CollectionNameUser, filter, 0, 0, &users, -1)
	if err != nil {
		log.Errorln(err)
		return err
	}

	for _, user := range users {
		name := user.Name
		conversation.ParticipantsName = append(conversation.ParticipantsName, name)
	}

	err = conversation.ValidateIndividualConversation(log)
	if err != nil {
		log.Errorln(err)
		return err
	}
	return nil
}

func (conversation *Conversation) CreateGroup(log *zap.SugaredLogger) error {

	err := conversation.ValidateGroupConversation(log)
	if err != nil {
		log.Errorln(err)
		return err
	}

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
		conversation.Type = system.ConversationTypeOne2One
		conversation.Participants = []string{message.SenderId, message.ReceiverId}
		err := conversation.Create(log)
		if err != nil && err != system.ErrConversationAlreadyExist {
			log.Errorln(err)
			return err
		}
		message.ConversationId = conversation.Id
	} else {
		err := conversation.SetConversationById(log)
		if err != nil {
			log.Errorln(err)
			return err
		}

		if !conversation.IsParticipant(log, message.SenderId) {
			err := system.ErrNotMemberOfConversation
			log.Errorln(err)
			return err
		}
	}

	conversation.RecentMessage = message.Body
	err := message.Send(log)
	if err != nil {
		log.Errorln(err)
		return err
	}
	go func() {
		err = conversation.UpdateRecentMessage(log)
		if err != nil {
			log.Errorln(err)
		}
	}()
	err = conversation.SendNotificationToParticipants(log, message.SenderId)
	if err != nil {
		log.Errorln(err)
		return err
	}
	return nil
}

func (conversation *Conversation) GetMessagesWithFriend(log *zap.SugaredLogger, userIds []string) ([]*Message, error) {
	filter := map[string]interface{}{"participants": map[string]interface{}{"$all": userIds}, "type": system.ConversationTypeOne2One}

	conversations, err := conversation.GetConversationsWithFilter(log, filter, 0, 1)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	var messages []*Message
	if len(conversations) != 0 {
		conversation = conversations[0]
		messages, err = conversation.GetMessages(log, 0, system.MessagesLimit)
		if err != nil {
			log.Errorln(err)
			return nil, err
		}
	}

	return messages, err
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

func (conversation *Conversation) ValidateConversation(log *zap.SugaredLogger) error {

	filter := map[string]interface{}{"_id": map[string]interface{}{"$in": conversation.Participants}}
	count, err := mongo_common_repo.GetDocumentCountsByFilter(log, system.CollectionNameUser, filter)
	if err != nil {
		log.Errorln(err)
		return err
	}

	if count != int64(len(conversation.Participants)) {
		err = system.ErrInvalidConversationParticipants
		log.Errorln(err)
		return err
	}

	return nil
}

func (conversation *Conversation) Create(log *zap.SugaredLogger) error {
	err := conversation.ValidateConversation(log)
	if err != nil {
		log.Errorln(err)
		return err
	}

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

func (conversation *Conversation) UpdateRecentMessage(log *zap.SugaredLogger) error {
	filter := map[string]interface{}{"_id": conversation.Id}
	dataToUpdate := map[string]interface{}{"recentMessage": conversation.RecentMessage}
	err := mongo_common_repo.UpdateDocumentByFilter(log, system.CollectionNameConversation, filter, dataToUpdate)
	if err != nil {
		log.Errorln(err)
		return err
	}

	return nil
}
func (conversation *Conversation) GetConversationsWithFilter(log *zap.SugaredLogger, filter map[string]interface{}, offset, limit int64) ([]*Conversation, error) {
	var conversations []*Conversation

	err := mongo_common_repo.GetDocumentsWithFilter(log, system.CollectionNameConversation, filter, offset, limit, &conversations, -1)
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
	if len(searchQuery) > 0 {
		filter["$or"] = []map[string]interface{}{{"participantsName": searchQuery}, {"name": searchQuery}}
	}

	conversations, err := conversation.GetConversationsWithFilter(log, filter, offset, limit)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	return conversations, nil
}
