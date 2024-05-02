package api

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/go-errors/errors"
	"github.com/jvanrhyn/brgeo/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// InitDatabase initializes the database connection and migrates the LookupRequest model.
//
// It retrieves the DSN from the environment variable "CONNECTION" and logs the connection details.
//
// It then opens a connection to the database using GORM and configures it. If an error occurs during this process, it panics.
//
// After the connection is established, it attempts to auto-migrate the LookupRequest model. If an error occurs during this process, it logs the error and continues.
//
// This function is called once at the start of the application to set up the database.
func InitDatabase() {
	dsn := os.Getenv("CONNECTION")
	slog.Info("Initializing database", "connection", dsn)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&model.LookupRequest{})
	if err != nil {
		fmt.Println(err)
	}
}

// Record records a new lookup request in the database.
//
// Parameters:
//   - lookupRequest: A pointer to a model.LookupRequest object containing the details of the lookup request.
//
// Returns:
//   - nil if the lookup request is successfully recorded, otherwise an error.
//
// Example:
//
//	```
//	lookupRequest := &model.LookupRequest{/* details */}
//	err := Record(lookupRequest)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	```
func Record(lookupRequest *model.LookupRequest) error {
	tx := db.Create(lookupRequest)
	if tx.Error != nil {
		stack := err.(*errors.Error).ErrorStack()
		slog.Error("error while recording lookup", "error", err, "stacktrace", stack)
		return err
	}
	return nil
}
