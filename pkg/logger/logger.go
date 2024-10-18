package logger

import (
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"time"

	"github.com/rs/zerolog"
)

type callerInfo struct {
	file         string
	line         int
	functionName string
}

var log *zerolog.Logger

func NewLogger(level zerolog.Level, writer io.Writer) *zerolog.Logger {
	zerolog.SetGlobalLevel(level)

	logger := zerolog.New(writer).With().Logger()

	log = &logger

	return &logger
}

func addCallerContext(logger zerolog.Logger) zerolog.Logger {
	callerInfo, ok := getCallerInfo()
	if ok {
		fileName := filepath.Base(callerInfo.file)
		funcName := filepath.Base(callerInfo.functionName)
		logger = logger.With().
			Str("time", time.Now().Format("2006-01-02 15:04:05.000")).
			Str("file", fileName).
			Int("line", callerInfo.line).
			Str("function", funcName).
			Logger()
	}
	return logger
}

func getCallerInfo() (callerInfo, bool) {
	pc, file, line, ok := runtime.Caller(3)
	if ok {
		functionName := runtime.FuncForPC(pc).Name()
		return callerInfo{file, line, functionName}, true
	}
	return callerInfo{}, false
}

func Debug(msg string) {
	logger := addCallerContext(*log)
	logger.Debug().Msg(msg)
}

func Debugf(format string, a ...any) {
	logger := addCallerContext(*log)
	msg := fmt.Sprintf(format, a...)
	logger.Debug().Msg(msg)
}

func Info(msg string) {
	logger := addCallerContext(*log)
	logger.Info().Msg(msg)
}

func Infof(format string, a ...any) {
	logger := addCallerContext(*log)
	msg := fmt.Sprintf(format, a...)
	logger.Info().Msg(msg)
}

func Warn(msg string) {
	logger := addCallerContext(*log)
	logger.Warn().Msg(msg)
}

func Warnf(format string, a ...any) {
	logger := addCallerContext(*log)
	msg := fmt.Sprintf(format, a...)
	logger.Warn().Msg(msg)
}

func Error(msg string) {
	logger := addCallerContext(*log)
	logger.Error().Msg(msg)
}

func Errorf(format string, a ...any) {
	logger := addCallerContext(*log)
	msg := fmt.Sprintf(format, a...)
	logger.Error().Msg(msg)
}

func Fatal(msg string) {
	logger := addCallerContext(*log)
	logger.Fatal().Msg(msg)
}

func Fatalf(format string, a ...any) {
	logger := addCallerContext(*log)
	msg := fmt.Sprintf(format, a...)
	logger.Fatal().Msg(msg)
}
