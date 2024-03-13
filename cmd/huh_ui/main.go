package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/charmbracelet/huh/spinner"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
)

type (
	GeoInfo struct {
		IPAddress string
		City      string `json:"city"`
		Region    string `json:"region"`
		Country   string `json:"country"`
	}
)

var ipAddress string

func main() {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("IP Address").
				Value(&ipAddress)))

	err := form.Run()
	if err != nil {
		log.Error("Error", err)
	}

	var geoInfo *GeoInfo

	action := func() {
		geoInfo, err = getGeoInfo(ipAddress)
		if err != nil {
			log.Error(err)
		}
	}

	_ = spinner.New().
		Title("Looking up IP-Address").
		Action(action).
		Run()

	form = huh.NewForm(
		huh.NewGroup(
			huh.NewText().Title("IP Address").Lines(1).
				Value(&ipAddress),
			huh.NewText().Title("Country").Lines(1).
				Value(&geoInfo.Country),
			huh.NewText().Title("Region").Lines(1).
				Value(&geoInfo.Region),
			huh.NewText().Title("City").Lines(1).
				Value(&geoInfo.City)))

	err = form.Run()
	if err != nil {
		log.Error(err)
	}
}

func getGeoInfo(ipAddress string) (*GeoInfo, error) {
	time.Sleep(time.Second * 2)
	resp, err := http.Get(fmt.Sprintf("http://localhost:3000/api/lookup/%s", ipAddress))
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var geoInfo GeoInfo
	err = json.Unmarshal(body, &geoInfo)
	if err != nil {
		return nil, err
	}

	return &geoInfo, nil
}
