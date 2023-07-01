package logger

import (
	"auth-svc/config"
	"os"

	"github.com/rs/zerolog"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Logger interface {
	InitLogger() error
	GetLogger() *zerolog.Logger
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
	tgBot  *tb.Bot
	logger zerolog.Logger
}

func NewAPILogger(cfg *config.Config) Logger {
	return &apiLogger{cfg: cfg}
}

func (a *apiLogger) GetLogger() *zerolog.Logger {
	return &a.logger
}

func (a *apiLogger) InitLogger() error {
	var w zerolog.LevelWriter
	if a.cfg.Logger.InFile {
		logFile, err := os.OpenFile(a.cfg.Logger.FilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return err
		}
		w = zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout}, logFile)
	} else {
		w = zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	if a.cfg.Logger.InTG {
		err := a.InitTG()
		if err != nil {
			return err
		}
	}
	a.logger = zerolog.New(w).Level(loggerLevelMap[a.cfg.Logger.Level]).With().
		CallerWithSkipFrameCount(a.cfg.Logger.SkipFrameCount).Timestamp().Logger().Hook(a)
	return nil
}

func (a *apiLogger) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if !a.cfg.Logger.InTG || loggerLevelMap[a.cfg.Logger.TGLevel] > level {
		return
	}
	go a.SendTGLogMessage(msg)
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
