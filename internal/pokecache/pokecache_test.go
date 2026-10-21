package pokecache

import (
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	cache := NewCache(5 * time.Second)

	key := "https://example.com"
	value := []byte("testdata")

	cache.Add(key, value)

	actual, ok := cache.Get(key)
	if !ok {
		t.Fatal("expected to find key in cache")
	}

	if string(actual) != string(value) {
		t.Errorf("expected %q, got %q", string(value), string(actual))
	}

}

func TestGetMissingKey(t *testing.T) {
	cache := NewCache(5 * time.Second)

	_, ok := cache.Get("missing-key")
	if ok {
		t.Fatal("expected key not to be found")
	}

}

func TestReapLoop(t *testing.T) {
	interval := 50 * time.Millisecond
	cache := NewCache(interval)

	cache.Add("test-key", []byte("testdata"))

	_, ok := cache.Get("test-key")
	if !ok {
		t.Fatal("expected to find key before reaping")
	}

	time.Sleep(150 * time.Millisecond)

	_, ok = cache.Get("test-key")
	if ok {
		t.Fatal("expected key to be removed by reapLoop")
	}

}
