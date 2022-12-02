package logger

import (
	"milk/config"
	"os"

	"github.com/rs/zerolog"
)

type Logger interface {
	InitLogger() error
	Debug(msg string)
	Debugf(template string, args ...interface{})
	Info(msg string)
	Infof(template string, args ...interface{})
	Warn(msg string)
	Warnf(template string, args ...interface{})
	Error(err error)
	Errorf(template string, args ...interface{})
	Fatal(msg string)
	Fatalf(template string, args ...interface{})
	Panic(msg string)
	Panicf(template string, args ...interface{})
	Request() *zerolog.Event
}

var loggerLevelMap = map[string]zerolog.Level{
	"debug":    zerolog.DebugLevel,
	"info":     zerolog.InfoLevel,
	"warn":     zerolog.WarnLevel,
	"error":    zerolog.ErrorLevel,
	"panic":    zerolog.PanicLevel,
	"fatal":    zerolog.FatalLevel,
	"noLevel":  zerolog.NoLevel,
	"disabled": zerolog.Disabled,
}

type apiLogger struct {
	cfg    *config.Config
	logger zerolog.Logger
}

func NewAPILogger(
	cfg *config.Config,
) *apiLogger {
	return &apiLogger{
		cfg: cfg,
	}
}

func (a *apiLogger) InitLogger() error {
	l := zerolog.New(os.Stdout)

	a.logger = l.Level(*a.cfg.Logger.Level).With().Timestamp().Logger()
	return nil
}

func (a *apiLogger) Debug(msg string) {
	a.logger.Debug().Msg(msg)
}

func (a *apiLogger) Debugf(template string, args ...interface{}) {
	a.logger.Debug().Msgf(template, args...)
}

func (a *apiLogger) Info(msg string) {
	a.logger.Info().Msg(msg)
}

func (a *apiLogger) Request() *zerolog.Event {
	return a.logger.Info()
}

func (a *apiLogger) Infof(template string, args ...interface{}) {
	a.logger.Info().Msgf(template, args...)
}

func (a *apiLogger) Warn(msg string) {
	a.logger.Warn().Msg(msg)
}

func (a *apiLogger) Warnf(template string, args ...interface{}) {
	a.logger.Warn().Msgf(template, args...)
}

func (a *apiLogger) Error(err error) {
	a.logger.Error().Msg(err.Error())
}

func (a *apiLogger) Errorf(template string, args ...interface{}) {
	a.logger.Error().Msgf(template, args...)
}

func (a *apiLogger) Panic(msg string) {
	a.logger.Panic().Msg(msg)
}

func (a *apiLogger) Panicf(template string, args ...interface{}) {
	a.logger.Panic().Msgf(template, args...)
}

func (a *apiLogger) Fatal(msg string) {
	a.logger.Fatal().Msg(msg)
}

func (a *apiLogger) Fatalf(template string, args ...interface{}) {
	a.logger.Fatal().Msgf(template, args...)
}
