package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
)

type Cache map[string]string

func Load() (Cache, error) {
	cache := make(Cache)

	filePath := ".cache"

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return cache, nil
	}

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return cache, fmt.Errorf("Error reading cache from file: %v", err)
	}
	if err = json.Unmarshal(bytes, &cache); err != nil {
		return cache, fmt.Errorf("Error unmarshalling cache: %v", err)
	}

	return cache, nil
}

func (c *Cache) Save() error {
	slog.Info("Saving cache to file...")

	bytes, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("Error marshalling cache: %v", err)
	}
	err = os.WriteFile("cache", bytes, 0644)
	if err != nil {
		return fmt.Errorf("Error writing cache to file: %v", err)
	}

	slog.Info("Cache saved successfully")

	return nil
}

func (c *Cache) Get(url string) (string, bool) {
	html, ok := (*c)[url]
	return html, ok
}

func (c *Cache) Set(url string, html string) {
	(*c)[url] = html
}
