package daplogger_test

import (
	"path"
	"testing"

	"github.com/KingDanx/daplogger"
)

func TestLog(t *testing.T) {
	logDir := path.Join("logs")

	logger := daplogger.Logger{
		Path:         logDir,
		LogName:      "test",
		LogFileCount: 30,
	}

	logger.LogInfo("test")
}
