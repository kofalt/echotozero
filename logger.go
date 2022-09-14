package echotozero

import (
	"context"
	"fmt"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// Enforce interface
var _ echo.Logger = (*Logger)(nil)

// Logger wraps `zerolog.Logger` and implements `echo.Logger`
type Logger struct {
	ZLog zerolog.Logger

	// Level used when Echo does not provide one
	DefaultLevel zerolog.Level
}

func New(log zerolog.Logger) *Logger {
	return &Logger{
		ZLog:         log,
		DefaultLevel: zerolog.DebugLevel,
	}
}

//
// Interface utility
//

func (l *Logger) Output() io.Writer {
	return l.ZLog
}

func (l *Logger) Level() echoLevel {
	return MapZeroToEcho[l.ZLog.GetLevel()]
}

func (l *Logger) WithContext(ctx context.Context) context.Context {
	return l.ZLog.WithContext(ctx)
}

//
// Intentionally omitting some parts of the Echo interface
// Use zerolog features instead
//

func (l *Logger) SetOutput(newOut io.Writer) {}
func (l *Logger) SetHeader(h string)         {}
func (l *Logger) SetLevel(level echoLevel)   {}
func (l *Logger) SetPrefix(newPrefix string) {}
func (l *Logger) Prefix() string {
	return ""
}

//
// Message
//

func (l *Logger) Debug(i ...interface{}) {
	l.ZLog.Debug().Msg(fmt.Sprint(i...))
}

func (l *Logger) Info(i ...interface{}) {
	l.ZLog.Info().Msg(fmt.Sprint(i...))
}

func (l *Logger) Warn(i ...interface{}) {
	l.ZLog.Warn().Msg(fmt.Sprint(i...))
}

func (l *Logger) Error(i ...interface{}) {
	l.ZLog.Error().Msg(fmt.Sprint(i...))
}

func (l *Logger) Fatal(i ...interface{}) {
	l.ZLog.Fatal().Msg(fmt.Sprint(i...))
}

func (l *Logger) Panic(i ...interface{}) {
	l.ZLog.Panic().Msg(fmt.Sprint(i...))
}

func (l *Logger) Print(i ...interface{}) {
	l.ZLog.WithLevel(l.DefaultLevel).Msg(fmt.Sprint(i...))
}

//
// Format message
//

func (l *Logger) Debugf(format string, i ...interface{}) {
	l.ZLog.Debug().Msgf(format, i...)
}

func (l *Logger) Infof(format string, i ...interface{}) {
	l.ZLog.Info().Msgf(format, i...)
}

func (l *Logger) Warnf(format string, i ...interface{}) {
	l.ZLog.Warn().Msgf(format, i...)
}

func (l *Logger) Errorf(format string, i ...interface{}) {
	l.ZLog.Error().Msgf(format, i...)
}

func (l *Logger) Fatalf(format string, i ...interface{}) {
	l.ZLog.Fatal().Msgf(format, i...)
}

func (l *Logger) Panicf(format string, i ...interface{}) {
	l.ZLog.Panic().Msgf(format, i...)
}

func (l *Logger) Printf(format string, i ...interface{}) {
	l.ZLog.WithLevel(l.DefaultLevel).Msgf(format, i...)
}

//
// JSON message
//

func (l *Logger) Debugj(j JSON) {
	l.logJSON(l.ZLog.Debug(), j)
}

func (l *Logger) Infoj(j JSON) {
	l.logJSON(l.ZLog.Info(), j)
}

func (l *Logger) Warnj(j JSON) {
	l.logJSON(l.ZLog.Warn(), j)
}

func (l *Logger) Errorj(j JSON) {
	l.logJSON(l.ZLog.Error(), j)
}

func (l *Logger) Fatalj(j JSON) {
	l.logJSON(l.ZLog.Fatal(), j)
}

func (l *Logger) Panicj(j JSON) {
	l.logJSON(l.ZLog.Panic(), j)
}

func (l *Logger) Printj(j JSON) {
	l.logJSON(l.ZLog.WithLevel(l.DefaultLevel), j)
}

func (l *Logger) logJSON(event *zerolog.Event, j JSON) {
	for k, v := range j {
		event = event.Interface(k, v)
	}

	event.Msg("")
}
