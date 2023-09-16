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

func Record(lookupRequest *model.LookupRequest) {
	tx := db.Create(&lookupRequest)
	if tx.Error != nil {
		stack := err.(*errors.Error).ErrorStack()
		slog.Error("error while recording lookup", "error", err, "stacktrace", stack)
	}
}
