package logger

import (
	"log"

	"github.com/natefinch/lumberjack"
)

func SetLogger(filename string) *log.Logger {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    10,   // Max file size before rotation
		MaxBackups: 3,    // Max old log files to retain
		MaxAge:     28,   // Max days to retain a log files
		Compress:   true, // Compress old files
	}

	return log.New(lumberjackLogger, "", log.LstdFlags|log.Lshortfile)
}
