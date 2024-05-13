package api

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/vincentfree/opentelemetry/otelslog"

	"github.com/jvanrhyn/brgeo/model"
)

var (
	logr = otelslog.New()
)

func init() {

}

// GetGeoInfo accepts an IP Address to perform a lookup
// of the geolocation information of the IP Address
// form KeyCDNs' Geo service.
// This service is rate limited to a maximum of 3 calls per second.
func GetGeoInfo(ipaddress string) (model.GeoData, int) {

	logr.Info("GetGeoInfo called")

	client := &http.Client{}
	path := os.Getenv("SERVICE_URL") + "?host=" + ipaddress
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		panic(err)
	}

	// Setting the User-Agent header
	req.Header.Set("User-Agent", os.Getenv("USER_AGENT"))

	// Define a re-usable response object needed to
	// the retry pattern implemented
	var resp *http.Response

	// Setup for a backoff retry pattern
	maxRetries, _ := strconv.Atoi(os.Getenv("MAX_RETRIES"))
	if maxRetries < 3 {
		maxRetries = 3
	}

	slog.Info("Max retries", "retries", maxRetries)

	baseInterval := 500 * time.Millisecond
	retryFactor := 2.0
	retry := 0

	for i := 0; i < maxRetries; i++ {
		slog.Info("Attempts Counter", "attempt", i)
		resp, err = client.Do(req)
		if err == nil {
			defer resp.Body.Close() // Ensure the body is always closed
			if resp.StatusCode == http.StatusOK {
				break
			}
		}

		// If this wasn't the last attempt, sleep for a while before retrying
		if i < maxRetries-1 {
			sleepDuration := time.Duration(float64(baseInterval) * float64(i+1) * retryFactor)
			retry = i

			logr.Info("Sleeping on error", "duration", sleepDuration, "error", err)
			slog.Info("Sleeping on error", "duration", sleepDuration, "error", err)
			time.Sleep(sleepDuration)
			continue
		} else {
			// If this was the last attempt, handle the error (e.g., by panicking)
			panic(err)
		}
	}

	if resp != nil {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error(err.Error())
		}
		var geoResponse model.Response
		err = json.Unmarshal(body, &geoResponse)
		if err != nil {
			slog.Error(err.Error())
		}

		return geoResponse.Data.Geo, retry
	}
	return model.GeoData{}, 0
}
