package system

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}

type uRIParam struct {
	Limit  int64 `form:"limit"`
	Offset int64 `form:"offset"`
}

func NowInUTC() time.Time {
	return time.Now().UTC()
}

func NowInUTCMicro() int64 {
	return NowInUTC().UnixMicro()
}

func GetRandomHashSalt(saltSize int64) ([]byte, error) {
	var salt = make([]byte, saltSize)

	_, err := rand.Read(salt)
	if err != nil {
		fmt.Println(err)
	}

	return salt, err
}

func ContainsString(array []string, element string) bool {
	for _, ele := range array {
		if ele == element {
			return true
		}
	}
	return false
}

func GetDifferentFields(existingData, newData interface{}) map[string]interface{} {

	existingDataValues := reflect.ValueOf(existingData)
	newDataValues := reflect.ValueOf(newData)

	existingDataTags := reflect.TypeOf(existingData)
	newDataTags := reflect.TypeOf(newData)

	difference := map[string]interface{}{}

	fieldsToSkip := []string{"_id", "createdOn", "updatedOn"}
	for i := 0; i < existingDataValues.NumField(); i++ {

		if ContainsString(fieldsToSkip, existingDataTags.Field(i).Tag.Get("bson")) {
			continue
		}
		if existingDataValues.Field(i).Interface() != newDataValues.Field(i).Interface() && existingDataTags.Field(i).Tag.Get("bson") == newDataTags.Field(i).Tag.Get("bson") {
			difference[existingDataTags.Field(i).Tag.Get("bson")] = newDataValues.Field(i).Interface()
		}
	}

	return difference

}

func HashPassword(salt []byte, password string) string {
	passwordBytes := []byte(password)

	var sha512Hash = sha512.New()

	// Append salt to password
	passwordBytes = append(passwordBytes, salt...)

	// Write password bytes to the hash
	sha512Hash.Write(passwordBytes)

	// Get the SHA-512 hashed password
	var hashedPasswordBytes = sha512Hash.Sum(nil)

	return base64.URLEncoding.EncodeToString(hashedPasswordBytes)
}

func GetHashedPassword(password string) ([]byte, string, error) {

	salt, err := GetRandomHashSalt(DefaultHashSaltSize)
	if err != nil {
		return []byte{}, "", err
	}

	hashedPass := HashPassword(salt, password)

	return salt, hashedPass, nil
}

func GetUserContextFromGinContext(c *gin.Context) *UserContext {
	userContext, ok := c.Get(AuthUserContext)
	if !ok || userContext == nil {
		return nil
	}

	return userContext.(*UserContext)
}

func GetUserContextFromAccessToken(accessToken string) (*UserContext, error) {
	db := MessengerContext.Redis

	result := db.Get(context.TODO(), fmt.Sprintf(AccessTokenToUser, accessToken))

	if result.Err() != nil {
		fmt.Printf("Error while getting the user from accessToken : %s Error : %s", accessToken, result.Err().Error())
		return nil, result.Err()
	}
	userDetailsByte, err := result.Bytes()
	if err != nil {
		fmt.Printf("Error while getting the user bytes from accessToken : %s Error : %s", accessToken, result.Err().Error())
		return nil, err
	}

	var userDetails UserContext
	err = json.Unmarshal(userDetailsByte, &userDetails)
	if err != nil {
		fmt.Printf("Error while getting the unmarshalling the user bytes from accessToken : %s Error : %s", accessToken, result.Err().Error())
		return nil, err
	}

	return &userDetails, nil
}

func getAccessTokenFromContext(c *gin.Context) string {
	bearerToken := c.GetHeader("Authorization")
	bearerTokenSplit := strings.Split(bearerToken, " ")

	if len(bearerTokenSplit) != 2 {
		return ""
	}

	return bearerTokenSplit[1]
}

func GetOffsetAndLimitFromContext(c *gin.Context, defaultLimit int64) (int64, int64) {
	urlParam := uRIParam{}
	err := c.Bind(&urlParam)
	if err != nil {
		fmt.Println(err.Error())
		return 0, 0
	}
	if urlParam.Limit == 0 {
		urlParam.Limit = defaultLimit
	}

	return urlParam.Offset, urlParam.Limit
}
