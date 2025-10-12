package util

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

// ReturnLogLevels : Return a comma-delimited list of log levels
func ReturnLogLevels(levelMap map[string]logrus.Level) string {
	logLevels := make([]string, 0, len(levelMap))
	for k := range levelMap {
		logLevels = append(logLevels, k)
	}
	sort.Strings(logLevels)

	return strings.Join(logLevels, ", ")
}

// ConfigureLogger : Configure the logger
func ConfigureLogger(logLevel logrus.Level, flagNoColor bool) (logger *logrus.Logger) {
	disableColors := false
	if flagNoColor {
		disableColors = true
	}
	logger = &logrus.Logger{
		Out:   os.Stderr,
		Level: logLevel,
		Formatter: &prefixed.TextFormatter{
			DisableColors:    disableColors,
			DisableTimestamp: true,
			TimestampFormat:  "2006-01-02 15:04:05",
			FullTimestamp:    true,
			ForceFormatting:  false,
		},
	}
	logger.SetLevel(logLevel)

	return logger
}

// Generate hex code
func GenerateHex() (string, error) {
	bytes := make([]byte, 32) // 32 bytes = 64 hex characters
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
