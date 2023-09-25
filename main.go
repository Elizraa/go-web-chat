package main

import (
	"os"

	"github.com/DLzer/go-product-api/api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	// Configure Zerolog to use the pretty formatter
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(&lumberjack.Logger{
		Filename:   "logs/app.log", // Specify the log file path
		MaxSize:    10,               // Max size (in megabytes) before log rotation
		MaxBackups: 10,               // Max number of old log files to retain
		MaxAge:     30,               // Max number of days to retain old log files
		Compress:   true,             // Whether to compress old log files
	})

	api.Run()
}

