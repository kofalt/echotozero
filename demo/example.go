package main

import (
	"net/http"
	"os"

	"github.com/kofalt/echotozero"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

func hello(c echo.Context) error {
	c.Logger().Info("Info from echo")
	c.Logger().Warn("Warn from echo")

	return c.String(http.StatusOK, "Hello, World!")
}

func main() {
	// Use any zerolog instance: this one uses pretty printing
	zl := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()

	// This would instead use json format
	// zl := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Convert to an echo inferface
	adapter := echotozero.New(zl)

	// Use any echo instance
	e := echo.New()

	// Log echo messages, and hide its marketing copy
	e.Logger = adapter
	e.HideBanner = true

	// Basic middleware that log requests
	// It's easy to make your own: just steal middleware.go
	e.Use(echotozero.Middleware(adapter))

	// Routes
	e.GET("/", hello)

	// Start
	e.Logger.Fatal(e.Start(":8080"))
}
