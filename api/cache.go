package api

import (
	"errors"
	"log/slog"
	"os"
	"strconv"
	"time"

	"brightrock.co.za/brgeo/model"
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

// init function is called before the main function
func init() {
	cacheTimeout, err = strconv.Atoi(os.Getenv("CACHE_TIMEOUT_SEC"))
	if err != nil {
		slog.Info("Could not retrieve CACHE_TIMEOUT_SEC", "error", err)
		cacheTimeout = 60
	}
}

// GetCacheById retrieves an item from the cache for the given key
func GetCacheById(id string) (*model.LookupResponse, error) {
	item, found := Cache.Get(id)
	if found {
		return item.(*model.LookupResponse), nil
	}
	return &model.LookupResponse{}, errors.New("not found")
}

// SetCacheItem sets an item in the cache for the given key
func SetCacheItem(id string, data *model.LookupResponse) {

	duration := time.Duration(cacheTimeout) * time.Second
	slog.Info("Cache durations set", "duration", duration)

	Cache.Set(id, data, duration)
}
