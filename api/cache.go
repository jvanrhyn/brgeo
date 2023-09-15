package api

import (
	"errors"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/jvanrhyn/brgeo/model"
	"github.com/patrickmn/go-cache"
)

var (
	Cache        = cache.New(5*time.Minute, 5*time.Minute)
	cacheTimeout int
	err          error
)

// CacheItem struct holds the string key and
// the cached pointer of the item being cached
type CacheItem []struct {
	ID   string                `json:"id"`
	Data *model.LookupResponse `json:"data"`
}

// GetCacheById retrieves an item from the cache for the given key
func GetCacheById(id string) (*model.LookupResponse, error) {
	slog.Info("Retrieving item from cache", "id", id)
	item, found := Cache.Get(id)
	if found {
		return item.(*model.LookupResponse), nil
	}
	return &model.LookupResponse{}, errors.New("not found")
}

// AddCacheItem sets an item in the cache for the given key
func AddCacheItem(id string, data *model.LookupResponse) error {

	if cacheTimeout == 0 {
		getTimeoutSeconds()
	}

	duration := time.Duration(cacheTimeout) * time.Second
	slog.Info("Cache durations set", "duration", duration)

	err := Cache.Add(id, data, duration)
	if err != nil {
		slog.Error("Error adding item to cache", "error", err)
		return err
	}
	slog.Info("Item added to cache", "id", id, "cache", Cache.Items())
	return nil
}

// init function is called before the main function
func getTimeoutSeconds() {
	ct := os.Getenv("CACHE_TIMEOUT_SEC")
	slog.Info("Initializing cache", "timeout", ct)
	cacheTimeout, err = strconv.Atoi(ct)
	if err != nil {
		slog.Info("Could not retrieve CACHE_TIMEOUT_SEC", "error", err)
		cacheTimeout = 60
	}
}
