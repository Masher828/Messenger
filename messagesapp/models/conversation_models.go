package models

import (
	"github.com/Masher828/MessengerBackend/common-packages/constants"
	"github.com/Masher828/MessengerBackend/common-packages/system"
)

type CreateConversationRequest struct {
	Conversation Conversation `json:"conversation"`
	Users        []int64      `json:"users"`
}

type Conversation struct {
	Id            string `json:"id" bson:"_id"`
	Name          string `json:"name" bson:"name"` //it will have value only in case of type as group
	Type          string `json:"type" bson:"type"`
	Description   string `json:"description" bson:"description"`
	RecentMessage string `json:"recentMessage" bson:"recentMessage"`
	MembersCount  string `json:"membersCount" bson:"membersCount"`
	Icon          string `json:"icon" bson:"icon"`
	CreatedBy     string `json:"createdBy" bson:"createdBy"`
	CreatedOn     int64  `json:"createdOn" bson:"createdOn"`
	UpdatedOn     int64  `json:"updatedOn" bson:"updatedOn"`
}

type UserConversation struct {
	Id             string `json:"id" bson:"_id"`
	UserId         int64  `json:"userId" bson:"userId"`
	ConversationId string `json:"conversationId" bson:"conversationId"`
	IsArchived     bool   `json:"isArchived" bson:"isArchived"`
	IsMuted        bool   `json:"isMuted" bson:"isMuted"`
}

func (conversation *Conversation) IsValid() (bool, error) {

	if conversation.Type != constants.ConversationTypeGroup || conversation.Type != constants.ConversationTypePersonal {
		return false, system.InavlidGroupType
	}

	if conversation.Type == constants.ConversationTypeGroup && conversation.Name == "" {
		return false, system.InavlidGroupName
	}

	return true, nil
}
