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
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
		spinner   spinner.Model
		err       error
		loading   bool
	}

	GeoResponseMsg struct {
		GeoInfo GeoInfo
		Err     error
	}
)

func main() {

	w := os.Stderr
	handler := log.New(w)
	handler.SetLevel(log.DebugLevel)
	handler.SetTimeFormat(time.Kitchen)
	handler.SetReportTimestamp(true)

	// set global logger with custom options
	slog.SetDefault(slog.New(
		handler))

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

	s := spinner.New()
	s.Spinner = spinner.Dot

	return Model{
		title:     "Geo-location lookup",
		textinput: ti,
		geoInfo:   GeoInfo{},
		spinner:   s,
		loading:   false,
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
			if !validateIP(v) {
				m.err = errors.New("invalid ip address")
			} else {
				m.err = nil // Clear the previous error
				cmd = handleGeoLookup(v)
				m.loading = true
			}

			m.textinput.Reset() // Reset the input field for new input
			return m, cmd
		case tea.KeyCtrlC:
			return m, tea.Quit
		}

	case GeoResponseMsg:
		m.geoInfo = msg.GeoInfo
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			m.err = nil
			m.textinput.Reset()
		}
		m.loading = false
		m.textinput.Reset()
		return m, nil

	case spinner.TickMsg: // Update the spinner on each tick
		m.spinner, cmd = m.spinner.Update(msg)
		if m.loading {
			m.spinner.Tick()
		}
		return m, cmd
	}

	m.textinput, cmd = m.textinput.Update(msg)
	m.spinner, _ = m.spinner.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	// Define styles
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#008ACC"))
	staticTextStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")) // Grey
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#EF9704"))
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#C62D16"))

	titleView := titleStyle.Render("Geo-location Lookup")
	inputView := ""

	if !m.loading {
		// Build the view string with styles
		inputView += "\n\n" + m.textinput.View() + "\n\n"
	}

	if m.err != nil {
		errorView := staticTextStyle.Render("Error : ") + errorStyle.Render(m.err.Error())
		return titleView + errorView + inputView
	}

	spinnerView := ""
	if m.loading {
		m.spinner.Update(m.spinner.Tick())
		spinnerView = "\n\n" + m.spinner.View() + "\n\n"
	}

	dataView := ""

	if m.geoInfo.Country != "" {

		dataView += "\n\n"
		dataView += staticTextStyle.Render("Looking up: ") + valueStyle.Render(m.geoInfo.IPAddress) + "\n" +
			staticTextStyle.Render("Region: ") + valueStyle.Render(m.geoInfo.Region) + "\n" +
			staticTextStyle.Render("City: ") + valueStyle.Render(m.geoInfo.City) + "\n" +
			staticTextStyle.Render("Country: ") + valueStyle.Render(m.geoInfo.Country)
	}
	return titleView + spinnerView + dataView + inputView
}

func handleGeoLookup(v string) tea.Cmd {
	return func() tea.Msg {
		geo, err := getGeoInfo(v)
		if err != nil {
			log.Error(err)
		}

		return GeoResponseMsg{
			GeoInfo: geo,
			Err:     err,
		}
	}
}

func getGeoInfo(ipAddress string) (GeoInfo, error) {
	var geoInfo GeoInfo
	var err error
	maxRetries := 3               // Maximum number of retries
	retryDelay := 1 * time.Second // Initial delay, doubled on each retry

	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, err := http.Get(fmt.Sprintf("http://localhost:3000/api/lookup/%s", ipAddress))
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err == nil {
					err = json.Unmarshal(body, &geoInfo)
					if err == nil {
						geoInfo.IPAddress = ipAddress
						return geoInfo, nil // Success
					}
				}
			}
		}

		time.Sleep(retryDelay)
		retryDelay *= 2 // Exponential backoff
	}

	return GeoInfo{}, fmt.Errorf("failed to get geo info after %d attempts: %v", maxRetries, err)
}

func validateIP(ip string) bool {
	address := net.ParseIP(ip)
	if address == nil {
		return false
	}

	// Check if the IP address is IPv4 or IPv6
	if ipv4 := address.To4(); ipv4 != nil {
		return true
	} else if ipv6 := address.To16(); ipv6 != nil {
		return true
	}

	return false
}
