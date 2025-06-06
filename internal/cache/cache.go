package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"sync"
)

const FILE_PATH = "./.cache"

type Cache struct {
	Data  map[string]string
	Mutex *sync.RWMutex
}

func Load() (Cache, error) {
	cache := Cache{
		Data:  make(map[string]string),
		Mutex: &sync.RWMutex{},
	}

	if _, err := os.Stat(FILE_PATH); errors.Is(err, os.ErrNotExist) {
		return cache, nil
	}

	bytes, err := os.ReadFile(FILE_PATH)
	if err != nil {
		return cache, fmt.Errorf("Error reading cache from file: %v", err)
	}
	if err = json.Unmarshal(bytes, &cache.Data); err != nil {
		return cache, fmt.Errorf("Error unmarshalling cache: %v", err)
	}

	return cache, nil
}

func (c *Cache) Save() error {
	slog.Info("Saving cache to file...")

	bytes, err := json.Marshal(c.Data)
	if err != nil {
		return fmt.Errorf("Error marshalling cache: %v", err)
	}
	err = os.WriteFile(FILE_PATH, bytes, 0644)
	if err != nil {
		return fmt.Errorf("Error writing cache to file: %v", err)
	}

	slog.Info("Cache saved successfully")

	return nil
}

func (c *Cache) Get(url string) (string, bool) {
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()

	html, ok := (*c).Data[url]
	return html, ok
}

func (c *Cache) Set(url string, html string) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	(*c).Data[url] = html
}
