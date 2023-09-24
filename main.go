package main

import (
	"os"

	"github.com/DLzer/go-product-api/api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Configure Zerolog to use the pretty formatter
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// You can configure the logger here if needed
	// For example, setting an output file, adding additional context, etc.
	// log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	log.Error().Msg("Error message")
	log.Warn().Msg("Warning message")
	log.Info().Msg("Info message")
	log.Debug().Msg("Debug message")
	log.Trace().Msg("Trace message")
	api.Run()
}
