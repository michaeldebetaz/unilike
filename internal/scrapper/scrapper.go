package scrapper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Cache map[string]string

func LoadCache() (Cache, error) {
	cache := make(Cache)

	if _, err := os.Stat("cache"); errors.Is(err, os.ErrNotExist) {
		return cache, nil
	}

	bytes, err := os.ReadFile("cache")
	if err != nil {
		return cache, fmt.Errorf("Error reading cache from file: %v", err)
	}
	if err = json.Unmarshal(bytes, &cache); err != nil {
		return cache, fmt.Errorf("Error unmarshalling cache: %v", err)
	}

	return cache, nil
}

func (c *Cache) get(url string) (string, bool) {
	html, ok := (*c)[url]
	return html, ok
}

func (c *Cache) set(url string, html string) {
	(*c)[url] = html
}

func (c *Cache) Save() error {
	fmt.Println("Saving cache to file")

	bytes, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("Error marshalling cache: %v", err)
	}
	err = os.WriteFile("cache", bytes, 0644)
	if err != nil {
		return fmt.Errorf("Error writing cache to file: %v", err)
	}
	return nil
}

func GetHtml(url string, cache Cache) (string, error) {
	if html, ok := cache.get(url); ok {
		fmt.Printf("Cache hit for %s\n", url)
		return html, nil
	}

	fmt.Printf("Visiting %s\n", url)

	res, err := http.Get(url)
	if err != nil {
		err := fmt.Errorf("Error while doing Get request: %v", err)
		return "", err
	}
	defer res.Body.Close()

	fmt.Printf("Response status: %s\n", res.Status)

	if res.StatusCode > 299 {
		err := fmt.Errorf("Error: %s", res.Status)
		return "", err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		err := fmt.Errorf("Error while reading body: %v", err)
		return "", err
	}

	html := string(body)
	cache.set(url, html)

	return html, nil
}
