package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

type (
	GeoInfo struct {
		IPAddress string
		City      string `json:"city"`
		Region    string `json:"region"`
		Country   string `json:"country"`
	}

	Model struct {
		title     string
		geoInfo   GeoInfo
		textinput textinput.Model
		err       error
	}

	GeoResponseMsg struct {
		GeoInfo GeoInfo
		Err     error
	}
)

func main() {
	m := NewModel()

	prog := tea.NewProgram(m, tea.WithAltScreen())

	_, err := prog.Run()
	if err != nil {
		slog.Error(err.Error())
	}
}

func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter IP Address"
	ti.Reset()
	ti.Focus()

	return Model{
		title:     "Geo-location lookup",
		textinput: ti,
		geoInfo:   GeoInfo{},
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			v := m.textinput.Value()
			if v == "quit" || v == "" {
				os.Exit(0)
			}

			if !ValidateIP(v) {
				m.err = errors.New("invalid ip address")
				return m, nil
			}

			return m, handGeoleLookup(v)
		case tea.KeyCtrlC:
			return m, tea.Quit
		}

	case GeoResponseMsg:
		m.geoInfo = msg.GeoInfo
		if msg.Err != nil {
			m.err = msg.Err
			return m, tea.Quit
		}

		return m, nil
	}

	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	s := m.textinput.View() + "\n\n"
	m.textinput.Reset()

	if m.err != nil {
		s += "Error : " + m.err.Error()
		return s
	}

	if m.geoInfo.Country != "" {
		s += "Looking up : " + m.geoInfo.IPAddress + "\n"
		s += "Region  : " + m.geoInfo.Region + "\n"
		s += "City    : " + m.geoInfo.City + "\n"
		s += "Country : " + m.geoInfo.Country + "\n"
	}
	return s
}

func handGeoleLookup(v string) tea.Cmd {
	return func() tea.Msg {
		geo, err := getGeoInfo(v)
		if err == nil {
			log.Error(err)
		}

		return GeoResponseMsg{
			GeoInfo: geo,
			Err:     err,
		}
	}
}

func getGeoInfo(ipAddress string) (GeoInfo, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:3000/api/lookup/%s", ipAddress))
	if err != nil {
		return GeoInfo{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return GeoInfo{}, err
	}

	var geoInfo GeoInfo
	err = json.Unmarshal(body, &geoInfo)
	if err != nil {
		return GeoInfo{}, err
	}

	geoInfo.IPAddress = ipAddress
	return geoInfo, nil
}

func ValidateIP(ip string) bool {
	// Try parsing as an IP address using net.ParseIP
	address := net.ParseIP(ip)
	if address != nil {
		// Successfully parsed as IP. Let's do additional checks if needed:
		if strings.Contains(ip, ".") {
			// It's likely IPv4
			return address.To4() != nil
		} else if strings.Contains(ip, ":") {
			// It's likely IPv6
			return true // Basic IPv6 check for now
		}
	}
	return false // Failed to parse or additional checks didn't pass
}
