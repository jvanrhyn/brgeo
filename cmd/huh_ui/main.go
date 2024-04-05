package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"

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

	var sb strings.Builder

	_, _ = fmt.Fprintf(&sb, "\n\nIP Address Searched : %s", ipAddress)
	_, _ = fmt.Fprintf(&sb, "\n\nIP City : %s", geoInfo.City)
	_, _ = fmt.Fprintf(&sb, "\n\nIP Address Searched : %s", geoInfo.Region)
	_, _ = fmt.Fprintf(&sb, "\n\nIP Address Searched : %s", geoInfo.Country)

	fmt.Println(
		lipgloss.NewStyle().
			Width(60).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2).
			Render(sb.String()),
	)
}

func getGeoInfo(ipAddress string) (*GeoInfo, error) {
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
