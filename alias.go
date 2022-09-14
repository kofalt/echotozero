package echotozero

import (
	echoLib "github.com/labstack/gommon/log"
	"github.com/rs/zerolog"
)

// Echo logging uses a separate "WIP" package for a few things
// Alias those away for clarity

type echoLevel = echoLib.Lvl // uint8
type JSON = echoLib.JSON     // map[string]any

var (
	MapEchoToZero = map[echoLevel]zerolog.Level{
		echoLib.DEBUG: zerolog.DebugLevel,
		echoLib.INFO:  zerolog.InfoLevel,
		echoLib.WARN:  zerolog.WarnLevel,
		echoLib.ERROR: zerolog.ErrorLevel,
		echoLib.OFF:   zerolog.NoLevel,
	}

	MapZeroToEcho = map[zerolog.Level]echoLevel{
		zerolog.DebugLevel: echoLib.DEBUG,
		zerolog.InfoLevel:  echoLib.INFO,
		zerolog.WarnLevel:  echoLib.WARN,
		zerolog.ErrorLevel: echoLib.ERROR,
		zerolog.NoLevel:    echoLib.OFF,
	}
)
