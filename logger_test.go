package daplogger_test

import (
	"testing"

	"github.com/KingDanx/daplogger"
)

// func TestLog(t *testing.T) {
// 	logDir := path.Join("logs")

// 	logger := daplogger.Logger{
// 		Path:         logDir,
// 		LogName:      "test",
// 		LogFileCount: 30,
// 	}

// 	logger.LogInfo("test")
// }

func TestCreate(t *testing.T) {
	logger := daplogger.CreateLogger("logs", "test", 1)

	logger.LogInfo("test 1")
}
