package zlog

import (
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	zLog "github.com/rs/zerolog/log"
	"golang.org/x/term"
)


var initOnce sync.Once

func InitZerolog() *zerolog.Logger {
	initOnce.Do(func() {
		zerolog.TimeFieldFormat = time.RFC3339Nano

		writer := io.Writer(os.Stdout)
		if term.IsTerminal(0) {
			writer = zerolog.NewConsoleWriter()
		}

		zLog.Logger = zerolog.New(writer).With().Timestamp().Logger()
		zerolog.DefaultContextLogger = &zLog.Logger

		log.SetOutput(&zLog.Logger)
	})

	return &zLog.Logger
}
