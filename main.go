package main

import (
	"fmt"

	"brightrock.co.za/brgeo/api"
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

	ipaddress := "169.1.245.236"
	geodata, retry := api.GetGeoInfo(ipaddress)

	fmt.Printf("You live in %s\nRetried %d times", geodata.City, retry)

}
