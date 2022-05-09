package services

import (
	"github.com/Masher828/MessengerBackend/authapp/models"
	"github.com/Masher828/MessengerBackend/authapp/repository"
	"github.com/sirupsen/logrus"
)

func CreateUser(user *models.UserModel, log *logrus.Entry) error {
	// user.Validate()
	// encrypt password
	user.Id = repository.GetNextHibernateSequence()
	err := repository.InsertUserToDB(user, log)
	if err != nil {
		log.Errorln(err)
	}
	return err

}
