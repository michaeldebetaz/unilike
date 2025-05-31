package scrapper

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/michaeldebetaz/unilike/internal/cache"
)

func GetHtml(url string, cache cache.Cache) (string, error) {
	if html, ok := cache.Get(url); ok {
		slog.Info("Cache hit", "url", url)
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
	cache.Set(url, html)

	return html, nil
}
