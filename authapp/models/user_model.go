package models

import (
	"context"
	"fmt"
	mongo_common_repo "github.com/Masher828/MessengerBackend/common-shared-package/mongo-common-repo"
	"github.com/Masher828/MessengerBackend/common-shared-package/system"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type RequestUser struct {
	EmailId  string `json:"emailId" bson:"emailId"`
	Password string `json:"password" bson:"password"`
}

type User struct {
	Id                     string `json:"id" bson:"_id"`
	FirstName              string `json:"firstName" binding:"required,min=2,max=200" bson:"firstName"`
	LastName               string `json:"lastName,omitempty" binding:"max=200" bson:"lastName,omitempty"`
	EmailId                string `json:"emailId" binding:"required,email,min=5,max=200" bson:"emailId"`
	Phone                  string `json:"phone,omitempty" bson:"phone,omitempty"`
	Deleted                bool   `json:"deleted,omitempty" bson:"deleted,omitempty"`
	Password               string `json:"password" binding:"required,alphanum,min=8,max=200" bson:"password"`
	Salt                   []byte `json:"salt" bson:"salt"`
	InCorrectPasswordCount int    `json:"inCorrectPasswordCount,omitempty" bson:"inCorrectPasswordCount,omitempty"`
	IsLocked               bool   `json:"isLocked" bson:"isLocked"`
	Gender                 string `json:"gender" binding:"max=20" bson:"gender"`
	UpdatedOn              int64  `json:"updatedOn" bson:"updatedOn"`
	CreatedOn              int64  `json:"createdOn" bson:"createdOn"`
	LastLogin              int64  `json:"lastLoginOn" bson:"lastLoginOn"`
	AccessToken            string `json:"accessToken" bson:"accessToken"`
}

func (user *User) Insert(log *zap.SugaredLogger) error {

	user.Id = uuid.NewString()
	user.UpdatedOn = system.NowInUTCMicro()
	user.CreatedOn = user.UpdatedOn
	err := user.getEncryptedPassword(log)
	if err != nil {
		log.Errorln(err)
		return err
	}

	err = mongo_common_repo.InsertDocument(log, system.UserCollectionName, user)
	if err != nil {
		log.Errorln(err)
		return err
	}

	return nil
}

func (user *User) getEncryptedPassword(log *zap.SugaredLogger) error {
	var err error = nil
	user.Salt, user.Password, err = system.GetHashedPassword(user.Password)
	if err != nil {
		log.Errorln(err)
		return err
	}
	return nil
}

func (user *User) Update(log *zap.SugaredLogger, updatedUser *User) error {

	dataToUpdate := system.GetDifferentFields(user, updatedUser)

	err := user.UpdateWithMap(log, dataToUpdate)
	if err != nil {
		log.Errorln(err)
		return err
	}

	return nil
}

func (user *User) UpdateWithMap(log *zap.SugaredLogger, dataToBeUpdated map[string]interface{}) error {
	err := mongo_common_repo.UpdateDocumentById(log, system.UserCollectionName, user.Id, dataToBeUpdated)
	if err != nil {
		log.Errorln(err)
		return err
	}

	return nil
}

func (user *User) IsEmailUnique(log *zap.SugaredLogger) (bool, error) {
	data := User{}
	err := mongo_common_repo.GetSingleDocumentByFilter(log, "d", map[string]interface{}{}, &data)
	if err != nil {
		log.Errorln(err)
		return false, err
	}

	return true, nil
}

func (user *User) SetUserByEmail(log *zap.SugaredLogger, email string) error {

	filter := map[string]interface{}{"emailId": email}
	err := mongo_common_repo.GetSingleDocumentByFilter(log, system.UserCollectionName, filter, &user)
	if err != nil {
		log.Errorln(err)
		return err
	}

	return nil
}

func (user *User) UpdateLastLogin(log *zap.SugaredLogger) {

	dataToUpdate := map[string]interface{}{"lastLogin": system.NowInUTCMicro()}

	err := user.UpdateWithMap(log, dataToUpdate)
	if err != nil {
		log.Errorln(err)
	}
}

func (user *User) invalidPasswordHandler(log *zap.SugaredLogger) {

}

func (user *User) GetUserContextDetails() *system.UserContext {
	userContext := system.UserContext{}

	userContext.UserId = user.Id
	userContext.AccessToken = user.AccessToken
	userContext.Name = user.FirstName + " " + user.LastName

	return &userContext
}
func (user *User) AddAccessTokenToUser(log *zap.SugaredLogger) error {
	db := system.MessengerContext.Redis

	key := fmt.Sprintf(system.AccessTokenToUser, user.AccessToken)
	err := db.Set(context.TODO(), key, user.GetUserContextDetails(), system.DefaultAccessTokenExpiry).Err()
	if err != nil {
		log.Errorln(err)
		return err
	}
	return nil
}

func (user *User) SignIn(log *zap.SugaredLogger, emailId, password string) (*system.UserContext, error) {

	hashedPassword := system.HashPassword(user.Salt, password)

	if hashedPassword != user.Password {
		go user.invalidPasswordHandler(log)
		fmt.Println("Invalid Password for emailId : ", emailId)
		return nil, system.ErrInvalidCredentials
	}

	dataToUpdate := map[string]interface{}{"lastLogin": system.NowInUTCMicro()}
	if user.InCorrectPasswordCount > 0 {
		dataToUpdate["inCorrectPasswordCount"] = 0
	}

	go func() {
		err := user.UpdateWithMap(log, dataToUpdate)
		if err != nil {
			log.Errorln(err)
			return
		}
	}()

	user.AccessToken = uuid.NewString()

	err := user.AddAccessTokenToUser(log)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	return user.GetUserContextDetails(), nil
}
