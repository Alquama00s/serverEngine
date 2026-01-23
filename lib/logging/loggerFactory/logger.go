package loggerFactory

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func GetLogger(context ...string) *zerolog.Logger {
	// Configure the logger to use the ConsoleWriter
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,             // Output to stdout (can be os.Stderr or a file)
		TimeFormat: "2006-01-02 15:04:05", // Optional: customize timestamp format
	}

	logger := log.Output(consoleWriter).With().
		Strs("context", context).
		Timestamp().
		Logger()

	return &logger
}
