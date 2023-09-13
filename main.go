package main

import (
	"brightrock.co.za/brgeo/controller"
	"github.com/joho/godotenv"
)

func init() {
	// Read Configuration data from the .env file in the project
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file : " + err.Error())
	}
}

func main() {
	controller.StartAndServe()
}
