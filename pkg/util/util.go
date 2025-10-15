package util

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"golang.org/x/term"
)

// ConfigureLogger : Configure the logger
func ConfigureLogger(flagNoColor bool, logFileName string) (logger *logrus.Logger) {
	var (
		disableColors    bool = false
		disableTimestamp bool = true
		err              error
		isTerminal       bool = false
		logfile          *os.File
		out              *os.File = os.Stderr
	)

	if term.IsTerminal(int(os.Stdout.Fd())) {
		isTerminal = true
	}

	if flagNoColor {
		disableColors = true
	}

	if !isTerminal {
		disableColors = true
		disableTimestamp = false
		logfile, err = os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err == nil {
			out = logfile
		}
	}

	logger = &logrus.Logger{
		Out:   out,
		Level: logrus.InfoLevel,
		Formatter: &prefixed.TextFormatter{
			DisableColors:    disableColors,
			DisableTimestamp: disableTimestamp,
			TimestampFormat:  "2006-01-02 15:04:05",
			FullTimestamp:    true,
			ForceFormatting:  false,
		},
	}

	return logger
}

// Return true if the specified path exists and false if it doesn't
func PathExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return true
	}
}

// Return true if the specified path is a directory and false if it isn't
func IsDir(path string) bool {
	if PathExists(path) {
		info, _ := os.Stat(path)
		if info.IsDir() {
			return true
		}
	}
	return false
}

// Return true if the specified path is writable and false if it isn't
func IsWritable(path string) bool {
	if IsDir(path) {
		tmpFile := "tmpfile"
		file, err := os.CreateTemp(path, tmpFile)
		if err != nil {
			return false
		}
		defer os.Remove(file.Name())
		defer file.Close()
		return true
	}
	return false
}

// Return true if the specified file size > 0 and false if it isn't
func FileSize(filename string) int64 {
	if info, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return -1
	} else {
		return info.Size()
	}
}

// Return a list of files that is older than olderThan
func FindOldFiles(pattern string, olderThan time.Duration) ([]string, error) {
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("glob error: %w", err)
	}

	var oldFiles []string
	cutoff := time.Now().Add(-olderThan)

	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			log.Printf("skipping %s: %v", file, err)
			continue
		}

		if !info.IsDir() && info.ModTime().Before(cutoff) {
			oldFiles = append(oldFiles, file)
		}
	}

	return oldFiles, nil
}

func GetHome() (homeDir string, err error) {
	userObj, err := user.Current()
	if err != nil {
		homeDir = os.Getenv("HOME")
		if homeDir == "" {
			return "", fmt.Errorf("failed to determine the path of your home directory: %s", err.Error())
		}
	}

	return userObj.HomeDir, nil
}

func VerifyDirectory(path string) error {
	if PathExists(path) {
		if IsDir(path) {
			if IsWritable(path) {
				return nil
			}
			return fmt.Errorf("%s exists and is a directory, but isn't writable", path)
		} else {
			return fmt.Errorf("%s exists but isn't a directory", path)
		}
	} else {
		err := os.MkdirAll(path, 0700)
		if err != nil {
			return fmt.Errorf("the path %s doesn't exist and couldn't be created: %s", path, err.Error())
		}
		return nil
	}
}
