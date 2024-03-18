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

	m := New()

	prog := tea.NewProgram(m)

	_, err := prog.Run()
	if err != nil {
		slog.Error(err.Error())
	}
}

// New creates and returns a new Model instance with default values.
// It initializes a text input with a placeholder for entering an IP address,
// a spinner for indicating loading status, and sets loading to false.
func New() Model {
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

// Init is a method of the Model struct that initializes the model.
// It returns a command that starts the spinner ticking.
func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update is a method of the Model struct that handles incoming messages and updates the model accordingly.
// It returns the updated model and a command.
// It handles key presses, geolocation responses, and spinner ticks.
// When the Enter key is pressed, it validates the input IP address and initiates a geolocation lookup if the IP is valid.
// When a geolocation response is received, it updates the geolocation info in the model and clears any errors.
// When a spinner tick message is received, it updates the spinner in the model.
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
				m.err = nil
				m.loading = true
				cmd = handleGeoLookup(v)
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

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	m.textinput, cmd = m.textinput.Update(msg)
	m.spinner, _ = m.spinner.Update(msg)
	return m, cmd
}

// View is a method of the Model struct that returns a string representation of the model.
// It defines styles for different parts of the view and uses these styles to render the view.
// The view includes a title, an input field (when not loading), an error message (if any), a loading spinner (when loading), and the geolocation data (when available).
// The method uses fmt.Sprintf for string formatting.
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
		inputView += fmt.Sprintf("\n\n%s\n\n", m.textinput.View())
	}

	if m.err != nil {
		errorView := fmt.Sprintf("%s %s", staticTextStyle.Render("Error : "), errorStyle.Render(m.err.Error()))
		return fmt.Sprintf("%s%s%s", titleView, errorView, inputView)
	}

	spinnerView := ""
	if m.loading {
		spinnerView += fmt.Sprintf("\n\n%s Looking up Geo location.\n\n", m.spinner.View())
	}

	dataView := ""

	if m.geoInfo.Country != "" {
		dataView += fmt.Sprintf("\n\n%s %s\n%s %s\n%s %s\n%s %s",
			staticTextStyle.Render("Looking up: "), valueStyle.Render(m.geoInfo.IPAddress),
			staticTextStyle.Render("Region: "), valueStyle.Render(m.geoInfo.Region),
			staticTextStyle.Render("City: "), valueStyle.Render(m.geoInfo.City),
			staticTextStyle.Render("Country: "), valueStyle.Render(m.geoInfo.Country))
	}
	return fmt.Sprintf("%s%s%s%s", titleView, spinnerView, dataView, inputView)
}

// handleGeoLookup is a function that takes an IP address as a string and returns a command that fetches the geolocation information for that IP.
// It calls the getGeoInfo function to fetch the geolocation information and logs any errors that occur during this process.
// It then returns a GeoResponseMsg with the fetched geolocation information and any error that occurred.
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

// getGeoInfo is a function that takes an IP address as a string and returns the geolocation information for that IP and any error that occurred.
// It makes a GET request to a local API endpoint with the IP address and reads the response.
// If the response status is OK, it reads the response body and unmarshals it into a GeoInfo struct.
// If any error occurs during this process, it retries the request up to a maximum number of times with an exponential backoff delay.
// If it still fails after the maximum number of retries, it returns an empty GeoInfo struct and an error indicating the failure.
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

// validateIP is a function that takes an IP address as a string and returns a boolean indicating whether the IP address is valid.
// It uses the net.ParseIP function to parse the IP address and checks if the result is nil (indicating an invalid IP address).
// If the result is not nil, it checks if the IP address is IPv4 or IPv6 using the To4 and To16 methods respectively.
// If the IP address is either IPv4 or IPv6, it returns true; otherwise, it returns false.
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
