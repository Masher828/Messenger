package log

import (
	"path"
	"runtime"
	"strconv"
	"strings"

	logger "github.com/sirupsen/logrus"
)

func GetDefaultLogger(userId int64, uri string, method string) *logger.Entry {
	logger.SetReportCaller(true)
	logger.SetFormatter(&logger.TextFormatter{
		DisableColors: false,
		ForceColors:   true,
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			fileName := path.Base(f.File) + ":" + strconv.Itoa(f.Line)
			s := strings.Split(f.Function, ".")
			return s[len(s)-1], fileName
		},
	})
	// logger.SetFormatter(&logger.JSONFormatter{
	// 	CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
	// 		fileName := path.Base(f.File) + ":" + strconv.Itoa(f.Line)
	// 		s := strings.Split(f.Function, ".")
	// 		return s[len(s)-1], fileName
	// 	},
	// })
	return logger.WithFields(logger.Fields{
		"userId": userId,
		"uri":    uri,
		"method": method,
	})
}
