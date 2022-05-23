package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Masher828/MessengerBackend/authapp/models"
	"github.com/Masher828/MessengerBackend/authapp/repository"
	"github.com/Masher828/MessengerBackend/common-packages/system"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func UserSignup(user *models.UserModel, log *logrus.Entry) error {

	if _, err := user.IsValid(); err != nil {
		log.Errorln(err)
		return err
	}

	userDetails, err := repository.GetUserByEmail(user.Email, log)
	if err == nil && userDetails.Id != 0 {
		return system.EmailAlreadyExists
	}

	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), int(time.Now().Month()))
	if err != nil {
		log.Errorln(err)
		return err
	}

	user.Id = repository.GetNextHibernateSequence()
	user.Password = string(passwordBytes)

	err = repository.InsertUserToDB(user, log)
	if err != nil {
		log.Errorln(err)
	}

	return err
}

func UserSignIn(user *models.UserLoginModel, log *logrus.Entry) (*models.UserModelContext, error) {
	if _, err := user.IsValid(); err != nil {
		log.Error(err)
		return nil, system.InvalidCredentialsErr
	}

	userDetails, err := repository.GetUserByEmail(user.Email, log)
	if err != nil || userDetails.Id == 0 {
		if err != nil {
			log.Errorln(err)
		}
		return nil, system.InvalidCredentialsErr
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDetails.Password), []byte(user.Password))

	if err != nil {
		log.Errorln(err)
		return nil, system.InvalidCredentialsErr
	}

	accessToken := uuid.NewString()
	redisDb := system.SocialContext.Redis

	userContext := userDetails.GetUserContext()
	userContext.AccessToken = accessToken
	data, _ := json.Marshal(userContext)
  
	redisDb.Set(context.TODO(), accessToken, data, 60*time.Minute)
	return userContext, nil
}

func GetAllUsers(log *logrus.Entry) ([]string, error) {
	// constants.SendMail()
	log.Errorln("hii")
	return repository.GetAllUsers(log)
}
