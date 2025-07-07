package logger

import (
	"context"
	"fmt"
	"github.com/tittuvarghese/go-core-wrappers/constants"
	"github.com/tittuvarghese/go-core-wrappers/timewrapper"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	timeFormat           = constants.TimestampFormat
	defaultLogRetentions = 24
)

type Attribute struct {
	Key   string
	Value string
}

// Logger :
type Logger interface {
	Panic(message string, args ...string)
	Fatal(message string, args ...string)
	Error(message string, err error, args ...string)
	Warn(message string, args ...string)
	Info(message string, args ...string)
	Debug(message string, args ...string)
}

// LoggingService :
type LoggingService struct {
	// Level          string
	ModuleName     string
	context        context.Context //nolint:containedctx
	Filename       string
	file           *os.File
	RotateDuration time.Duration
}

type LogOption struct {
	NodeId    string
	ReplicaId int
	ShardId   int
	Module    string
}

type LoggingServiceOptions struct {
	FilenamePrefix    string
	RetentionDuration time.Duration
}

func NewLogger(module string, args ...LoggingServiceOptions) Logger {

	zerolog.TimeFieldFormat = timeFormat
	SetLogLevel(constants.DebugLevel)

	var rotateDuration time.Duration
	var fileName string

	for _, arg := range args {

		if arg.FilenamePrefix != "" {
			fileName = arg.FilenamePrefix
		} else {
			fileName = "logs"
		}

		if arg.RetentionDuration > 0 {
			rotateDuration = arg.RetentionDuration
		} else {
			rotateDuration = defaultLogRetentions * time.Hour
		}

		break
	}
	logService := &LoggingService{
		ModuleName: module,
		//  Level:          level,
		Filename:       fileName,
		RotateDuration: rotateDuration,
	}
	return logService
}

func SetLogLevel(level string) {
	switch level {
	case constants.PanicLevel:
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case constants.FatalLevel:
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case constants.ErrorLevel:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case constants.WarnLevel:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case constants.InfoLevel:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case constants.DebugLevel:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case constants.TraceLevel:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

// Panic :
func (loggingService *LoggingService) Panic(message string, args ...string) {

	logString := log.Ctx(loggingService.context).Panic().Str("module", loggingService.ModuleName)
	var attribute Attribute
	for len(args) > 0 {
		attribute, args = argsToAttr(args)
		logString.Str(attribute.Key, attribute.Value)
	}
	logString.Msg(message)
}

// Fatal :
func (loggingService *LoggingService) Fatal(message string, args ...string) {

	logString := log.Ctx(loggingService.context).Fatal().Str("module", loggingService.ModuleName)
	var attribute Attribute
	for len(args) > 0 {
		attribute, args = argsToAttr(args)
		logString.Str(attribute.Key, attribute.Value)
	}
	logString.Msg(message)

}

// Error :
func (loggingService *LoggingService) Error(message string, err error, args ...string) {

	logString := log.Error().Str("module", loggingService.ModuleName).Ctx(loggingService.context).Caller(1).Stack().Err(err)
	var attribute Attribute
	for len(args) > 0 {
		attribute, args = argsToAttr(args)
		logString.Str(attribute.Key, attribute.Value)
	}
	logString.Msg(message)
}

// Warn :
func (loggingService *LoggingService) Warn(message string, args ...string) {
	logString := log.Warn().Str("module", loggingService.ModuleName).Ctx(loggingService.context).Caller(1).Stack()
	var attribute Attribute
	for len(args) > 0 {
		attribute, args = argsToAttr(args)
		logString.Str(attribute.Key, attribute.Value)
	}
	logString.Msg(message)
}

// Info :
func (loggingService *LoggingService) Info(message string, args ...string) {
	logString := log.Info().Str("module", loggingService.ModuleName).Ctx(loggingService.context)
	var attribute Attribute
	for len(args) > 0 {
		attribute, args = argsToAttr(args)
		logString.Str(attribute.Key, attribute.Value)
	}
	logString.Msg(message)
}

// Debug :
func (loggingService *LoggingService) Debug(message string, args ...string) {

	logString := log.Debug().Str("module", loggingService.ModuleName).Ctx(loggingService.context).Caller(1).Stack()
	var attribute Attribute
	for len(args) > 0 {
		attribute, args = argsToAttr(args)
		logString.Str(attribute.Key, attribute.Value)
	}
	logString.Msg(message)
}

// Rotate :
func (loggingService *LoggingService) Rotate() error {

	return loggingService.openNew()
}

// openNew :
func (loggingService *LoggingService) openNew() (err error) {
	if loggingService.Filename != "" {
		newName := nextName(loggingService.Filename)
		loggingService.file, err = os.OpenFile(newName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		return
	}
	loggingService.file = os.Stdout // default output

	return
}

// Close :
func (loggingService *LoggingService) Close() (err error) {
	err = fmt.Errorf("%w", loggingService.file.Close())

	return
}

// getCtx :
func (loggingService *LoggingService) getCtx() (context.Context, error) {
	if _, err := os.Stat(loggingService.file.Name()); err != nil {
		err := loggingService.openNew()
		if err != nil {

			return nil, err
		}
	}

	return zerolog.New(loggingService.file).With().Timestamp().
		CallerWithSkipFrameCount(3).
		Logger().
		WithContext(context.Background()), nil
}

// nextName : Generic function to get filename with timestamp
func nextName(name string) string {
	filename := name
	ext := filepath.Ext(filename)
	prefix := filename[:len(filename)-len(ext)]
	timestamp := timewrapper.NewTime().GetCurrentTime().Format(timeFormat)

	return filepath.Join(fmt.Sprintf("%s-%s%s", prefix, timestamp, ext)) //nolint:gocritic
}

// watchdog : Structure which keeps track of duration after
type watchdog struct {
	interval time.Duration
}

// Run :
func (watchdog *watchdog) Run(l *LoggingService) {
	ticker := time.NewTicker(watchdog.interval)
	for range ticker.C {
		if err := l.Rotate(); err != nil {

			return
		}
	}
}

// runWatchdog :
func runWatchdog(loggingService *LoggingService, interval time.Duration) {
	wd := watchdog{interval: interval}
	go wd.Run(loggingService)
}

func argsToAttr(args []string) (Attribute, []string) {
	if len(args) == 1 {
		return Attribute{Key: constants.BadKey, Value: args[0]}, nil
	}
	return Attribute{Key: args[0], Value: args[1]}, args[2:]
}
