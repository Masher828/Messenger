package system

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}

func NowInUTC() time.Time {
	return time.Now().UTC()
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
