package main

import (
	"os"
	"time"

	"log/slog"

	"github.com/joho/godotenv"
	"github.com/jvanrhyn/brgeo/api"
	"github.com/jvanrhyn/brgeo/controller"
	"github.com/lmittmann/tint"
)

func init() {
	// Read Configuration data from the .env file in the project
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file : " + err.Error())
	}
}

func main() {
	w := os.Stderr

	// set global logger with custom options
	slog.SetDefault(slog.New(
		tint.NewHandler(w, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.RFC3339Nano,
		}),
	))
	slog.Info("Starting the application")

	api.InitDatabase()
	controller.StartAndServe()
}
