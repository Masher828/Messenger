package log

import (
	"fmt"
	"io"
	"net/smtp"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	logger "github.com/sirupsen/logrus"
)

func GetDefaultLogger(userId int64, uri string, method string) *logger.Entry {
	logger.SetReportCaller(true)

	f, err := os.OpenFile("testlogrus.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
	}

	// defer f.Close()

	// logger.SetFormatter(&logger.TextFormatter{
	// 	DisableColors: false,
	// 	ForceColors:   true,
	// 	CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
	// 		fileName := path.Base(f.File) + ":" + strconv.Itoa(f.Line)
	// 		s := strings.Split(f.Function, ".")
	// 		return s[len(s)-1], fileName
	// 	},
	// })

	logger.SetFormatter(&logger.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			fileName := path.Base(f.File) + ":" + strconv.Itoa(f.Line)
			s := strings.Split(f.Function, ".")
			return s[len(s)-1], fileName
		},
		PrettyPrint: true,
	})

	mw := io.MultiWriter(os.Stdout, f)

	logger.SetOutput(mw)

	logger.AddHook(NewExtraFieldHook("local"))

	return logger.WithFields(logger.Fields{
		"method": method,
		"uri":    uri,
		"userId": userId,
	})
}

type ExtraFieldHook struct {
	env string
	pid int
}

func NewExtraFieldHook(env string) *ExtraFieldHook {
	return &ExtraFieldHook{
		env: env,
		pid: os.Getpid(),
	}
}

func (h *ExtraFieldHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *ExtraFieldHook) Fire(entry *logrus.Entry) error {
	fmt.Println(entry.Data, entry.Level, entry.Message)
	SendMail("entry.Data[].(string)")
	return nil
}

func SendMail(message string) {
	from := "******"
	password := "**********"
	toList := []string{""}
	host := "smtp.gmail.com"

	port := "587"

	body := []byte(message)

	auth := smtp.PlainAuth("", from, password, host)

	err := smtp.SendMail(host+":"+port, auth, from, toList, body)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
