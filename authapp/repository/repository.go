package repository

import (
	"github.com/Masher828/MessengerBackend/authapp/models"
	"github.com/Masher828/MessengerBackend/common-packages/system"
	"github.com/sirupsen/logrus"
)

func GetNextHibernateSequence() int64 {
	query := "select nextval('hibernate_sequence')"
	db := system.SocialContext.PostgresDB
	var id int64
	db.QueryRow(query).Scan(&id)
	return id
}

func InsertUserToDB(user *models.UserModel, log *logrus.Entry) error {
	query := `INSERT INTO social_user (id, name, email, password, contact, country_code, country, 
		date_of_birth, last_updated, date_created) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	db := system.SocialContext.PostgresDB

	now := system.GetUTCTime()
	_, err := db.Exec(query, user.Id, user.FullName, user.Email, user.Password, user.Contact, user.CountryCode,
		user.Country, user.DateOfBirth, now, now)

	if err != nil {
		log.Errorln(err)
	}
	return err
}
