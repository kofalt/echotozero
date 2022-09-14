package echotozero

import (
	"strconv"
	"time"

	"github.com/rs/zerolog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// MiddleWare logs successful requests to debug, failed req to error, and never skips.
func Middleware(logger *Logger) echo.MiddlewareFunc {
	return MiddlewareWithOptions(
		logger,
		zerolog.DebugLevel,
		zerolog.ErrorLevel,
		middleware.DefaultSkipper,
	)
}

func MiddlewareWithOptions(logger *Logger, successLvl, failLvl zerolog.Level, skipper middleware.Skipper) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {

			// Allows some requests to skip logging
			if skipper(c) {
				return next(c)
			}

			start := time.Now()
			req := c.Request()
			res := c.Response()

			// Invokes next middleware or handler
			err := next(c)

			var evt *zerolog.Event
			if err != nil {
				// Invokes registered HTTP error handler, if any
				c.Error(err)

				evt = logger.ZLog.WithLevel(failLvl).Err(err)
			} else {
				evt = logger.ZLog.WithLevel(successLvl)
			}

			cl := req.Header.Get(echo.HeaderContentLength)
			if cl == "" {
				cl = "0"
			}

			elapsed := time.Since(start)
			evt.Dur("elapsed", elapsed)
			evt.Str("host", req.Host)
			evt.Str("method", req.Method)
			evt.Str("remote_ip", c.RealIP())
			evt.Str("req_bytes", cl)
			evt.Str("res_bytes", strconv.FormatInt(res.Size, 10))
			evt.Int("status", res.Status)
			evt.Str("uri", req.RequestURI)
			evt.Str("user_agent", req.UserAgent())

			evt.Msg("HTTP")
			return err
		}
	}
}
