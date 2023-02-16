package models

import (
	mongo_common_repo "github.com/Masher828/MessengerBackend/common-shared-package/mongo-common-repo"
	"github.com/Masher828/MessengerBackend/common-shared-package/system"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Message struct {
	Id             string `json:"id" bson:"_id"`
	ConversationId string `json:"conversationId" bson:"conversationId"`
	SenderId       string `json:"senderId" bson:"senderId"`
	Body           string `json:"body,omitempty" bson:"body,omitempty"`
	MessageType    string `json:"messageType" bson:"messageType"`
	CreatedOn      int64  `json:"createdOn" bson:"createdOn"`
	UpdatedOn      int64  `json:"updatedOn" bson:"updatedOn"`

	ReceiverId string `json:"receiverId,omitempty"` //only for the first message of conversation and will not be stored oin the db. Used to create conversation
}

func (message *Message) Get(log *zap.SugaredLogger, offset, limit int64) ([]*Message, error) {

	filter := map[string]interface{}{"conversationId": message.ConversationId}
	messagesList, err := message.getMessageWithFilter(log, filter, offset, limit)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	return messagesList, nil
}

func (message *Message) getMessageWithFilter(log *zap.SugaredLogger, filter map[string]interface{}, offset, limit int64) ([]*Message, error) {

	var messagesList []*Message

	err := mongo_common_repo.GetDocumentsWithFilter(log, system.CollectionNameMessages, filter, offset, limit, &messagesList, 1)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	return messagesList, nil
}

func (message *Message) Send(log *zap.SugaredLogger) error {

	err := message.insert(log)
	if err != nil {
		log.Errorln(err)
		return err
	}

	return nil

}

func (message *Message) Delete(log *zap.SugaredLogger) error {
	filter := map[string]interface{}{"_id": message.Id, "senderId": message.SenderId}
	err := mongo_common_repo.DeleteSingleDocumentByFilter(log, system.CollectionNameMessages, filter)
	if err != nil {
		log.Errorln(err)
		return err
	}

	return nil
}

func (message *Message) insert(log *zap.SugaredLogger) error {
	message.Id = uuid.NewString()
	message.CreatedOn = system.NowInUTCMicro()
	message.UpdatedOn = message.CreatedOn
	err := mongo_common_repo.InsertDocument(log, system.CollectionNameMessages, message)
	if err != nil {
		log.Errorln(err)
		return err
	}
	return nil
}
