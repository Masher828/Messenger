package ServicesTest

import (
	"github.com/Masher828/MessengerBackend/common-shared-package/log"
	"go.uber.org/zap"
	"testing"
)

func Test_HH(t *testing.T) {
	logger := log.GetDefaultLogger()
	logger.Core().With([]zap.Field{{Key: "test", String: "value"}})

	logger.Error("helllo")
}
