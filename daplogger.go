package daplogger

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Logger struct {
	Path         string
	LogName      string
	LogFileCount int
	LogFiles     LogFiles
}

type LogFiles struct {
	CurrentDay string
	Latest     string
}

type DAPFile struct {
	file os.FileInfo
	epoc int64
}

func (l *Logger) Log(message string, messageType string) {
	logFiles, err := l.createLogFile()
	if err != nil {
		fmt.Println("Failed to get the log file")
	}

	logFile, err := os.OpenFile(logFiles.CurrentDay, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Println("Failed to open the log file")
	}
	defer logFile.Close()

	latestFile, err := os.OpenFile(logFiles.Latest, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Println("Failed to open the log file")
	}
	defer latestFile.Close()

	now := time.Now()
	mins := fmt.Sprintf("%02d", now.Minute())
	secs := fmt.Sprintf("%02d", now.Second())

	formatedMessage := fmt.Sprintf("[%s] %d-%d-%d - %d:%v:%v - %s\r\n", strings.ToUpper(messageType), now.Month(), now.Day(), now.Year(), now.Hour(), mins, secs, message)

	_, writeErr := logFile.WriteString(formatedMessage)
	if writeErr != nil {
		fmt.Println("Failed to write to the log file")
	}

	_, latestWriteErr := latestFile.WriteString(formatedMessage)
	if latestWriteErr != nil {
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

func (l *Logger) createLogFile() (LogFiles, error) {

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
			return LogFiles{}, err
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
			return LogFiles{}, err
		}
		defer file.Close()
	} else if latestFile.Size() > 50*1024*1024 { //? if the file is over 50mb we will delete it and make a new file
		removeErr := os.Remove(fullPathLatest)
		if removeErr != nil {
			fmt.Println("Failed to delete file")
		}
		file, err := os.OpenFile(fullPathLatest, os.O_CREATE|os.O_RDWR, os.ModePerm)
		if err != nil {
			return LogFiles{}, err
		}
		defer file.Close()
	}

	logFiles := LogFiles{
		CurrentDay: fullPath,
		Latest:     fullPathLatest,
	}

	l.LogFiles = logFiles

	return logFiles, nil
}

func (l *Logger) cleanLogs() {
	hour, min, sec, nsec := 0, 0, 0, 50
	for {
		files := []DAPFile{}
		err := filepath.Walk(l.Path, func(path string, file os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !file.IsDir() && !strings.Contains(file.Name(), "latest") {
				modTime := file.ModTime()
				epoc := modTime.Unix()

				dapFile := DAPFile{
					file: file,
					epoc: epoc,
				}

				files = append(files, dapFile)
			}

			return nil
		})

		if len(files) > l.LogFileCount {
			sort.Slice(files, func(i, j int) bool {
				return files[i].epoc > files[j].epoc
			})

			filesToDelete := files[l.LogFileCount:]

			for _, file := range filesToDelete {
				fmt.Println("I must delete this file: ", file.file.Name())
				err := os.Remove(path.Join(l.Path, file.file.Name()))
				if err != nil {
					fmt.Println("Error deleting: ", file.file.Name())
				}
			}
		}

		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), hour, min, sec, nsec, now.Location())

		if now.After(next) {
			next = next.Add(24 * time.Hour)
		}

		duration := time.Until(next)
		fmt.Printf("Sleeping for %v until next trigger at %v\n", duration, next)

		time.Sleep(duration)

		if err != nil {
			panic(fmt.Sprintf("Error walking the path %q: %v\n", l.Path, err))
		}
	}
}

func CreateLogger(path, logName string, logCount int) Logger {
	logger := Logger{
		Path:         path,
		LogName:      logName,
		LogFileCount: logCount,
	}

	go logger.cleanLogs()

	return logger
}
