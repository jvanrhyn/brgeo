package api

import (
	"log/slog"
	"os"
	"testing"
)

func TestCanGetLocation(t *testing.T) {

	t.Parallel()

	_ = os.Setenv("USER_AGENT", "keycdn-tools:https://www.b.co.za")
	_ = os.Setenv("SERVICE_URL", "https://tools.keycdn.com/geo.json")
	_ = os.Setenv("MAX_RETRIES", "3")

	testCases := map[string]struct {
		value string
	}{
		"afrihost":   {value: "169.1.245.236"},
		"mazda":      {value: "172.67.188.116"},
		"capetalk":   {value: "52.85.24.88"},
		"vw":         {value: "75.2.9.61"},
		"brightrock": {value: "41.160.113.136"},
	}

	for n, tc := range testCases {
		n, tc := n, tc

		t.Run(n, func(t *testing.T) {
			t.Parallel()
			got, retry := GetGeoInfo(tc.value)
			slog.Info("Found ", "Country", got.CountryCode, "retries", retry)
			if got.CountryCode == "" {
				t.Error("expected city name to be populated")
			}

		})
	}
}
