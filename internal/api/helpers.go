// Package api provides functions to get the environment file path.
package api

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
)

// GetEnvFilePath returns the path to the .env file in the current executable's directory.
// If the executable's path is not available, it falls back to the caller's path.
//
// If an error occurs while retrieving the executable's path, it logs the error and returns an empty string.
//
// The function does not return any error, as it only logs errors and does not propagate them.
//
// Example:
//
//	envPath := api.GetEnvFilePath()
//	fmt.Println(envPath)
func GetEnvFilePath() string {

	ex, err := os.Executable()
	if err != nil {
		slog.Error(err.Error())
	}

	if ex != "" {
		exPath := filepath.Dir(ex)
		envPath := filepath.Join(exPath, ".env")
		return envPath
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("No caller information")
	}

	dir := filepath.Dir(filename)
	envPath := filepath.Join(dir, ".env")
	return envPath
}
