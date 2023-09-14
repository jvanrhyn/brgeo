package api

import (
	"fmt"
	"log/slog"
	"os"

	"brightrock.co.za/brgeo/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDatabase() {
	dsn := os.Getenv("CONNECTION")
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
		slog.Error("error while recording lookuo", "error", err)
	}
}
