package api

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"brightrock.co.za/brgeo/model"
)

// GetGeoInfo accepts an IP Address to perform a lookup
// of the geo-location infomation of the IP Address
// form KeyCDN's Geo service.
// This service is rate limited to a maximum of 3 calls per second.
func GetGeoInfo(ipaddress string) (model.GeoData, int) {

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
	maxRetries := 3
	baseInterval := 500 * time.Millisecond
	retryFactor := 2.0
	retry := 0

	for i := 0; i < maxRetries; i++ {
		resp, err = client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}

		// If this wasn't the last attempt, sleep for a while before retrying
		if i < maxRetries-1 {
			sleepDuration := time.Duration(float64(baseInterval) * float64(i+1) * retryFactor)
			retry = i
			slog.Info("Sleeping on error", "duration", sleepDuration)
			time.Sleep(sleepDuration)
			continue
		} else {
			// If this was the last attempt, handle the error (e.g., by panicking)
			panic(err)
		}
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var geoResponse model.Response
	err = json.Unmarshal(body, &geoResponse)
	if err != nil {
		panic(err)
	}

	return geoResponse.Data.Geo, retry
}
