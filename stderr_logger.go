package cache

import (
	"log"
	"os"
)

type stderrLogger struct {
	logger *log.Logger
}

func (sl stderrLogger) Errorf(format string, args ...interface{}) {
	sl.logger.Printf(format, args...)
}

func StderrLogger() LoggerInterface {
	return stderrLogger{logger: log.New(os.Stderr, "ERROR CACHE ", log.LstdFlags)}
}
