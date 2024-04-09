package api

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
)

// GetEnvFilePath retrieves the absolute file path of the .env file in the directory
// where the calling function is located.
// If the caller information cannot be obtained, a message will be printed to the console.
// The function combines the directory path and the name of the .env file using filepath.Join.
// Returns the absolute file path of the .env file.
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
