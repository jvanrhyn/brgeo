package api

import (
	"os"
	"testing"

	"brightrock.co.za/brgeo/model"
)

func TestAddCacheItem(t *testing.T) {

	os.Setenv("CACHE_TIMEOUT_SEC", "60")
	Cache.Flush()

	AddCacheItem("1", &model.LookupResponse{})

	if Cache.ItemCount() != 1 {
		t.Error("Cache item count should be 1")
	}
}

func TestAddMultipleCacheItem(t *testing.T) {

	os.Setenv("CACHE_TIMEOUT_SEC", "60")
	Cache.Flush()
	AddCacheItem("1", &model.LookupResponse{})
	AddCacheItem("2", &model.LookupResponse{})
	AddCacheItem("3", &model.LookupResponse{})

	if Cache.ItemCount() != 3 {
		t.Error("Cache item count should be 3")
	}
}

func TestAddMultipleWithDuplicateCacheItem(t *testing.T) {

	os.Setenv("CACHE_TIMEOUT_SEC", "60")
	Cache.Flush()

	AddCacheItem("1", &model.LookupResponse{})
	AddCacheItem("1", &model.LookupResponse{})
	AddCacheItem("3", &model.LookupResponse{})

	if Cache.ItemCount() != 2 {
		t.Errorf("Cache item count should be 2 but was %d", Cache.ItemCount())
	}
}
