package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"log/slog"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/jvanrhyn/brgeo/controller"
	"github.com/jvanrhyn/brgeo/internal/api"
)

// init is the function that initializes the application by reading Configuration data
// from the .env file in the project.
// If the .env file cannot be found in the project directory, it attempts to find the file
// using the getEnvFilePath function.
// If the .env file still cannot be found, an error message is printed to the console and the application panics.
// The function uses the godotenv package to load the .env file.
func init() {

	// Read Configuration data from the .env file in the project
	err := godotenv.Load()
	if err != nil {
		envPath := getEnvFilePath()
		err = godotenv.Load(envPath)
		if err != nil {
			panic("Error loading .env file : " + err.Error())
		}
	}
}

// main is the entry point of the application.
// It initializes the logger, sets up the global logger with custom options,
// and starts the application by calling the InitDatabase and StartAndServe functions.
// It uses the slog package for logging.
//
// The logger is initialized with the log.New function and set with the provided options.
// The SetLevel function sets the logging level to Debug.
// The SetTimeFormat function sets the time format to time.Kitchen.
// The SetReportTimestamp function enables reporting timestamps in logs.
//
// The slogger package is used to set the default logger to the created logger.
//
// The Info log message "Starting the application" is printed using the slog.Info function.
//
// The Debug log message "InitDatabase called" is printed using the slog.Debug function.
// The InitDatabase function (api.InitDatabase) is then called to initialize the database.
//
// The StartAndServe function (controller.StartAndServe) is called to start and serve the application.
//
// This function does not return anything.
func main() {
	w := os.Stderr
	handler := log.New(w)
	handler.SetLevel(log.DebugLevel)
	handler.SetTimeFormat(time.Kitchen)
	handler.SetReportTimestamp(true)

	// set global logger with custom options
	slog.SetDefault(slog.New(
		handler))
	slog.Info("Starting the application")

	slog.Debug("InitDatabase called")

	api.InitDatabase()
	controller.StartAndServe()
}

// getEnvFilePath retrieves the absolute file path of the .env file in the directory
// where the calling function is located.
// If the caller information cannot be obtained, a message will be printed to the console.
// The function combines the directory path and the name of the .env file using filepath.Join.
// Returns the absolute file path of the .env file.
func getEnvFilePath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("No caller information")
	}

	dir := filepath.Dir(filename)
	envPath := filepath.Join(dir, ".env")
	return envPath
}
