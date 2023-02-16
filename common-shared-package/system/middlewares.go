package system

import (
	"encoding/json"
	"fmt"
	"github.com/Masher828/MessengerBackend/common-shared-package/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"reflect"
	"runtime/debug"
	"strconv"
	"time"
)

type Application struct {
}

func (application *Application) Route(controller interface{}, controllerName string, isPublic bool) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		authFailed, ok := c.Get(AuthFailed)
		if !ok {
			authFailed = false
		}

		//if url is not public and authentication is failed then return 401
		if !isPublic && (!ok || authFailed.(bool)) {
			c.Writer.WriteHeader(http.StatusUnauthorized)
		} else {
			var user *UserContext
			var userid string
			if !isPublic {
				userInterface, ok := c.Get(AuthUserContext)
				if !ok {
					c.Redirect(http.StatusUnauthorized, "/")
				}
				user = userInterface.(*UserContext)
			}

			if user != nil {
				userid = user.UserId
			}
			logger := log.GetDefaultLogger(userid, c.Request.RequestURI, c.Request.Method)

			methodInterface := reflect.ValueOf(controller).MethodByName(controllerName).Interface()

			method := methodInterface.(func(c *gin.Context, log *zap.SugaredLogger) ([]byte, error))
			response, err := method(c, logger)
			code := http.StatusOK

			if err != nil {
				responseMap := map[string]string{}
				if IsFunctionalError(err) {
					responseMap["err"] = err.Error()
					code = http.StatusPreconditionFailed
				} else {
					responseMap["err"] = "something went wrong please try again later"
					code = http.StatusInternalServerError
				}
				response, err = json.Marshal(responseMap)
			}

			c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
			c.Header("Access-Control-Allow-Headers", "*")
			c.Header("Access-Control-Allow-Methods", "*")
			if _, ok := c.Get("Content-Type"); ok {
				c.Header("Content-Type", c.GetString("Content-Type"))
			} else {
				c.Header("Content-Type", c.GetHeader("Content-Type"))
			}
			c.Writer.WriteHeader(code)
			c.Writer.Write(response)

		}

	}

	return fn
}

func (application *Application) Recovery(c *gin.Context) {
	defer func() {
		err := recover()

		if err != nil {
			debug.PrintStack()
			response := map[string]interface{}{"success": false, "err": ErrInternalServer}
			c.JSON(http.StatusInternalServerError, response)
		}

	}()
	c.Next()
}

func (application *Application) PerformanceMeasure() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		fmt.Println("Request for url : "+c.Request.URL.RequestURI()+" is started at ", c.Request.URL, NowInUTC())
		c.Set(RequestStartTime, NowInUTC())

		c.Next()

		timeTakenByRequestInString := strconv.FormatInt(time.Since(c.GetTime(RequestStartTime)).Microseconds(), 10)

		fmt.Printf("Request for url : %s|%s is completed in %s microseconds with status %d.\n", c.Request.Method, c.Request.URL, timeTakenByRequestInString, c.Writer.Status())

	}
	return fn
}

func (application *Application) Cors() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
			c.Status(http.StatusOK)
			c.Writer.Write([]byte{})
		} else {
			c.Next()
		}
	}

	return fn
}

func (application *Application) ApplyAuth() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		accessToken := getAccessTokenFromContext(c)
		if len(accessToken) == 0 {
			c.Set(AuthFailed, true)
		} else {
			userContext, err := GetUserContextFromAccessToken(accessToken)
			if err != nil {
				fmt.Println(err.Error())
			}

			if userContext == nil {
				c.Set(AuthFailed, true)
			} else {
				c.Set(AuthFailed, false)
				c.Set(AuthUserContext, userContext)
			}
		}

		c.Next()
	}
	return fn
}
