package system

import (
	mongo_common_repo "github.com/Masher828/MessengerBackend/common-shared-package/mongo-common-repo"
	"go.uber.org/zap"
)

type Controller struct {
}

type UserContext struct {
	UserId      string
	Name        string
	AccessToken string
}

type UserProfile struct {
	Id        string `json:"id" bson:"_id"`
	FirstName string `json:"firstName" binding:"required,min=2,max=200" bson:"firstName"`
	LastName  string `json:"lastName,omitempty" binding:"max=200" bson:"lastName,omitempty"`
	EmailId   string `json:"emailId" binding:"required,email,min=5,max=200" bson:"emailId"`
	Phone     string `json:"phone,omitempty" bson:"phone,omitempty"`
	Status    string `json:"status,omitempty" bson:"status,omitempty"`
	Gender    string `json:"gender" binding:"max=20" bson:"gender"`
	UpdatedOn int64  `json:"updatedOn" bson:"updatedOn"`
	CreatedOn int64  `json:"createdOn" bson:"createdOn"`
	LastLogin int64  `json:"lastLoginOn" bson:"lastLoginOn"`
}

func (profile *UserProfile) GetUserByFilter(log *zap.SugaredLogger, filter map[string]interface{}) ([]*UserProfile, error) {

	var users []*UserProfile
	err := mongo_common_repo.GetDocumentsWithFilter(log, CollectionNameUser, filter, 0, 0, &users)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	return users, nil
}