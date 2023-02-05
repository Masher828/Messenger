package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func GetDefaultLogger(userId int64, url string, method string) *zap.SugaredLogger {

	logg := zap.NewProductionConfig()
	logg.Encoding = "console"
	//logg.Enco
	logg.DisableStacktrace = true
	logg.EncoderConfig.CallerKey = "caller"
	logg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	logg.EncoderConfig.TimeKey = "timestamp"
	logg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	logg.EncoderConfig.MessageKey = "message"
	logg.EncoderConfig.LevelKey = "level"
	logg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logg.OutputPaths = []string{"logg.log", "stderr"}
	logger, err := logg.Build()
	if err != nil {
		fmt.Println(err)
	}
	return logger.With([]zap.Field{{Key: "UserId", Type: zapcore.Int64Type, Integer: userId}, {Key: "url", Type: zapcore.StringType, String: url}, {Key: "method", Type: zapcore.StringType, String: method}}...).Sugar()
}
