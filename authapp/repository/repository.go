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

func GetUserByEmail(email string, log *logrus.Entry) (*models.UserModel, error) {
	query := `SELECT id, name, email, password FROM social_user WHERE email = $1`

	db := system.SocialContext.PostgresDB

	var user models.UserModel

	err := db.QueryRow(query, email).Scan(&user.Id, &user.FullName, &user.Email, &user.Password)

	if err != nil {
		log.Errorln(err)
	}
	return &user, err
}

func GetAllUsers(log *logrus.Entry) ([]string, error) {
	query := "SELECT name FROM social_user"

	db := system.SocialContext.PostgresDB

	rows, err := db.Query(query)
	if err != nil {
		log.Errorln(err)
		return []string{}, err
	}

	var names []string
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			log.Error(err)
			continue
		}

		names = append(names, name)
	}
	return names, nil
}
