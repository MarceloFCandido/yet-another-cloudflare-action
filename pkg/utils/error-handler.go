package utils

import (
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"yaca/pkg/logger"
)

func HandleError(err error, msg string, args ...any) {
	if err != nil {
		args = append(args, slog.String("error", err.Error()))
		args = append(args, slog.String("error_type", fmt.Sprintf("%T", err)))
		
		logger.Error(msg, args...)
		
		if os.Getenv("ENVIRONMENT") != "production" {
			logger.Debug("Error stack trace", 
				slog.String("stack", string(debug.Stack())))
		}
		
		os.Exit(1)
	}
}

func LogAndExit(code int, msg string, args ...any) {
	if code == 0 {
		logger.Info(msg, args...)
	} else {
		logger.Error(msg, args...)
	}
	os.Exit(code)
}
