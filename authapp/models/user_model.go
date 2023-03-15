package models

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	mongocommonrepo "github.com/Masher828/MessengerBackend/common-shared-package/mongo-common-repo"
	"github.com/Masher828/MessengerBackend/common-shared-package/system"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"go.uber.org/zap"
)

type RequestUser struct {
	EmailId  string `json:"emailId,omitempty"`
	Password string `json:"password,omitempty"`
}

type User struct {
	Id                     string `json:"id" bson:"_id"`
	Name                   string `json:"name,omitempty" bson:"name,omitempty"`
	EmailId                string `json:"emailId,omitempty" binding:"required,email,min=5,max=200" bson:"emailId,omitempty"`
	Phone                  string `json:"phone,omitempty" bson:"phone,omitempty"`
	Status                 string `json:"status,omitempty" bson:"status,omitempty"`
	Gender                 string `json:"gender,omitempty" binding:"max=20" bson:"gender,omitempty"`
	UpdatedOn              int64  `json:"updatedOn,omitempty" bson:"updatedOn,omitempty"`
	CreatedOn              int64  `json:"createdOn,omitempty" bson:"createdOn,omitempty"`
	LastLogin              int64  `json:"lastLoginOn,omitempty" bson:"lastLoginOn,omitempty"`
	Deleted                bool   `json:"deleted,omitempty" bson:"deleted,omitempty"`
	Password               string `json:"password,omitempty" binding:"required,alphanum,min=8,max=200" bson:"password,omitempty"`
	Salt                   []byte `json:"salt,omitempty" bson:"salt,omitempty"`
	InCorrectPasswordCount int    `json:"inCorrectPasswordCount,omitempty" bson:"inCorrectPasswordCount,omitempty"`
	IsLocked               bool   `json:"isLocked,omitempty" bson:"isLocked,omitempty"`
	AccessToken            string `json:"-"`

	// To reset password
	ResetPasswordToken string `json:"resetPasswordToken,omitempty"`
}

func (user *User) Insert(log *zap.SugaredLogger) error {

	user.EmailId = strings.ToLower(user.EmailId)
	emailExists, err := user.IsEmailUnique(log)
	if err != nil {
		log.Errorln(err)
		return err
	}

	if !emailExists {
		err = system.ErrEmailAlreadyExists
		return err
	}
	user.Id = uuid.NewString()
	user.UpdatedOn = system.NowInUTCMicro()
	user.CreatedOn = user.UpdatedOn
	err = user.encryptedPassword(log)
	if err != nil {
		log.Errorln(err)
		return err
	}

	err = mongocommonrepo.InsertDocument(log, system.CollectionNameUser, user)
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

	dataToBeUpdated["updatedOn"] = system.NowInUTCMicro()
	err := mongocommonrepo.UpdateDocumentById(log, system.CollectionNameUser, user.Id, dataToBeUpdated)
	if err != nil {
		log.Errorln(err)
		return err
	}

	return nil
}

func (user *User) IsEmailUnique(log *zap.SugaredLogger) (bool, error) {
	filter := map[string]interface{}{"emailId": user.EmailId}
	count, err := mongocommonrepo.GetDocumentCountsByFilter(log, system.CollectionNameUser, filter)
	if err != nil {
		log.Errorln(err)
		return false, err
	}

	return count == 0, nil
}

func (user *User) SetUserByEmail(log *zap.SugaredLogger, email string) error {

	email = strings.ToLower(email)
	filter := map[string]interface{}{"emailId": email}
	err := mongocommonrepo.GetSingleDocumentByFilter(log, system.CollectionNameUser, filter, &user)
	if err != nil {
		log.Errorln(err)
		return err
	}

	return nil
}

func (user *User) SetUserById(log *zap.SugaredLogger, userId string) error {

	selectedFields := map[string]interface{}{"name": 1, "emailId": 1, "createdOn": 1, "updatedOn": 1, "phone": 1, "status": 1, "gender": 1}
	err := mongocommonrepo.GetSelectedFieldsDocumentsById(log, system.CollectionNameUser, userId, selectedFields, user)
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
	user.InCorrectPasswordCount += 1

	dataToBeUpdated := map[string]interface{}{"inCorrectPasswordCount": user.InCorrectPasswordCount}

	if user.InCorrectPasswordCount >= system.MaxPasswordRetries {
		dataToBeUpdated["isLocked"] = true
	}

	err := user.UpdateWithMap(log, dataToBeUpdated)
	if err != nil {
		log.Errorln(err)
		return
	}
}

func (user *User) GetUserContextDetails() *system.UserContext {
	userContext := system.UserContext{}

	userContext.UserId = user.Id
	userContext.AccessToken = user.AccessToken
	userContext.Name = user.Name

	return &userContext
}

func (user *User) AddAccessTokenToUser(log *zap.SugaredLogger) error {
	db := system.MessengerContext.Redis

	key := fmt.Sprintf(system.AccessTokenToUser, user.AccessToken)
	data, _ := json.Marshal(user.GetUserContextDetails())
	err := db.Set(context.TODO(), key, data, system.DefaultAccessTokenExpiry).Err()
	if err != nil {
		log.Errorln(err)
		return err
	}
	return nil
}

func (user *User) encryptedPassword(log *zap.SugaredLogger) error {
	var err error = nil
	user.Salt, user.Password, err = system.GetHashedPassword(user.Password)
	if err != nil {
		log.Errorln(err)
		return err
	}
	return nil
}

func (user *User) SignIn(log *zap.SugaredLogger, emailId, password string) (*system.UserContext, error) {

	emailId = strings.ToLower(emailId)
	err := user.SetUserByEmail(log, emailId)
	if err != nil {
		if err.Error() == system.ErrNoMongoDocument.Error() {
			err = system.ErrInvalidCredentials
		}
		log.Errorln(err)
		return nil, err
	}

	if user.IsLocked {
		return nil, system.ErrUserIsLockedOut
	}

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

	err = user.AddAccessTokenToUser(log)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	return user.GetUserContextDetails(), nil
}

func (user *User) validateResetPasswordToken() bool {
	db := system.MessengerContext.Redis

	key := fmt.Sprintf(system.ResetPasswordTokenKey, user.ResetPasswordToken)
	redisUserId := db.Get(context.TODO(), key).String()

	return user.Id == redisUserId
}

func (user *User) UpdatePassword(log *zap.SugaredLogger) error {

	if len(user.ResetPasswordToken) > 0 && !user.validateResetPasswordToken() {
		return system.ErrInvalidPasswordToken
	}
	err := user.encryptedPassword(log)
	if err != nil {
		log.Errorln(err)
		return err
	}

	dataToUpdate := map[string]interface{}{"salt": user.Salt, "password": user.Password, "inCorrectPasswordCount": 0, "isLocked": false}
	err = user.UpdateWithMap(log, dataToUpdate)
	if err != nil {
		log.Errorln(err)
		return err
	}

	return nil
}

func (user *User) SearchUsers(log *zap.SugaredLogger, searchQuery string, userId string, offset, limit int64) ([]*User, error) {

	searchQuery = "^" + regexp.QuoteMeta(searchQuery) + ".*"

	nameFilter := map[string]interface{}{"name": bson.M{"$regex": bsonx.Regex(searchQuery, "i")}}
	emailFilter := map[string]interface{}{"emailId": bson.M{"$regex": bsonx.Regex(searchQuery, "i")}}
	filter := map[string]interface{}{"$or": []map[string]interface{}{nameFilter, emailFilter}, "_id": map[string]interface{}{"$ne": userId}}

	selectedFields := map[string]interface{}{"name": 1, "emailId": 1}

	users, err := user.GetSelectedFieldsWithFilter(log, selectedFields, filter, offset, limit)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	return users, nil
}

func (user *User) GetSelectedFieldsWithFilter(log *zap.SugaredLogger, selectedFields, filter map[string]interface{}, offset, limit int64) ([]*User, error) {
	var users []*User
	err := mongocommonrepo.GetSelectedFieldsDocumentsWithFilter(log, system.CollectionNameUser, selectedFields, filter, offset, limit, &users)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	return users, nil
}

func (user *User) ResetPassword(log *zap.SugaredLogger) (string, error) {
	//TODO handle with publisher to send otp
	resetPasswordToken := uuid.NewString()
	fmt.Println("Reset Password Token : ", resetPasswordToken)

	dataToUpdate := map[string]interface{}{"password": "", "salt": ""}
	err := user.UpdateWithMap(log, dataToUpdate)
	if err != nil {
		log.Errorln(err)
		return "", err
	}

	db := system.MessengerContext.Redis
	key := fmt.Sprintf(system.ResetPasswordTokenKey, user.Id)
	err = db.Set(context.TODO(), key, true, system.ResetPasswordTokenExpiry).Err()
	if err != nil {
		log.Errorln(err)
		return "", err
	}

	return resetPasswordToken, nil
}

func (user *User) UpdateStatus(log *zap.SugaredLogger, status string) error {

	dataToUpdate := map[string]interface{}{"status": status}

	err := user.UpdateWithMap(log, dataToUpdate)
	if err != nil {
		log.Errorln(err)
		return err
	}
	return nil
}
