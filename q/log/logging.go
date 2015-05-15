package log

import (
	"log"
	"os"
)

func init() {
	log.SetFlags(0)
}

func Verbose(format string, args ...interface{}) {
	if IsVerbose() {
		log.Printf(format, args...)
	}
}

func IsVerbose() bool {
	return os.Getenv("Q_VERBOSE") != ""
}

func Debug(format string, args ...interface{}) {
	if IsDebug() {
		log.Printf(format, args...)
	}
}

func IsDebug() bool {
	return os.Getenv("Q_DEBUG") != ""
}
