package models

import (
	"regexp"
	"strings"
	"time"

	"github.com/Masher828/MessengerBackend/common-packages/system"
)

type UserModel struct {
	Id          int64     `json:"id"`
	FullName    string    `json:"name" column:"name"`
	Email       string    `json:"email" column:"email"`
	Password    string    `json:"password" column:"password"`
	Contact     string    `json:"contact" column:"contact"`
	CountryCode string    `json:"contry_code" column:"country_code"`
	Country     string    `json:"country" column:"country"`
	DateOfBirth time.Time `json:"date_of_birth" column:"date_of_birth"`
	DateCreated time.Time `json:"date_created" column:"date_created"`
	LastUpdated time.Time `json:"last_updated" column:"last_updated"`
}

type UserLoginModel struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserModelContext struct {
	Id          int64  `json:"id"`
	FullName    string `json:"name" column:"name"`
	Email       string `json:"email" column:"email"`
	AccessToken string `json:"accessToken"`
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func (user *UserModel) IsValid() (bool, error) {
	if len(user.FullName) < 2 {
		return false, system.InvalidNameErr
	}

	if !emailRegex.MatchString(user.Email) {
		return false, system.InvalidEmailErr
	}

	if len(user.Password) < 8 || len(user.Password) > 20 {
		return false, system.InvalidPasswordFormatErr
	}

	if len(user.Contact) != 10 {
		return false, system.InvalidContactNumberErr
	}

	user.Email = strings.ToLower(user.Email)

	return true, nil
}

func (user *UserModel) GetUserContext() *UserModelContext {
	var userContext UserModelContext
	userContext.Email = user.Email
	userContext.FullName = user.FullName
	userContext.Id = user.Id
	return &userContext
}

func (user *UserLoginModel) IsValid() (bool, error) {
	if !emailRegex.MatchString(user.Email) {
		return false, system.InvalidEmailErr
	}

	if len(user.Password) < 8 || len(user.Password) > 20 {
		return false, system.InvalidPasswordFormatErr
	}

	user.Email = strings.ToLower(user.Email)

	return true, nil
}
