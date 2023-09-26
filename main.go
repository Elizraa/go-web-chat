package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Elizraa/go-web-chat/api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	// Create a multi-output logger
	multiLogger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Add a file output using lumberjack
	fileLogger := zerolog.New(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   true,
	}).With().Timestamp().Logger()

	// Create a multi-output writer
	multiOutput := zerolog.MultiLevelWriter(
		zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
			FormatMessage: func(i interface{}) string {
				return fmt.Sprintf("| %s |", i)
			},
		},
		fileLogger,
	)

	// Set the logger to use the multi-output writer
	log.Logger = multiLogger.Output(multiOutput)

	api.Run()
}
