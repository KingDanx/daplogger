package daplogger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Logger struct {
	Path         string
	LogName      string
	LogFileCount int
}

func (l Logger) Log(message string, messageType string) {
	logFilePath, err := l.createLogFile()
	if err != nil {
		fmt.Println("Failed to get the log file")
	}

	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Println("Failed to open the log file")
	}
	defer logFile.Close()

	now := time.Now()
	mins := fmt.Sprintf("%02d", now.Minute())
	secs := fmt.Sprintf("%02d", now.Second())

	formatedMessage := fmt.Sprintf("%d:%v:%v - [%s] - %s\r\n", now.Hour(), mins, secs, strings.ToUpper(messageType), message)

	_, writeErr := logFile.WriteString(formatedMessage)
	if writeErr != nil {
		fmt.Println("Failed to write to the log file")
	}
}

func (l Logger) LogInfo(message string) {
	l.Log(message, "info")
}

func (l Logger) LogError(message string) {
	l.Log(message, "error")
}

func (l Logger) LogWarning(message string) {
	l.Log(message, "warning")
}

func (l Logger) createLogFile() (string, error) {

	//? Check if directory does not exist
	_, dirErr := os.Stat(l.Path)
	if os.IsNotExist(dirErr) {
		//? Directory does not exist
		os.MkdirAll(l.Path, os.ModePerm)
	}

	formatLogName := fmt.Sprintf("%d_%d_%d-%s.log", time.Now().Year(), time.Now().Month(), time.Now().Day(), l.LogName)

	fullPath := filepath.Join(l.Path, formatLogName)

	//? Check if file does not exist
	_, fileErr := os.Stat(fullPath)
	if os.IsNotExist(fileErr) {
		//? Create the file in the directory
		file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
		if err != nil {
			return "", err
		}
		defer file.Close()
	}

	formatLatestLogName := fmt.Sprintf("latest-%s.log", l.LogName)

	fullPathLatest := filepath.Join(l.Path, formatLatestLogName)

	latestFile, fileErrLatest := os.Stat(fullPathLatest)
	if os.IsNotExist(fileErrLatest) {
		//? Create the file in the directory
		file, err := os.OpenFile(fullPathLatest, os.O_CREATE|os.O_RDWR, os.ModePerm)
		if err != nil {
			return "", err
		}
		defer file.Close()
	}

	if latestFile.Size() > 50*1024*1024 {
		removeErr := os.Remove(fullPathLatest)
		if removeErr != nil {
			fmt.Println("Failed to delete file")
		}
		file, err := os.OpenFile(fullPathLatest, os.O_CREATE|os.O_RDWR, os.ModePerm)
		if err != nil {
			return "", err
		}
		defer file.Close()
	}

	return fullPath, nil
}
