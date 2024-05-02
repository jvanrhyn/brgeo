package api

import (
	"github.com/jvanrhyn/brgeo/model"
	"os"
	"testing"
)

func TestAddCacheItem(t *testing.T) {

	os.Setenv("CACHE_TIMEOUT_SEC", "60")
	Cache.Flush()

	_ = AddCacheItem("1", &model.LookupResponse{})

	if Cache.ItemCount() != 1 {
		t.Error("Cache item count should be 1")
	}
}

func TestAddMultipleCacheItem(t *testing.T) {

	os.Setenv("CACHE_TIMEOUT_SEC", "60")
	Cache.Flush()
	_ = AddCacheItem("1", &model.LookupResponse{})
	_ = AddCacheItem("2", &model.LookupResponse{})
	_ = AddCacheItem("3", &model.LookupResponse{})

	if Cache.ItemCount() != 3 {
		t.Error("Cache item count should be 3")
	}
}

func TestAddMultipleWithDuplicateCacheItem(t *testing.T) {

	os.Setenv("CACHE_TIMEOUT_SEC", "60")
	Cache.Flush()

	_ = AddCacheItem("1", &model.LookupResponse{})
	_ = AddCacheItem("1", &model.LookupResponse{})
	_ = AddCacheItem("3", &model.LookupResponse{})

	if Cache.ItemCount() != 2 {
		t.Errorf("Cache item count should be 2 but was %d", Cache.ItemCount())
	}
}
