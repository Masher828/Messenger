package log

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/smtp"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
	if viper.GetString("mail.sendActualMail") == "Y" {
		go SendMail(entry)
	}
	return nil
}

func GetMessage(entry *logrus.Entry, from, password string) string {

	timeFormat := "02-January-2006 :: 15:04:05.000000"
	subject := "Subject: Exception occurred in " + os.Getenv("ENV") + "env at " + time.Now().Format(timeFormat) + "\r\n"

	data := struct {
		FileName     string
		FunctionName string
		Error        string
		Level        string
	}{
		FileName:     path.Base(entry.Caller.File) + ":" + strconv.Itoa(entry.Caller.Line),
		FunctionName: entry.Caller.Function,
		Error:        entry.Message,
		Level:        entry.Level.String(),
	}

	t, err := template.New("").Parse(`
	File : {{.FileName}},
	Function : {{.FunctionName}},
	Error : {{.Error}},
	Level : {{.Level}}
	`)
	if err != nil {
		fmt.Println(err)
	}

	var message bytes.Buffer

	err = t.Execute(&message, data)
	if err != nil {
		fmt.Println(err)
	}

	msg := ""
	msg += fmt.Sprintf("From: %s\r\n", from)
	msg += fmt.Sprintf("To: %s\r\n", from)
	msg += fmt.Sprintf("%s\r\n", subject)
	msg += fmt.Sprintf("%s\r\n", message.String())

	return msg
}

func SendMail(entry *logrus.Entry) {

	from := viper.GetString("mail.from")
	password := viper.GetString("mail.password")

	toList := []string{from}

	host := viper.GetString("mail.host")

	port := viper.GetString("mail.port")

	message := GetMessage(entry, from, password)

	auth := smtp.PlainAuth("", from, password, host)

	err := smtp.SendMail(host+":"+port, auth, from, toList, []byte(message))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
